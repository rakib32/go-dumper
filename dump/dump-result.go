package dump

// DumpResult is the result of an export operation.
type DumpResult struct {
	// Path to exported file
	Path string

	TarFilePath string
}
