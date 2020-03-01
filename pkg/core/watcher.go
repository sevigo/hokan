package core

type Watcher interface {
	StartDirWatcher()
	GetDirsToWatch() error
}
