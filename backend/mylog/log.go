package mylog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"sdim_pc/backend/config"
)

var (
	rootLg *zerolog.Logger
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
}

func InitLogger(cfg config.LogConfig, profile string) error {
	ll, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(ll)

	lg := zerolog.New(zerolog.MultiLevelWriter(
		createStdoutWriter(cfg.ColorfulStd),
		createFileWriter(cfg),
	)).With().
		Timestamp().
		Caller().
		Str("logger", "appLogger").
		Int("pid", os.Getpid()).Str("profile", profile).
		Logger()

	rootLg = &lg

	return nil
}

func GetLogger() *zerolog.Logger {
	return rootLg
}

func createStdoutWriter(colorfulStdout bool) io.Writer {
	if colorfulStdout {
		return zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    false,
			TimeFormat: "2006-01-02 15:04:05.000",
		}
	} else {
		return os.Stdout
	}
}

func createFileWriter(cfg config.LogConfig) io.Writer {
	logNamePaths := []string{
		cfg.FilePath,
		"point.log",
	}

	return &lumberjack.Logger{
		Filename:   filepath.Join(logNamePaths...),
		MaxSize:    cfg.MaxFileSize,
		MaxAge:     cfg.HistoryDays,
		MaxBackups: cfg.MaxBackup,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}
}
