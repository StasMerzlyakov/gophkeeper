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
		buf:  make([]byte, fileChunkSize),
		file: file,
		size: inf.Size(),
	}, nil
}

var _ StreamFileReader = (*streamFileReader)(nil)

type streamFileReader struct {
	buf  []byte
	file *os.File
	size int64
}

func (sf *streamFileReader) FileSize() int64 {
	return sf.size
}

func (sf *streamFileReader) Next() ([]byte, error) {
	n, err := sf.file.Read(sf.buf)
	if err != nil {
		return sf.buf[:n], err
	}
	return sf.buf[:n], nil
}
func (sf *streamFileReader) Close() {
	err := sf.file.Close()
	if err != nil {
		panic(err)
	}
}
