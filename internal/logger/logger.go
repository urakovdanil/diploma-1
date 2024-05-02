package logger

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/types"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var l *slog.Logger

func New(cfg *config.StartUp) {
	l = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource:   true,
				Level:       logLevelsMap[cfg.GetLogLevel()],
				ReplaceAttr: replaceAttr,
			}),
	)
}

func Debug(ctx context.Context, message string) {
	log(ctx, message, slog.LevelDebug)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	log(ctx, fmt.Sprintf(format, args...), slog.LevelDebug)
}

func Info(ctx context.Context, message string) {
	log(ctx, message, slog.LevelInfo)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	log(ctx, fmt.Sprintf(format, args...), slog.LevelInfo)
}

func Warn(ctx context.Context, message string) {
	log(ctx, message, slog.LevelWarn)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	log(ctx, fmt.Sprintf(format, args...), slog.LevelWarn)
}

func Error(ctx context.Context, message string) {
	log(ctx, message, slog.LevelError)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log(ctx, fmt.Sprintf(format, args...), slog.LevelError)
}

func Fatal(ctx context.Context, message string) {
	log(ctx, "FATAL: "+message, slog.LevelError+1)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	log(ctx, fmt.Sprintf("FATAL: "+format, args...), slog.LevelError+1)
}

func log(ctx context.Context, message string, level slog.Level) {
	if !l.Enabled(context.Background(), level) {
		return
	}
	curLogger := l
	internalRequestID := ctx.Value(types.CtxKeyRequestID)
	if internalRequestID != nil {
		curLogger = l.With(slog.String(string(types.CtxKeyRequestID), internalRequestID.(string)))
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, message, pcs[0])
	_ = curLogger.Handler().Handle(ctx, r)
	if level == slog.LevelError+1 {
		os.Exit(1)
	}
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
	}
	return a
}

var logLevelsMap = map[string]slog.Level{
	slog.LevelDebug.String(): slog.LevelDebug,
	slog.LevelInfo.String():  slog.LevelInfo,
	slog.LevelWarn.String():  slog.LevelWarn,
	slog.LevelError.String(): slog.LevelError,
}
