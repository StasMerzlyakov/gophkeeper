package domain

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . Logger,StreamSender,StreamFileReader

type Logger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
}

// StreamSender used for sending big files in chunk
type StreamSender interface {
	Send(chunk []byte) error
	CloseAndRecv() error
}

// StreamFileReader is used for reading file chunk
type StreamFileReader interface {
	FileSize() int64
	Next() ([]byte, error)
	Close()
}
