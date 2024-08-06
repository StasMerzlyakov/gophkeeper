package domain

import "os"

func NewStreamFileReader(path string, fileChunkSize int) (*streamFileReader, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	inf, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &streamFileReader{
		file:          file,
		size:          inf.Size(),
		fileChunkSize: fileChunkSize,
	}, nil
}

var _ StreamFileReader = (*streamFileReader)(nil)

type streamFileReader struct {
	file          *os.File
	size          int64
	fileChunkSize int
}

func (sf *streamFileReader) FileSize() int64 {
	return sf.size
}

func (sf *streamFileReader) Next() ([]byte, error) {
	buf := make([]byte, sf.fileChunkSize)
	n, err := sf.file.Read(buf)
	if err != nil {
		return buf[:n], err
	}
	return buf[:n], nil
}
func (sf *streamFileReader) Close() {
	err := sf.file.Close()
	if err != nil {
		panic(err)
	}
}
