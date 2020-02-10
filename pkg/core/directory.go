package core

type Directory struct {
	Path        string
	Recursive   bool
	IgnoreFiles []string
	Target      []string
}

type DirectoryService interface{}

type DirectoryStore interface{}
