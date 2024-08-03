package domain

import "context"

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . Logger,StreamFileWriter,StreamFileReader

type Logger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
}

// StreamFileWriter used for sending big files in chunk
type StreamFileWriter interface {
	WriteChunk(ctx context.Context, name string, chunk []byte) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// StreamFileReader is used for reading file chunk
type StreamFileReader interface {
	FileSize() int64
	Next() ([]byte, error)
	Close()
}
