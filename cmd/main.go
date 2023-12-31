package main

import (
	"context"
	"danilamukhin/serv_go/internal/api"
	"danilamukhin/serv_go/internal/pgx"
	"danilamukhin/serv_go/internal/service"
	"danilamukhin/serv_go/pkg/server"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	// connecting to postgres
	pgxConnURL := "postgresql://test:test@localhost:5432/postgres?&sslmode=disable"
	//pgxConnURL := "postgresql://index:Yfhenj@localhost:5432/index?&sslmode=disable"
	pool, err := pgxpool.New(ctx, pgxConnURL)
	if err != nil {
		log.Fatalln(err)
	}

	repo := pgx.New(pool)

	srv := service.NewService(repo)

	a := api.InitApi(srv)
	httpServer := server.NewServer(a.InitRouter())
	go func() {
		if err := httpServer.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := httpServer.Stop(ctx); err != nil {
		log.Fatalln(err)
	}
}
