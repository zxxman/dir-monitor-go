package monitor

import (
	"dir-monitor-go/internal/model"
)

// Watcher 监视器接口
type Watcher interface {
	// Watch 开始监控指定目录
	Watch(dir string) error
	// Unwatch 停止监控指定目录
	Unwatch(dir string) error
	// Events 返回事件通道
	Events() <-chan model.FileEvent
	// Close 关闭监视器
	Close() error
}
