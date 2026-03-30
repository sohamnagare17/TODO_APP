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
	http.HandleFunc("GET /user/{userid}",handlers.GetUserById(db))
	http.HandleFunc("GET /users/{userid}/tasks",handlers.GetTaskByUserId(db))
	http.HandleFunc("GET /users/{userid}/tasks/{status}",handlers.GetTaskByStatus(db))
	http.HandleFunc("PATCH /task/status",handlers.UpdateStatusOfTask(db))
	http.HandleFunc("GET /tasks/{userid}/{sort}",handlers.GetTasksBySorted(db))
	http.HandleFunc("DELETE /users/{userid}/tasks/{taskid}",handlers.DeleteTask(db))
	http.HandleFunc("PATCH /users/{userid}/tasks/{taskid}",handlers.RenameTask(db))

	

}