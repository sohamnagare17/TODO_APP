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
	http.HandleFunc("/task/name",handlers.GetTasksByName(db))
	http.HandleFunc("/task/createdAt",handlers.GetTaskByCreatedAt(db))
	http.HandleFunc("/task/updatedAt",handlers.GetTaskUpdatedAt(db))
	
}