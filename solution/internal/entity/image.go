package entity

type Image struct {
	Name        string
	Buffer      []byte
	Size        int64
	ContentType string
}
