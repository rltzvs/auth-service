package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(mode string) (*zap.Logger, error) {
	var config zap.Config

	switch mode {
	case "dev":
		config = zap.NewDevelopmentConfig()

		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case "prod":
		config = zap.NewProductionConfig()
	default:
		return nil, fmt.Errorf("unknown logger mode")
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()

	if err != nil {
		return nil, fmt.Errorf("logger build err: %w", err)
	}

	return logger, nil
}
