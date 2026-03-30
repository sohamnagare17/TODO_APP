package routes

import (
	"net/http"
	"go-sqlite/handlers"
	"database/sql"
)

func SetupRoutes(db *sql.DB) {

	http.HandleFunc("/insert", handlers.InsertTask(db))
	// http.HandleFunc("/getall", handlers.GetAll(db))
	// http.HandleFunc("/get", handlers.GetTask(db))
	http.HandleFunc("/rename", handlers.RenameTask(db))
	http.HandleFunc("/changeStatus", handlers.ChangeStatus(db))
	http.HandleFunc("/delete", handlers.DeleteTask(db))
	// http.HandleFunc("/deleteCompletedTask", handlers.DeleteCompletedTask(db))
	http.HandleFunc("/insertmany", handlers.Insertmany(db))
	http.HandleFunc("/showtask", handlers.ShowTask(db))
	http.HandleFunc("/viewtask", handlers.ViewTask(db))


	http.HandleFunc("/task", handlers.InsertTask(db))
	http.HandleFunc("/user",handlers.InsertUser(db))
	http.HandleFunc("/users",handlers.GetAllUsers(db))
	http.HandleFunc("/user/id",handlers.GetUserById(db))
	http.HandleFunc("/task/userid",handlers.GetTaskByUserId(db))
	http.HandleFunc("/task/status/userid",handlers.GetTask(db))
	http.HandleFunc("/task/status",handlers.UpdateStatusOfTask(db))
	http.HandleFunc("/tasks",handlers.GetTasksBySorted(db))
	http.HandleFunc("/task/user",handlers.DeleteTask(db))



         //have to update each every route handler function
		// http.HandleFunc("POST /users", handlers.InsertUser(db))
		// http.HandleFunc("POST /users/{userId}/tasks", handlers.InsertTask(db))
		// http.HandleFunc("GET /users", handlers.GetAllUsers(db))
		// http.HandleFunc("GET /users/{id}", handlers.GetUserById(db))
		// http.HandleFunc("GET /users/{userId}/tasks", handlers.GetTaskByUserId(db))
		// http.HandleFunc("GET /tasks", handlers.GetTasksBySorted(db))
		// http.HandleFunc("PATCH /tasks/{taskId}/status", handlers.UpdateStatusOfTask(db))
		// http.HandleFunc("DELETE /users/{userId}/tasks/{taskId}", handlers.DeleteTask(db))
	

}