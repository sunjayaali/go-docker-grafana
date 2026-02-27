package main

import (
	"log/slog"
	"math/rand"
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

	sleepChooser := lo.Must(weightedrand.NewChooser(
		weightedrand.NewChoice(func() time.Duration { return randBetween(time.Millisecond, 100*time.Millisecond) }, 5),
		weightedrand.NewChoice(func() time.Duration { return randBetween(100*time.Millisecond, 300*time.Millisecond) }, 3),
		weightedrand.NewChoice(func() time.Duration { return randBetween(300*time.Millisecond, 1000*time.Millisecond) }, 2),
	))

	go func() {
		for {
			level := chooser.Pick()
			logger.Info("",
				slog.String("level", level.String()),
				slog.String("level2", level.String()),
				slog.String("type", "http_request"),
			)
			time.Sleep(sleepChooser.Pick()())
		}
	}()
	http.ListenAndServe(":8080", nil)
}

func randBetween(min, max time.Duration) time.Duration {
	return time.Duration(rand.Intn(int(max-min+1)) + int(min))
}
