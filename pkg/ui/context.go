package ui

import (
	"context"
)

type uiCtxKey struct{}

func ToContext(ctx context.Context, ui *UI) context.Context {
	return context.WithValue(ctx, uiCtxKey{}, ui)
}

func FromContext(ctx context.Context) (*UI, bool) {
	if value := ctx.Value(uiCtxKey{}); value != nil {
		return value.(*UI), true
	}
	return nil, false
}

func Update(ctx context.Context, msg string) {
	if ui, ok := FromContext(ctx); ok {
		ui.Update(msg)
	}
}

func Info(ctx context.Context, msg string) {
	if ui, ok := FromContext(ctx); ok {
		ui.Info(msg)
	}
}

func Error(ctx context.Context, msg string) {
	if ui, ok := FromContext(ctx); ok {
		ui.Error(msg)
	}
}

func Stop(ctx context.Context) {
	if ui, ok := FromContext(ctx); ok {
		ui.Stop()
	}
}
