package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger

func Init(debug bool) error {
	cfg := zap.NewProductionConfig()
	if debug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	lg, err := cfg.Build()
	if err != nil {
		return err
	}
	sugar = lg.Sugar()
	return nil
}

func Logger() *zap.SugaredLogger {
	return sugar
}
