package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alleswebdev/marketplace-3d-factory/internal/api"
	"github.com/alleswebdev/marketplace-3d-factory/internal/config"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/queue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/savepoint"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/orders_updater"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/queuer"
)

func main() {
	cfg := config.GetAppConfig()

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		ServerHeader:  "go-app",
		AppName:       "Marketplace 3d factory",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1, http://localhost, http://127.0.0.1:4173, http://80.76.35.119",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Static("/", "./web/factory-front/dist")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wbClient := wb.NewClient(cfg.WbToken)

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	orderStore := order.New(dbpool)
	cardStore := card.New(dbpool)
	queueStore := queue.New(dbpool)
	savepointsStore := savepoint.New(dbpool)

	if err != nil {
		log.Fatal(err)
	}

	ordersUpdater := orders_updater.NewWorker(wbClient, orderStore)
	go ordersUpdater.Run(ctx)
	queueUpdater := queuer.NewWorker(dbpool, orderStore, savepointsStore, queueStore, cardStore)
	go queueUpdater.Run(ctx)

	appAPI := api.New(queueStore, cardStore, wbClient)
	app.Get("/api/list-queue", appAPI.ListQueue)
	app.Get("/api/update-cards", appAPI.UpdateCards) //todo выпилить
	app.Post("/api/set-complete", appAPI.SetComplete)
	app.Post("/api/set-printing", appAPI.SetPrinting)

	err = app.Listen(":" + strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}
