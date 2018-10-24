package store

// Storer takes an `DumpResult` and move it somewhere! To a cloud storage service, for instance...
type Storer interface {
	Store(fileKey string) error
	Download(fileKey string) (string, error)
}
