package filesystem

import "time"

// The interface for the filesystem that must be implemented
type FS interface {
	Put(fileName string, folder string, acl ...string) error
	Get(destination string, items ...string) error
	List(prefix string) ([]Listing, error)
	Delete(itemsToDelete []string) bool
}

// Describes one file on a remote file system
type Listing struct {
	Etag         string
	LastModified time.Time
	Key          string
	Size         float64
	IsDir        bool
}
