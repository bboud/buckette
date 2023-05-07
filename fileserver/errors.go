package fileserver

type FileExists struct{}

func (e *FileExists) Error() string {
	return "file already exists"
}

type MaxRecordsReached struct{}

func (e *MaxRecordsReached) Error() string {
	return "a max number of records that the server can cache has been reached"
}
