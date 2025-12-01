package monitor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adhocore/gronx"

	"dir-monitor-go/internal/config"
	"dir-monitor-go/internal/logger"
	"dir-monitor-go/internal/model"
)

const (
	DefaultMaxConcurrentOperations = 5
	DropLogThrottleTime            = 10 * time.Second
	CleanupInterval                = 1 * time.Minute
	DedupCacheExpiration             = 10 * time.Minute
)

type Monitor struct {
	config       *config.Config
	logger       *logger.Logger
	watcher      *FsnotifyWatcher
	stopChan     chan struct{}
	wg           sync.WaitGroup
	watchedDirs  map[string]bool
	workingDir   string
	mu           sync.Mutex
	eventChannel chan model.FileEvent
	specificDir  string

	dedupCache map[string]time.Time
	dedupMu    sync.Mutex

	dirTimers  map[string]*time.Timer
	dirBuffers map[string]map[string]model.FileEvent
	dirMu      sync.Mutex

	stopped int32

	dropLog map[string]time.Time
	dropMu  sync.Mutex

	cleanupStop chan struct{}

	opCtx    context.Context
	opCancel context.CancelFunc

	opSem chan struct{}
}

func NewMonitor(cfg *config.Config, log *logger.Logger) (*Monitor, error) {
	watcher := NewFsnotifyWatcher(log)
	if watcher == nil {
		return nil, fmt.Errorf("failed to create FsnotifyWatcher")
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}

	opMax := cfg.Settings.MaxConcurrentOperations
	if opMax <= 0 {
		opMax = DefaultMaxConcurrentOperations
	}

	opCtx, opCancel := context.WithCancel(context.Background())

	monitor := &Monitor{
		config:       cfg,
		logger:       log,
		watcher:      watcher,
		watchedDirs:  make(map[string]bool),
		stopChan:     make(chan struct{}),
		eventChannel: make(chan model.FileEvent, cfg.Settings.EventChannelBufferSize),
		workingDir:   workingDir,
		dedupCache:   make(map[string]time.Time),
		dirTimers:    make(map[string]*time.Timer),
		dirBuffers:   make(map[string]map[string]model.FileEvent),
		dropLog:      make(map[string]time.Time),
		cleanupStop:  make(chan struct{}),
		opCtx:        opCtx,
		opCancel:     opCancel,
		opSem:        make(chan struct{}, opMax),
	}

	return monitor, nil
}

func (m *Monitor) Start() error {
	m.logger.Info("[Monitor] Starting directory monitor")
	
	// 记录监控配置摘要
	enabledCount := 0
	for _, monitor := range m.config.Monitors {
		if monitor.Enabled {
			enabledCount++
		}
	}
	m.logger.Info("[Monitor] 监控配置摘要 - 总监控项: %d, 启用监控项: %d, 禁用监控项: %d", 
		len(m.config.Monitors), enabledCount, len(m.config.Monitors)-enabledCount)

	if err := m.startWatching(); err != nil {
		return fmt.Errorf("failed to start watching directories: %v", err)
	}

	m.wg.Add(1)
	go m.eventProcessor()

	m.wg.Add(1)
	go m.cleanupDaemon()

	m.logger.Info("[Monitor] Directory monitor started successfully")
	return nil
}

func (m *Monitor) Stop() error {
	m.logger.Info("Stopping directory monitor")

	atomic.StoreInt32(&m.stopped, 1)

	if m.opCancel != nil {
		m.opCancel()
	}

	close(m.stopChan)

	if m.cleanupStop != nil {
		close(m.cleanupStop)
	}

	if m.watcher != nil {
		m.watcher.Stop()
	}

	m.wg.Wait()

	m.logger.Info("Directory monitor stopped successfully")
	return nil
}

func (m *Monitor) startWatching() error {
	m.logger.Info("[Monitor] Starting to watch directories")

	dirsToWatch := make(map[string]bool)
	monitorInfo := make(map[string][]string) // 记录每个目录对应的监控项信息

	if m.specificDir != "" {
		dirsToWatch[m.specificDir] = true
	} else {
		for _, monitor := range m.config.Monitors {
			if monitor.Enabled {
				dirsToWatch[monitor.Directory] = true
				monitorInfo[monitor.Directory] = append(monitorInfo[monitor.Directory], 
					fmt.Sprintf("%s(%s)", monitor.Name, monitor.Command))
			}
		}
	}

	watchCount := 0
	for dir := range dirsToWatch {
		if err := m.watcher.Watch(dir, func(event model.FileEvent) {
			select {
			case m.eventChannel <- event:
				m.logger.Info("[Monitor] 文件事件已发送到通道: %s", event.Path)
			default:
				if m.shouldLogDrop(event.Path) {
					m.logger.Info("[Monitor] 事件通道已满，丢弃事件: %s", event.Path)
				}
			}
		}); err != nil {
			m.logger.Error("Failed to watch directory %s: %v", dir, err)
			continue
		}

		m.watchedDirs[dir] = true
		watchCount++
		
		// 记录详细的监控目录信息
		if monitors, exists := monitorInfo[dir]; exists {
			m.logger.Info("[Monitor] 开始监控目录: %s, 关联监控项: %v", dir, monitors)
		} else {
			m.logger.Info("[Monitor] 开始监控目录: %s", dir)
		}
	}

	if watchCount == 0 {
		return fmt.Errorf("[Monitor] No directories to watch")
	}

	m.logger.Info("[Monitor] Started watching %d directories", watchCount)
	return nil
}

func (m *Monitor) eventProcessor() {
	defer m.wg.Done()
	
	m.logger.Info("[Monitor] 事件处理器已启动，等待文件事件...")
	eventCount := 0

	for {
		select {
		case <-m.stopChan:
			m.logger.Info("[Monitor] 事件处理器正在停止，处理了 %d 个事件", eventCount)
			return
		case event := <-m.eventChannel:
			if atomic.LoadInt32(&m.stopped) == 1 {
				continue
			}

			eventCount++
			m.logger.Info("[Monitor] 事件处理器接收到第 %d 个事件: %s - %s", eventCount, event.Type, event.Path)
			m.processEvent(event)
		}
	}
}

func (m *Monitor) processEvent(event model.FileEvent) {
	m.logger.Info("[Monitor] 接收到文件事件: 类型=%s, 路径=%s, 目录=%s", event.Type, event.Path, event.Directory)

	if _, err := os.Stat(event.Path); err != nil {
		m.logger.Info("[Monitor] 文件未找到，跳过处理: %s", event.Path)
		return
	}

	m.handleDirectoryAggregation(event)
}

func (m *Monitor) handleDirectoryAggregation(event model.FileEvent) {
	m.dirMu.Lock()
	defer m.dirMu.Unlock()

	dir := filepath.Dir(event.Path)

	if _, exists := m.dirBuffers[dir]; !exists {
		m.dirBuffers[dir] = make(map[string]model.FileEvent)
		m.logger.Info("[Monitor] 初始化目录缓冲区: %s", dir)
	}

	m.dirBuffers[dir][event.Path] = event

	m.logger.Info("[Monitor] 目录事件聚合 - 目录: %s, 新增文件: %s, 缓冲区文件数: %d",
		dir, event.Path, len(m.dirBuffers[dir]))

	if timer, exists := m.dirTimers[dir]; exists {
		timer.Stop()
		m.logger.Info("[Monitor] 重置目录定时器: %s", dir)
	}

	quietMs := time.Duration(m.config.Settings.DirectoryStabilityQuietMs) * time.Millisecond
	m.dirTimers[dir] = time.AfterFunc(quietMs, func() {
		m.logger.Info("[Monitor] 目录静默期结束，开始处理目录事件: %s (静默期: %v)", dir, quietMs)
		m.processDirectoryEvents(dir)
	})

	m.logger.Info("[Monitor] 设置目录稳定性定时器: 目录=%s, 静默期=%v, 触发文件=%s",
		dir, quietMs, event.Path)

	timeoutMs := time.Duration(m.config.Settings.DirectoryStabilityTimeoutSeconds) * time.Second
	if timeoutMs > 0 && quietMs < timeoutMs {
		time.AfterFunc(timeoutMs, func() {
			m.dirMu.Lock()
			defer m.dirMu.Unlock()

			if events, exists := m.dirBuffers[dir]; exists && len(events) > 0 {
				m.logger.Warn("[Monitor] 目录稳定性检测超时，强制处理: %s (文件数量: %d, 超时: %v)",
					dir, len(events), timeoutMs)
				m.processDirectoryEvents(dir)
			}
		})
	}
}

func (m *Monitor) processDirectoryEvents(dir string) {
	m.dirMu.Lock()
	defer m.dirMu.Unlock()

	events := m.dirBuffers[dir]
	if len(events) == 0 {
		return
	}

	m.logger.Info("[Monitor] 目录已稳定，开始处理目录中的文件事件: %s (文件数量: %d)", dir, len(events))
	
	// 记录监控目录触发详情
	fileList := make([]string, 0, len(events))
	for filePath := range events {
		fileList = append(fileList, filepath.Base(filePath))
	}
	m.logger.Info("[Monitor] 监控目录触发详情 - 目录: %s, 文件列表: %v, 事件总数: %d", dir, fileList, len(events))

	delete(m.dirBuffers, dir)
	delete(m.dirTimers, dir)

	// 收集匹配的监控项，同一监控项只执行一次
	matchedMonitors := make(map[string]config.Monitor)
	var firstMatchingEvent model.FileEvent
	firstEventSet := false

	// 遍历所有事件，收集匹配的监控项
	for _, event := range events {
		for _, monitor := range m.config.Monitors {
			if !monitor.Enabled {
				continue
			}

			if event.Directory != monitor.Directory {
				continue
			}

			if !m.matchesFilePattern(event.Path, monitor.FilePatterns) {
				continue
			}

			if monitor.Schedule != "" && !m.isScheduleActive(monitor.Schedule) {
				continue
			}

			// 记录第一个匹配的事件，用于执行命令
			if !firstEventSet {
				firstMatchingEvent = event
				firstEventSet = true
			}

			// 记录匹配的监控项，使用命令作为唯一标识
			matchedMonitors[monitor.Command] = monitor
		}
	}

	// 对每个匹配的监控项执行一次命令
	for _, monitor := range matchedMonitors {
		m.logger.Info("[Monitor] 批量处理目录事件，执行命令: %s, 目录: %s, 文件数量: %d", monitor.Command, dir, len(events))
		m.executeCommand(monitor, firstMatchingEvent)
	}
}



func (m *Monitor) processFileEventInStableDir(event model.FileEvent) {
	m.logger.Info("[Monitor] 目录已稳定，开始处理文件事件: %s - %s", event.Type, event.Path)

	for _, monitor := range m.config.Monitors {
		if !monitor.Enabled {
			m.logger.Info("[Monitor] 监控项已禁用，跳过: %s", monitor.Directory)
			continue
		}

		if event.Directory != monitor.Directory {
			m.logger.Info("[Monitor] 目录不匹配，跳过: 监控目录=%s, 事件目录=%s", monitor.Directory, event.Directory)
			continue
		}

		if !m.matchesFilePattern(event.Path, monitor.FilePatterns) {
			m.logger.Info("[Monitor] 文件模式不匹配，跳过: 文件=%s, 期望模式=%v", event.Path, monitor.FilePatterns)
			continue
		}

		if monitor.Schedule != "" {
			m.logger.Info("[Monitor] 检查调度: 监控项=%s, 调度表达式=%s, 文件=%s",
				monitor.Command, monitor.Schedule, event.Path)
			if !m.isScheduleActive(monitor.Schedule) {
				m.logger.Info("[Monitor] 调度不匹配，跳过监控项: 命令=%s, 调度=%s, 文件=%s",
					monitor.Command, monitor.Schedule, event.Path)
				continue
			} else {
				m.logger.Info("[Monitor] 调度匹配，继续执行: 命令=%s, 调度=%s, 文件=%s",
					monitor.Command, monitor.Schedule, event.Path)
			}
		} else {
			m.logger.Info("[Monitor] 无调度限制，继续执行: 命令=%s, 文件=%s",
				monitor.Command, event.Path)
		}

		m.logger.Info("[Monitor] 找到匹配的监控项，准备执行命令: 监控名称=%s, 命令=%s, 文件=%s", monitor.Name, monitor.Command, event.Path)
		m.executeCommand(monitor, event)
	}
}

func (m *Monitor) matchesFilePattern(filePath string, patterns []string) bool {
	fileName := filepath.Base(filePath)
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}
	return false
}

func (m *Monitor) isScheduleActive(schedule string) bool {
	if schedule == "" {
		m.logger.Debug("[Monitor] 调度检查: 无调度表达式，总是激活")
		return true
	}

	gx := gronx.New()
	now := time.Now().Truncate(time.Minute)
	due, err := gx.IsDue(schedule, now)
	if err != nil {
		m.logger.Error("[Monitor] 调度表达式解析错误: %v, 表达式: %s", err, schedule)
		return false
	}

	m.logger.Debug("[Monitor] 调度检查详情: 表达式=%s, 当前时间=%v, 星期=%d, 小时=%d, 分钟=%d, 是否匹配=%v",
		schedule, now, now.Weekday(), now.Hour(), now.Minute(), due)

	if !due {
		m.logger.Info("[Monitor] 调度不匹配，跳过执行: 表达式=%s, 当前时间=%v", schedule, now)
	}

	return due
}

func (m *Monitor) isFileStable(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	if info.Size() == 0 {
		m.logger.Debug("[Monitor] 文件大小为0: %s", filePath)
	}

	minStabilityTime := time.Duration(m.config.Settings.MinStabilityTimeMs) * time.Millisecond
	fileAge := time.Since(info.ModTime())

	if fileAge < 0 {
		m.logger.Debug("[Monitor] 文件修改时间在未来，但仍处理: %s", filePath)
		return true
	}

	isStable := fileAge >= minStabilityTime
	if !isStable {
		m.logger.Debug("[Monitor] 文件尚未稳定: %s (age: %v, required: %v)", filePath, fileAge, minStabilityTime)
	}
	return isStable
}

func (m *Monitor) executeCommand(monitor config.Monitor, event model.FileEvent) {
	m.logger.Info("[Monitor] 开始执行命令: %s", monitor.Command)
	m.logger.Info("[Monitor] 命令执行详情 - 监控名称: %s, 目录: %s, 文件: %s, 事件类型: %s", 
		monitor.Name, monitor.Directory, event.Path, event.Type)

	if m.isDuplicate(monitor.Command, event.Path) {
		m.logger.Info("[Monitor] 检测到重复执行，跳过: 命令=%s, 文件=%s", monitor.Command, event.Path)
		return
	}

	executor := NewCommandExecutor(m.logger, monitor.Directory)

	executor.SetEnvVar("FILE_PATH", event.Path)
	executor.SetEnvVar("FILE_NAME", filepath.Base(event.Path))
	executor.SetEnvVar("FILE_DIR", filepath.Dir(event.Path))
	executor.SetEnvVar("EVENT_TYPE", string(event.Type))

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.logger.Info("[Monitor] 启动命令执行goroutine: 命令=%s", monitor.Command)

		select {
		case m.opSem <- struct{}{}:
			defer func() { <-m.opSem }()
			m.logger.Info("[Monitor] 获取到操作信号量: 命令=%s", monitor.Command)
		case <-m.opCtx.Done():
			m.logger.Info("[Monitor] 操作被取消，无法获取信号量: 命令=%s", monitor.Command)
			return
		}

		m.logger.Info("[Monitor] 开始执行命令: %s (超时: %d秒)", monitor.Command, monitor.Timeout)
		if err := executor.ExecuteCommandWithContext(m.opCtx, monitor.Command, &event, monitor.Timeout); err != nil {
			m.logger.Error("[Monitor] 命令执行失败: %v", err)
			return
		}

		m.logger.Info("[Monitor] 命令执行成功: %s", monitor.Command)
	}()
}

func (m *Monitor) isDuplicate(command, filePath string) bool {
	key := command + "|" + filePath

	m.dedupMu.Lock()
	defer m.dedupMu.Unlock()

	now := time.Now()
	if lastExec, exists := m.dedupCache[key]; exists {
		if now.Sub(lastExec) < time.Duration(m.config.Settings.ExecutionDedupIntervalSeconds)*time.Second {
			return true
		}
	}

	m.dedupCache[key] = now
	return false
}

func (m *Monitor) shouldLogDrop(filePath string) bool {
	m.dropMu.Lock()
	defer m.dropMu.Unlock()

	now := time.Now()
	if lastLog, exists := m.dropLog[filePath]; exists {
		if now.Sub(lastLog) < DropLogThrottleTime {
			return false
		}
	}

	m.dropLog[filePath] = now
	return true
}

func (m *Monitor) cleanupDaemon() {
	defer m.wg.Done()

	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-m.cleanupStop:
			return
		case <-ticker.C:
			m.cleanup()
		}
	}
}

func (m *Monitor) cleanup() {
	m.dedupMu.Lock()
	now := time.Now()
	for key, lastExec := range m.dedupCache {
		if now.Sub(lastExec) > DedupCacheExpiration {
			delete(m.dedupCache, key)
		}
	}
	m.dedupMu.Unlock()

	m.dropMu.Lock()
	for filePath, lastLog := range m.dropLog {
		if now.Sub(lastLog) > DedupCacheExpiration {
			delete(m.dropLog, filePath)
		}
	}
	m.dropMu.Unlock()
}
