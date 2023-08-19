package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// logがどんな感じに出力されるかを確認する用
func TestInitializeLogger(t *testing.T) {
	l := InitializeLogger()

	assert.IsType(t, &zap.Logger{}, l)

	// 例示のために、いくつかのログを出力します。
	l.Debug("test debug message", zap.String("foo", "bar"))
	l.Info("test info message", zap.String("foo", "bar"))
	l.Warn("test warning message", zap.Int("baz", 42))
	l.Error("test error message", zap.Int("baz", 42))
	// l.DPanic("test dpanic message", zap.Int("baz", 42))
}
