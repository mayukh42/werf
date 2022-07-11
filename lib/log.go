package lib

import (
	"path/filepath"

	"github.com/mayukh42/logx/logx"
	"github.com/mayukh42/werf/config"
)

func NewLogger(cfg *config.LogCfg) *logx.Logger {
	path := filepath.Join(cfg.Location, cfg.Service)

	logger := logx.NewLogger().
		SetMaxLevel(cfg.Level).
		SetFormatter(logx.JSONFormatterFn).
		AddFileHandler(path, cfg.File).
		ConsoleOut(true)

	return logger
}
