package routes

import (
	"net/http"
	"go-sqlite/handlers"
	"database/sql"
)

func SetupRoutes(db *sql.DB) {


	http.HandleFunc("POST /users/{userId}/tasks", handlers.InsertTask(db))
	http.HandleFunc("POST /user",handlers.InsertUser(db))
	http.HandleFunc("GET /users",handlers.GetAllUsers(db))
	http.HandleFunc("GET /users/{userid}",handlers.GetUserById(db))

	http.HandleFunc("GET /users/{userid}/tasks",handlers.GetTaskByUserId(db))


	http.HandleFunc("PATCH /users/{userid}/tasks/{taskid}",handlers.UpdateTask(db))

	http.HandleFunc("DELETE /users/{userid}/tasks/{taskid}",handlers.DeleteTask(db))

	

}