package monitor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"dir-monitor-go/internal/logger"
	"dir-monitor-go/internal/model"

	"github.com/fsnotify/fsnotify"
)

type renamePair struct {
	old string
	ts  time.Time
}

const (
	// 移动对窗口时间
	movePairWindow = 500 * time.Millisecond

	// 默认事件通道缓冲区大小
	DefaultEventChannelBuffer = 100
)

type handlerEntry struct {
	baseDir string
	fn      func(model.FileEvent)
}

// FsnotifyWatcher File system monitor based on native fsnotify library (event-driven implementation)
type FsnotifyWatcher struct {
	logger      *logger.Logger
	fsWatcher   *fsnotify.Watcher // Use actual fsnotify watcher
	watchedDirs map[string]bool
	mu          sync.RWMutex
	events      chan model.FileEvent
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	handlers    []handlerEntry

	// recent REMOVE cache for pairing with CREATE -> renamed
	movePairs map[string]renamePair // key: directory path

	// stop flag to prevent any further event logging/dispatch after Stop
	stopping int32
}

// NewFsnotifyWatcher Create a new monitor based on native fsnotify (event-driven implementation)
func NewFsnotifyWatcher(logger *logger.Logger) *FsnotifyWatcher {
	ctx, cancel := context.WithCancel(context.Background())

	// Create actual fsnotify watcher
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("[FsnotifyWatcher] Failed to create fsnotify watcher: %v", err)
		cancel()
		return nil
	}

	// Use default debounce time
	// Use default event channel buffer size
	eventChannelBuffer := DefaultEventChannelBuffer

	return &FsnotifyWatcher{
		logger:      logger,
		fsWatcher:   fsWatcher,
		watchedDirs: make(map[string]bool),
		events:      make(chan model.FileEvent, eventChannelBuffer),
		ctx:         ctx,
		cancel:      cancel,
		handlers:    make([]handlerEntry, 0),
		movePairs:   make(map[string]renamePair),
	}
}

// Watch starts monitoring specified directory
func (fw *FsnotifyWatcher) Watch(baseDir string, handler func(model.FileEvent)) error {
	fw.mu.Lock()
	fw.handlers = append(fw.handlers, handlerEntry{baseDir: baseDir, fn: handler})
	fw.mu.Unlock()

	// Set up directory watch (non-recursive)
	if err := fw.setupWatch(baseDir); err != nil {
		return fmt.Errorf("failed to setup watch for %s: %w", baseDir, err)
	}

	fw.wg.Add(1)
	go fw.processEvents()

	return nil
}

// Stop stops monitoring
func (fw *FsnotifyWatcher) Stop() error {
	// mark stopping first to avoid race logging in handlers
	atomic.StoreInt32(&fw.stopping, 1)
	fw.cancel()
	if fw.fsWatcher != nil {
		_ = fw.fsWatcher.Close()
	}
	fw.wg.Wait()
	close(fw.events)
	return nil
}







func (fw *FsnotifyWatcher) removeWatchRecursive(path string) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if _, exists := fw.watchedDirs[path]; exists {
		delete(fw.watchedDirs, path)
		if err := fw.fsWatcher.Remove(path); err != nil {
			fw.logger.Error("[FsnotifyWatcher] Failed to remove watch for path: %s, error: %v", path, err)
		}
	}

	for watchedPath := range fw.watchedDirs {
		if strings.HasPrefix(watchedPath, path+string(filepath.Separator)) {
			delete(fw.watchedDirs, watchedPath)
			if err := fw.fsWatcher.Remove(watchedPath); err != nil {
				fw.logger.Error("[FsnotifyWatcher] Failed to remove watch for sub-path: %s, error: %v", watchedPath, err)
			}
		}
	}
}

func (fw *FsnotifyWatcher) setupWatch(baseDir string) error {
	info, err := os.Stat(baseDir)
	if err != nil {
		return fmt.Errorf("stat failed for %s: %w", baseDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", baseDir)
	}

	err = filepath.Walk(baseDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			fw.mu.RLock()
			_, alreadyWatched := fw.watchedDirs[p]
			fw.mu.RUnlock()
			if alreadyWatched {
				return nil
			}
			if err := fw.fsWatcher.Add(p); err != nil {
				fw.logger.Warn("[FsnotifyWatcher] Failed to add watch for directory %s: %v", p, err)
				return err
			}
			fw.mu.Lock()
			fw.watchedDirs[p] = true
			fw.mu.Unlock()
		}
		return nil
	})
	return err
}



func (fw *FsnotifyWatcher) processEvents() {
	defer fw.wg.Done()

	for {
		select {
		case <-fw.ctx.Done():
			return
		case event, ok := <-fw.fsWatcher.Events:
			if !ok {
				return
			}
			fw.handleFsnotifyEvent(event)

		case err, ok := <-fw.fsWatcher.Errors:
			if !ok {
				return
			}
			fw.logger.Error("[FsnotifyWatcher] Monitor error: %v", err)
		}
	}
}

func (fw *FsnotifyWatcher) handleFsnotifyEvent(event fsnotify.Event) {
	if atomic.LoadInt32(&fw.stopping) == 1 {
		return
	}

	if fw.shouldIgnoreFile(event.Name) {
		return
	}

	var eventType model.FileEventType
	var oldPath string
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		dir := filepath.Dir(event.Name)
		paired := false
		fw.mu.Lock()
		if p, ok := fw.movePairs[dir]; ok && time.Since(p.ts) <= movePairWindow {
			paired = true
			oldPath = p.old
			delete(fw.movePairs, dir)
		}
		fw.mu.Unlock()
		if paired {
			eventType = model.FileRenamed
			fw.logger.Info("[FsnotifyWatcher] 检测到重命名事件(配对): %s -> %s", oldPath, event.Name)
		} else {
			eventType = model.FileCreated
			fw.logger.Info("[FsnotifyWatcher] 检测到文件创建事件: %s", event.Name)
		}
	case event.Op&fsnotify.Write == fsnotify.Write:
		eventType = model.FileModified
		fw.logger.Info("[FsnotifyWatcher] 检测到文件修改事件: %s", event.Name)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		eventType = model.FileDeleted
		fw.logger.Info("[FsnotifyWatcher] 检测到文件删除事件: %s", event.Name)
		dir := filepath.Dir(event.Name)
		fw.mu.Lock()
		fw.movePairs[dir] = renamePair{old: event.Name, ts: time.Now()}
		fw.mu.Unlock()
		fw.removeWatchRecursive(event.Name)
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		eventType = model.FileRenamed
		fw.logger.Info("[FsnotifyWatcher] 检测到重命名事件: %s", event.Name)
		fw.removeWatchRecursive(event.Name)
	default:
		fw.logger.Info("[FsnotifyWatcher] 未识别的事件类型，忽略: %s", event.Op)
		return
	}

	fileEvent := model.FileEvent{
		Type:      eventType,
		Path:      event.Name,
		Directory: filepath.Dir(event.Name),
		Timestamp: time.Now(),
	}
	if eventType == model.FileRenamed && oldPath != "" {
		fileEvent.OldPath = oldPath
	}

	select {
	case fw.events <- fileEvent:
		fw.logger.Info("[FsnotifyWatcher] 文件事件已发送到事件通道: %s (类型: %s)", event.Name, eventType)
	default:
		fw.logger.Info("[FsnotifyWatcher] 事件通道已满，丢弃事件: %s (类型: %s)", event.Name, eventType)
	}

	fw.dispatchEvent(fileEvent)
}

func (fw *FsnotifyWatcher) shouldIgnoreFile(path string) bool {
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") || strings.HasSuffix(base, "~") || strings.HasSuffix(base, ".tmp") {
		return true
	}
	if strings.HasSuffix(base, ".swp") || strings.HasSuffix(base, ".swo") || strings.HasSuffix(base, ".swn") {
		return true
	}
	if strings.HasSuffix(base, ".lock") || strings.HasSuffix(base, ".bak") {
		return true
	}
	return false
}

// dispatchEvent Dispatch event to all handlers
func (fw *FsnotifyWatcher) dispatchEvent(event model.FileEvent) {
	fw.mu.RLock()
	handlers := make([]handlerEntry, len(fw.handlers))
	copy(handlers, fw.handlers)
	fw.mu.RUnlock()

	for _, h := range handlers {
		if strings.HasPrefix(event.Path, h.baseDir) {
			h.fn(event)
		}
	}
}
