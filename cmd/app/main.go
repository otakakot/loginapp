package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/otakakot/loginapp/internal/base/firebase"
	"github.com/otakakot/loginapp/internal/base/pocketbase"
	"github.com/otakakot/loginapp/internal/base/supabase"
	"github.com/otakakot/loginapp/internal/handler"
)

func main() {
	port := cmp.Or(os.Getenv("PORT"), "8080")

	pockatbaseURL := cmp.Or(os.Getenv("POCKETBASE_URL"), "http://localhost:7070")

	pb := pocketbase.New(pockatbaseURL)

	supabaseProjectReference := cmp.Or(os.Getenv("SUPABASE_PROJECT_REFERENCE"), "http://localhost:7070")

	supbaseAPIKey := cmp.Or(os.Getenv("SUPABASE_API_KEY"), "")

	sb := supabase.New(supabaseProjectReference, supbaseAPIKey)

	firebaseAPIKey := cmp.Or(os.Getenv("FIREBASE_API_KEY"), "")

	fb := firebase.New(firebaseAPIKey)

	auhtn := handler.New(
		fb,
		sb,
		pb,
	)

	health := &handler.Health{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", auhtn.Handle)

	mux.HandleFunc("GET /health", health.Handle)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		slog.Info("start server listen")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	slog.Info("start server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	slog.Info("done server shutdown")
}
