package main

import (
	"checker/api"
	"checker/config"
	"checker/logs"
	"checker/pkg"
	"checker/service"
	"checker/storage"
	"checker/storage/postgres"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	logger := logs.NewLogger()

	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := storage.NewStoragePro(db, logger)
	client := pkg.NewClients(cfg)
	service := service.NewService(logger, storage, client)

	router := api.Router(logger, service)
	log.Println("Checker Service is running on port 50054")
	err = router.Run(cfg.CHECKER_SERVICE)
	if err != nil {
		logger.Error(err.Error())
	}
}
