package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"danilamukhin/serv_go/internal/api"
	"danilamukhin/serv_go/internal/model"
	"danilamukhin/serv_go/internal/pgx"
	"danilamukhin/serv_go/internal/service"
	"danilamukhin/serv_go/pkg/server"

	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ConfigURL = "/Users/danilamukhin/Desktop/Работа/serv_go/cmd/config.toml"

func main() {
	ctx := context.Background()

	//config toml
	conf, err := pgx.ConfigData(ConfigURL)
	if err != nil {
		fmt.Println(err)
	}

	// connecting to postgres
	pgxConnURL := "postgresql://" + conf.Login + ":" + conf.Password + "@" + conf.IP + "/" + conf.DBName + "?&sslmode=" + conf.Sslmode
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

	// Archive
	// Инициализация планировщика
	s := gocron.NewScheduler(time.UTC)

	// Первый запуск при запуске программы
	Archive(repo, ctx, conf)

	// Крон-выражение для запуска функции Archive каждые 3 месяца
	s.Every(3).Months().Do(func() {
		Archive(repo, ctx, conf)
	})

	// s.Every(5).Second().Do(func() {
	// 	fmt.Println("5 Seconds")
	// })

	// Запуск планировщика
	s.StartAsync()

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

func Archive(repo pgx.Repo, ctx context.Context, conf model.TableConfig) {
	var year, quarter []string
	var nowQuarter, nowYear int
	var err error
	var tableName string

	// Получаем значения минимального квартала и года из таблицы
	year, quarter, err = pgx.Repo.MinTimestamp(repo, ctx, conf)
	if err != nil {
		fmt.Println(err)
	}

	// Переменные с годом и кварталом в формате int
	yearInt := make([]int, len(year))
	quarterInt := make([]int, len(quarter))

	// Записываем в эти переменные значения из формата str в int
	for i := range year {
		yr, err := strconv.Atoi(year[i])
		if err != nil {
			fmt.Println(err)
		}
		yearInt[i] = yr
		qr, err := strconv.Atoi(quarter[i])
		if err != nil {
			fmt.Println(err)
		}
		quarterInt[i] = qr
	}

	// Получаем текущий год и квартал
	curentTime := time.Now()
	nowYear = curentTime.Year()
	nowQuarter = (int(curentTime.Month())-1)/3 + 1

	// Разница между текущим и минимальным кварталом
	diffQuarter := make([]int, len(conf.Tables))

	// Проходимся по всем таблицам, которые есть в conf.Tables
	for i := range conf.Tables {
		Table := conf.Tables[i]
		tableName = Table.TableName

		// Разница между текущим кварталом и минимальным
		diffQuarter[i] = (nowYear-yearInt[i])*4 + nowQuarter - quarterInt[i]

		// Пока разница >= 2 делаем
		for diffQuarter[i] >= 2 {

			// Создаём Архив и таблицу в архиве
			err := pgx.Repo.CreateSchemaAndTable(repo, ctx, conf, year[i], quarter[i], tableName, i)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			// Переносим данные за квартал в архивную таблицу
			err = pgx.Repo.MoveQuarter(repo, ctx, conf, year[i], quarter[i], tableName)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			// Удаляем данные за квартал из исходной таблицы
			err = pgx.Repo.DeleteQuarter(repo, ctx, conf, year[i], quarter[i], tableName)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			// Находим новый минимальный год и квартал, после чего находим разницу между новым минимальным и текущим
			year, quarter, err = pgx.Repo.MinTimestamp(repo, ctx, conf)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			yearInt[i], err = strconv.Atoi(year[i])
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			quarterInt[i], err = strconv.Atoi(quarter[i])
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			diffQuarter[i] = (nowYear-yearInt[i])*4 + nowQuarter - quarterInt[i]
		}
	}

}
