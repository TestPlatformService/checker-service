package main

import (
	"checker/api"
	"checker/config"
	"checker/logs"
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

	router := api.Router(logger)
	log.Println("Checker Service is running on port 50054")
	err = router.Run(cfg.CHECKER_SERVICE)
	if err != nil {
		logger.Error(err.Error())
	}
}
