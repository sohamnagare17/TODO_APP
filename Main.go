package main

import (
	"go-sqlite/database"
	"go-sqlite/handlers"
	"go-sqlite/repository"
	"go-sqlite/routes"
	"go-sqlite/services"
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	dbconn := db.Dbinit()

	defer dbconn.Close()

	repo := repository.NewTaskRepository(dbconn)
	service := services.NewTaskServices(repo)
	taskhandler := handlers.NewTaskHandler(service)

	repouser := repository.NewUserRepository(dbconn)
	userservices := services.NewUserServices(repouser)
	userhandler := handlers.NewUserHandler(userservices)

	routes.SetupRoutes(taskhandler, userhandler, dbconn)

	log.Println("server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
