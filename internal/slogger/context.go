package slogger

import (
	"context"
	"log/slog"
)

var disabledLogger *slog.Logger

type disabledHandler struct{}

func (disabledHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (disabledHandler) Handle(context.Context, slog.Record) error { return nil }
func (h disabledHandler) WithAttrs([]slog.Attr) slog.Handler      { return h }
func (h disabledHandler) WithGroup(string) slog.Handler           { return h }

func init() {
	disabledLogger = slog.New(disabledHandler{})
}

type contextKey struct{}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

func With(ctx context.Context, args ...any) context.Context {
	return WithContext(ctx, Ctx(ctx).With(args...))
}
func Ctx(ctx context.Context) *slog.Logger {
	l := ctx.Value(contextKey{})
	if l == nil {
		return disabledLogger
	}
	return l.(*slog.Logger)
}
