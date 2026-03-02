package main

import (
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/mroth/weightedrand/v3"
	"github.com/samber/lo"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	var span trace.Span
	_ = span

	chooser := lo.Must(weightedrand.NewChooser(
		weightedrand.NewChoice(slog.LevelInfo, 5),
		weightedrand.NewChoice(slog.LevelWarn, 3),
		weightedrand.NewChoice(slog.LevelError, 2),
	))

	sleepChooser := lo.Must(weightedrand.NewChooser(
		weightedrand.NewChoice(func() time.Duration { return randBetween(time.Millisecond, 100*time.Millisecond) }, 5),
		weightedrand.NewChoice(func() time.Duration { return randBetween(100*time.Millisecond, 150*time.Millisecond) }, 3),
		weightedrand.NewChoice(func() time.Duration { return randBetween(150*time.Millisecond, 300*time.Millisecond) }, 2),
	))
	_ = sleepChooser

	httpMethodChooser := lo.Must(weightedrand.NewChooser(
		weightedrand.NewChoice(semconv.HTTPRequestMethodGet, 50),
		weightedrand.NewChoice(semconv.HTTPRequestMethodPost, 40),
		weightedrand.NewChoice(semconv.HTTPRequestMethodDelete, 10),
	))

	userCount := 0

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				level := chooser.Pick()
				httpMethod := httpMethodChooser.Pick()

				logger.Info("",
					slog.String("level2", level.String()),
					slog.String("type", "http_request"),
					slog.String(string(httpMethod.Key), httpMethod.Value.AsString()),
					slog.Int("user_count", userCount),
				)
				time.Sleep(sleepChooser.Pick()())
				userCount++
			}
		}()
	}

	http.ListenAndServe(":8080", nil)
}

func randBetween(min, max time.Duration) time.Duration {
	return time.Duration(rand.Intn(int(max-min+1)) + int(min))
}
