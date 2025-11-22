package model

import "time"

type FileEventType string

const (
	FileCreated  FileEventType = "created"
	FileModified FileEventType = "modified"
	FileDeleted  FileEventType = "deleted"
	FileRenamed  FileEventType = "renamed"
)

type FileEvent struct {
	Type      FileEventType
	Path      string
	OldPath   string
	Timestamp time.Time
	Size      int64
	ModTime   time.Time
	Directory string
}
