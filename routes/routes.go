package routes

import (
	"database/sql"
	"go-sqlite/handlers"
	"net/http"
)

func SetupRoutes(taskhandler *handlers.TaskHandler, userhandler *handlers.UserHandler, db *sql.DB) {

	http.HandleFunc("POST /users/{userid}/tasks", taskhandler.InsertTask)

	http.HandleFunc("POST /user", userhandler.InsertUser)

	http.HandleFunc("GET /users", userhandler.GetAllUsers)

	http.HandleFunc("GET /users/{userid}", userhandler.GetUserById)

	http.HandleFunc("GET /users/{userid}/tasks", taskhandler.GetTaskByUserId)

	http.HandleFunc("PATCH /users/{userid}/tasks/{taskid}", taskhandler.UpdateTask)

	http.HandleFunc("DELETE /users/{userid}/tasks/{taskid}", taskhandler.DeleteTask)
}
