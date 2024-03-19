package main

import (
	"context"
	"fmt"
	"os"
	"task/app/database"
	"task/app/flood-control"

	"github.com/joho/godotenv"
)

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

var _ FloodControl = &floodcontroller.FloodController{}

func main() {
	N := 1000
	K := 2
	env, _ := godotenv.Read(".env")
	user := env["POSTGRES_USER"]
	password := env["POSTGRES_PASSWORD"]
	host := env["POSTGRES_HOSTNAME"]
	port := env["POSTGRES_PORT"]
	name := env["POSTGRES_DB"]
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, name)
	psgl, err := database.NewPostgresConnection(dbURL, 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db := database.NewDatabase(psgl)

	floodController := floodcontroller.NewFloodController(N, K, db)

	allowed, err := floodController.Check(context.TODO(), 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(allowed)
}
