package domain

import "fmt"

// store last tailSize bytes
func NewTailStorage(tailSize int) *tailStorage {
	if tailSize < 0 {
		panic("tail size is negative")
	}

	return &tailStorage{
		tailSize: tailSize,
		tail:     make([]byte, tailSize),
		tailLen:  0,
	}
}

type tailStorage struct {
	tailSize int
	tailLen  int

	tail []byte
}

func (ts *tailStorage) Write(chunk []byte) []byte {

	if len(chunk) >= ts.tailSize {
		result := make([]byte, len(chunk)+ts.tailLen-ts.tailSize)
		copy(result, ts.tail[:ts.tailLen])
		copy(result[ts.tailLen:], chunk[:len(chunk)-ts.tailLen])
		copy(ts.tail, chunk[len(chunk)-ts.tailSize:])
		ts.tailLen = ts.tailSize
		return result
	}
	if len(chunk)+ts.tailLen >= ts.tailSize {
		resultSize := len(chunk) + ts.tailLen - ts.tailSize
		result := make([]byte, resultSize)
		copy(result, ts.tail[:resultSize])

		copy(ts.tail, ts.tail[resultSize:])
		copy(ts.tail[ts.tailSize-len(chunk):], chunk)
		ts.tailLen = ts.tailSize
		return result
	} else {
		copy(ts.tail[ts.tailLen:], chunk)
		ts.tailLen += len(chunk)
		return nil
	}

}

func (ts *tailStorage) Finish() ([]byte, error) {
	if ts.tailLen < ts.tailSize {
		return nil, fmt.Errorf("tail is not complete")
	}
	return ts.tail, nil
}
