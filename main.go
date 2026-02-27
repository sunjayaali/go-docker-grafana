package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mroth/weightedrand/v3"
	"github.com/samber/lo"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	chooser := lo.Must(weightedrand.NewChooser(
		weightedrand.NewChoice(slog.LevelInfo, 5),
		weightedrand.NewChoice(slog.LevelWarn, 3),
		weightedrand.NewChoice(slog.LevelError, 2),
	))

	for range time.NewTicker(500 * time.Millisecond).C {
		level := chooser.Pick()
		logger.Info("",
			slog.String("level", level.String()),
			slog.String("level2", level.String()),
		)
	}
	http.ListenAndServe(":8080", nil)
}
