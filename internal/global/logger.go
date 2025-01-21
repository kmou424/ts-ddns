package global

import (
	"github.com/gookit/slog"
)

func initLogger() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = false
		f.TimeFormat = "2006-01-02 15:04:05.000"
		if !Debug.Load() {
			f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}\n")
		} else {
			f.SetTemplate("[{{datetime}}] [{{level}}] [{{caller}}] {{message}} {{data}} {{extra}}\n")
		}
	})
}
