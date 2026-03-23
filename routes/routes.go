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
}