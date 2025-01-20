package zapslog

import (
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

func Init(debug bool) (*slog.Logger, func()) {
	cfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	level := zap.InfoLevel
	fileName := "logs/app.log"
	if debug {
		level = zap.DebugLevel
		fileName = "logs/debug.log"
	}

	fileEncoder := zapcore.NewJSONEncoder(cfg)
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10, // Mb
		MaxAge:     7,  // days
		MaxBackups: 3,
	})

	core := zapcore.NewCore(fileEncoder, file, level)
	if debug {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, file, level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
	}

	logger := zap.New(core)

	sl := slog.New(zapslog.NewHandler(logger.Core()))

	return sl, func() {
		logger.Sync()
	}
}
