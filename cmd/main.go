package main

import (
	"crypto/rand"
	"denet/db"
	"denet/server"
	"log"
	"log/slog"
)

func generateSecretKey() []byte {
	key := make([]byte, 32) // 256 бит
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}
func main() {
	// migrations - or use migrations via docker and binary file
	// this way is looks simpler
	db.InsertMigrations()
	slog.Info("Migrations inserted")

	// Приватная обертка над драйвером
	db := db.New("postgres://denet:password@db:5432/mydb?sslmode=disable")
	defer db.Close()
	slog.Info("Connected to database")

	// Приватная обертка над сервером + JWT
	s := server.New(":8080", db, generateSecretKey())
	slog.Info("Server started on :8080")
	log.Fatal(s.ListenAndServe())
}
