package dump

// Dumper is expected to export "something" to a file and return a complete `DumpResult` struct (`Path`, `MIME`)
type Dumper interface {
	Dump() (*DumpResult, error)
}
