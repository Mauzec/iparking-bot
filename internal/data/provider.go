package data

type DistanceProvider interface {
	ReadDistance() (float64, error)
}

type FileProvider struct {
	Path string
}

func (f *FileProvider) ReadDistance() (float64, error) {
	return ReadDistance(f.Path)
}

func (f *FileProvider) NewFileProvider(path string) *FileProvider {
	return &FileProvider{
		Path: path,
	}
}
