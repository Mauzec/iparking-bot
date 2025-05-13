package data

type DistanceProvider interface {
	ReadDistance() (int, error)
}

type FileProvider struct {
	Path string
}

func (f *FileProvider) ReadDistance() (int, error) {
	return ReadDistance(f.Path)
}

func (f *FileProvider) NewFileProvider(path string) *FileProvider {
	return &FileProvider{
		Path: path,
	}
}
