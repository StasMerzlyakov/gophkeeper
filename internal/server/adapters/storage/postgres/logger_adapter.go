package postgres

import (
	"context"
	"log"
	"sort"

	"github.com/jackc/pgx/v5/tracelog"
)

func NewLogAdapter(logger Logger) *loggerAdapter {
	return &loggerAdapter{
		logger: logger,
	}
}

type loggerAdapter struct {
	logger Logger
}

func (la *loggerAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {

	pl := make(pairList, len(data))
	i := 0
	for k, v := range data {
		pl[i] = pair{k, v}
		i++
	}
	sort.Sort(pl)

	keyAndValues := make([]any, 2*len(data))

	for i, pair := range pl {
		keyAndValues[2*i] = pair.Key
		keyAndValues[2*i+1] = pair.Value
	}

	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		la.logger.Debugw(msg, keyAndValues...)
	case tracelog.LogLevelInfo, tracelog.LogLevelWarn:
		la.logger.Infow(msg, keyAndValues...)
	case tracelog.LogLevelError:
		la.logger.Errorw(msg, keyAndValues...)
	default:
		log.Printf("invalid level %d", la)
	}
}

type pair struct {
	Key   string
	Value any
}
type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Key < p[j].Key }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
