package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/pangolin-do-golang/thumb-processor-worker/docs"
	dbAdapter "github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/db"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/rest/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Thumb processor worker
// @version 0.1.0
// @description Hackathon

// @host localhost:8080
// @BasePath /
func main() {
	_, err := initDb()
	if err != nil {
		panic(err)
	}

	/**

	_ := dbAdapter.NewPostgresThumbRepository(db)

	_ := thumb.NewThumbService()

	_ := dbAdapter.NewPostgresThumbRepository(db)

	**/

	restServer := server.NewRestServer(&server.RestServerOptions{})

	restServer.Serve()
}

func initDb() (*gorm.DB, error) {
	_ = godotenv.Load()
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	err = db.AutoMigrate(
		&dbAdapter.ThumbPostgres{},
	)
	if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}
