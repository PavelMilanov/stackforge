package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PavelMilanov/stackforge/config"
	"github.com/PavelMilanov/stackforge/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:00",
	})
	env, err := config.NewEnv()
	if err != nil {
		logrus.Fatal(err)
	}
	handler := handlers.NewHandler(env)
	router := handler.InitRouters()
	s := http.Server{
		Addr:              ":1323",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		logrus.Info("Сервер запущен")
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}
}
