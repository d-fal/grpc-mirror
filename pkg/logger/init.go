package logger

import (
	"go.uber.org/zap"
)

var (
	zapLoggersMap map[string]*zap.Logger
)

func init() {
	// Log as JSON instead of the default ASCII formatter.

	zapLoggersMap = make(map[string]*zap.Logger)
}
