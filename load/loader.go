package load

// Dumper is expected to export "something" to a file and return a complete `DumpResult` struct (`Path`, `MIME`)
type Loader interface {
	Load() error
	Update() error
}
