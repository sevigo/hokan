package core

type Watcher interface {
	StartDirWatcher()
	StartFileWatcher()
	GetDirsToWatch() error
}
