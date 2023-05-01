package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gerladeno/chat-service/internal/buildinfo"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
	sentryDSN      string
	env            string
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
		return fmt.Errorf("parsing log level: %w", err)
	}
	config := zap.NewProductionEncoderConfig()
	config.NameKey = "component"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "T"
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(config)
	if opts.productionMode {
		Atom.SetLevel(zap.InfoLevel)
		config.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(config)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, Atom),
	}
	if opts.sentryDSN != "" {
		cfg := zapsentry.Configuration{
			Level: zapcore.WarnLevel,
		}
		client, err := NewSentryClient(opts.sentryDSN, opts.env, buildinfo.GetSentryVersion())
		if err != nil {
			return fmt.Errorf("initing sentry client: %w", err)
		}
		core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))
		if err != nil {
			return fmt.Errorf("adding sentry to zap: %w", err)
		}
		cores = append(cores, core)
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
