package data

import (
	"context"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Emit wraps runtime.EventsEmit with nil guard and recover to support CLI/service usage.
func Emit(ctx context.Context, event string, data ...interface{}) {
	if os.Getenv("GO_STOCK_CLI") == "1" {
		return
	}
	if ctx == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	runtime.EventsEmit(ctx, event, data...)
}
