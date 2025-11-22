package monitor

import (
	"context"
	"sync"

	"dir-monitor-go/internal/logger"
)

type MonitorManager struct {
	monitors []*Monitor
	logger   *logger.Logger
	wg       sync.WaitGroup
	mu       sync.Mutex
}

func NewMonitorManager(log *logger.Logger) *MonitorManager {
	return &MonitorManager{
		monitors: make([]*Monitor, 0),
		logger:   log,
	}
}

func (mm *MonitorManager) AddMonitor(monitor *Monitor) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.monitors = append(mm.monitors, monitor)
}

func (mm *MonitorManager) Start(ctx context.Context) error {
	var startErr error

	mm.mu.Lock()
	snapshot := make([]*Monitor, len(mm.monitors))
	copy(snapshot, mm.monitors)
	mm.mu.Unlock()

	for _, monitor := range snapshot {
		if err := monitor.Start(); err != nil {
			mm.logger.Error("[MonitorManager] Failed to start monitor: %v", err)
			startErr = err
		}
	}

	if startErr != nil {
		return startErr
	}

	return nil
}

func (mm *MonitorManager) Stop(ctx context.Context) error {
	mm.mu.Lock()
	snapshot := make([]*Monitor, len(mm.monitors))
	copy(snapshot, mm.monitors)
	mm.mu.Unlock()

	for _, monitor := range snapshot {
		monitor.Stop()
	}

	mm.wg.Wait()
	mm.cleanupResources()

	return nil
}

// cleanupResources 清理资源
func (mm *MonitorManager) cleanupResources() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.monitors = nil
}

func (mm *MonitorManager) GetMonitorCount() int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	return len(mm.monitors)
}

func (mm *MonitorManager) Wait() {
	mm.wg.Wait()
}
