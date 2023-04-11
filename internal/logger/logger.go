package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

var Atom zap.AtomicLevel

func Init(opts Options) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf("validate options: %v", err)
	}
	Atom, err = zap.ParseAtomicLevel(opts.level)
	if err != nil {
		return fmt.Errorf("err parsing log level: %w", err)
	}
	config := zap.NewProductionEncoderConfig()
	config.NameKey = "component"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "T"
	encoder := zapcore.NewConsoleEncoder(config)
	if opts.productionMode {
		Atom = zap.NewAtomicLevelAt(zap.InfoLevel)
		encoder = zapcore.NewJSONEncoder(config)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, Atom),
	}
	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
