package models

type FileInfo struct {
	Name     string
	FilePath string
}

type File struct {
	FileInfo
	Key  string
	Data []byte
}
