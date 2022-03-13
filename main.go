package main

import (
	"byd/api"
	"byd/messaging"
	"byd/repo"
	"byd/tool"
	"byd/workers"
	"log"
)

func main() {
	busCs, err := tool.GetBusConnectionString()
	if err != nil {
		log.Fatalln(err)
	}
	bus, err := messaging.NewRabbitMQBus(busCs, "download")
	if err != nil {
		log.Fatalln(err)
	}
	dbCs, err := tool.GetDbConnectionString()
	if err != nil {
		log.Fatalln(err)
	}
	rep, err := repo.NewPostgresRepo(dbCs)
	if err != nil {
		log.Fatalln(err)
	}
	wm := workers.NewWorkersManager(bus, rep)
	wm.SpawnWorker(workers.DOWNLOAD_WORKER)
	apiUrl, err := tool.GetApiUrl()
	if err != nil {
		log.Fatalln(err)
	}
	api := api.NewApiServer(apiUrl, bus, rep)
	api.StartApi()
}
