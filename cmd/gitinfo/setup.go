package main

import (
	"errors"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jessevdk/go-flags"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Version bool `long:"version"                       description:"Show version and exit"`
	Debug   bool `long:"debug"                         description:"Show debug data"`
}

var (
	// ErrGotHelp returned after showing requested help
	ErrGotHelp = errors.New("help printed")
	// ErrBadArgs returned after showing command args error message
	ErrBadArgs = errors.New("option error printed")
)

// SetupConfig loads flags from args (if given) or command flags and ENV otherwise
func SetupConfig(args ...string) (*Config, error) {
	cfg := &Config{}
	p := flags.NewParser(cfg, flags.Default|flags.PrintErrors) //  HelpFlag | PrintErrors | PassDoubleDash
	var err error
	if len(args) == 0 {
		_, err = p.Parse()
	} else {
		_, err = p.ParseArgs(args)
	}
	if err == nil {
		return cfg, nil
	}
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		return nil, ErrGotHelp
	}
	return nil, ErrBadArgs
}

// SetupLog creates logger
func SetupLog(withDebug bool, opts ...zap.Option) logr.Logger {
	var zapLog *zap.Logger
	if withDebug {
		aa := zap.NewDevelopmentEncoderConfig()
		zo := append(opts, zap.AddCaller())
		aa.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLog = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(aa),
			zapcore.AddSync(colorable.NewColorableStdout()),
			zapcore.DebugLevel,
		),
			zo...,
		)
	} else {
		zapLog, _ = zap.NewProduction(opts...)
	}
	return zapr.NewLogger(zapLog)
}

// Shutdown runs exit after deferred cleanups have run
func Shutdown(exitFunc func(code int), e error, log logr.Logger) {
	if e != nil {
		var code int
		switch e {
		case ErrGotHelp:
			code = 3
		case ErrBadArgs:
			code = 2
		default:
			log.Error(e, "Run error")
			code = 1
		}
		exitFunc(code)
	}
}
