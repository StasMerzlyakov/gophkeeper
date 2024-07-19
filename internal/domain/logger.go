package domain

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var mainLogger *zap.SugaredLogger

func GetCtxLogger(ctx context.Context) Logger {
	if v := ctx.Value(LoggerKey); v != nil {
		lg, ok := v.(Logger)
		if !ok {
			return GetApplicationLogger()
		}
		return lg
	}
	return GetApplicationLogger()
}

func GetApplicationLogger() *zap.SugaredLogger {

	if mainLogger != nil {
		return mainLogger
	} else {
		log.Default().Println("[WARN] application logger is not set")
		logger := zap.NewNop()
		mainLogger = logger.Sugar()
	}

	return mainLogger
}

func SetApplicationLogger(logger *zap.SugaredLogger) {
	mainLogger = logger
}

type ContextKey string

const UserIDKey = ContextKey("UserID")

const LoggerKey = ContextKey("Logger")
const LoggerKeyRequestID = "requestID"

const LoggerKeyUserID = "userID"

func EnrichWithUserID(ctx context.Context, userID UserID) context.Context {
	resultCtx := context.WithValue(ctx, UserIDKey, userID)
	return resultCtx
}

func GetUserID(ctx context.Context) (UserID, error) {
	if v := ctx.Value(UserIDKey); v != nil {
		requestID, ok := v.(UserID)
		if !ok {
			return UserID(""), fmt.Errorf("%w: unexpected userID type", ErrServerInternal)
		}
		return requestID, nil
	}
	return UserID(""), fmt.Errorf("%w: can't extract userID", ErrNotAuthorized)
}

func EnrichWithRequestIDLogger(ctx context.Context, requestID uuid.UUID, logger Logger) context.Context {
	requestIDLogger := &requestIDLogger{
		internalLogger: logger,
		requestID:      requestID.String(),
	}
	resultCtx := context.WithValue(ctx, LoggerKey, requestIDLogger)
	return resultCtx
}

var _ Logger = (*requestIDLogger)(nil)

type requestIDLogger struct {
	requestID      string
	internalLogger Logger
}

func (l *requestIDLogger) Debugw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, LoggerKeyRequestID, l.requestID)
	l.internalLogger.Debugw(msg, keysAndValues...)
}

func (l *requestIDLogger) Infow(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, LoggerKeyRequestID, l.requestID)
	l.internalLogger.Infow(msg, keysAndValues...)
}

func (l *requestIDLogger) Errorw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, LoggerKeyRequestID, l.requestID)
	l.internalLogger.Infow(msg, keysAndValues...)
}

func EnrichWithUserIDLogger(ctx context.Context, userID UserID, logger Logger) context.Context {
	requestIDLogger := &userIDLogger{
		internalLogger: logger,
		userID:         string(userID),
	}
	resultCtx := context.WithValue(ctx, LoggerKey, requestIDLogger)
	return resultCtx
}

type userIDLogger struct {
	userID         string
	internalLogger Logger
}

func (l *userIDLogger) Debugw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, LoggerKeyUserID, l.userID)
	l.internalLogger.Debugw(msg, keysAndValues...)
}

func (l *userIDLogger) Infow(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, LoggerKeyUserID, l.userID)
	l.internalLogger.Infow(msg, keysAndValues...)
}

func (l *userIDLogger) Errorw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, LoggerKeyUserID, l.userID)
	l.internalLogger.Infow(msg, keysAndValues...)
}
