package main

import (
	"context"
	"fmt"
	"log"
	"note_api/db/connection"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
	"note_api/http"
	"note_api/note"
)

// todo: сделать возможным удаление версий заметки
// решить проблему: 
/*
	1. когда меняю и название и содержимое заметки, то в таблицу note_versions летит начальная версия заметки И версия заметки с измененным только названием, но прежним содержимым
	2. сделать индитификатором не название заметки, а сделать прям айди, чтобы было удобней получать доступ ко всему и различать, так как в note_versions этого очень не хватает после добавления бд
*/

func main()  {
	ctx := context.Background()
	pool, err := connection.CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

		driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	

	repo := note.NewRepository(pool)
	httpHandler := http.NewHTTPHandlers(repo)
	httpServer := http.NewHTTPServer(httpHandler)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start http server:", err)
	}
}