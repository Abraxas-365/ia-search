package main

import (
	"context"
	"log"
	"os"

	"github.com/Abraxas-365/ia-search/internal/application"
	"github.com/Abraxas-365/ia-search/internal/infrastructure/database"
	"github.com/Abraxas-365/ia-search/internal/infrastructure/repository"
	"github.com/Abraxas-365/ia-search/internal/infrastructure/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	ctx := context.Background()

	postgresUri := os.Getenv("POSTGRES_URI")
	openIaKEy := os.Getenv("OPENIA_KEY")
	conn, err := database.NewConnection(postgresUri, 5432, "myuser", "mypassword", "mydb")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	defer conn.Close()

	err = conn.RunMigrations(context.Background())
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewParagraphRepository(conn.Pool)

	app := application.NewApplication(repo, openIaKEy)

	if err := app.ParseFile(ctx, "example.txt", 800, 10); err != nil {
		log.Fatal(err)
	}

	appF := fiber.New()
	appF.Use(cors.New())
	appF.Use(logger.New())
	rest.ControllerFactory(appF, app)

	appF.Listen(":3000")

}
