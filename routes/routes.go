package routes

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-sqlite/handlers"
	"go-sqlite/middelware"
	"net/http"
)

func SetupRoutes(taskhandler *handlers.TaskHandler, userhandler *handlers.UserHandler, db *sql.DB) {

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("POST /users/{userid}/tasks", middleware.MetricsMiddleware(http.HandlerFunc(taskhandler.InsertTask)))
	//http.HandleFunc("POST /users/{userid}/tasks", taskhandler.InsertTask)

	http.Handle("POST /user", middleware.MetricsMiddleware(http.HandlerFunc(userhandler.InsertUser)))
	//http.HandleFunc("POST /user", userhandler.InsertUser)

	http.Handle("GET /users", middleware.MetricsMiddleware(http.HandlerFunc(userhandler.GetAllUsers)))
	//http.HandleFunc("GET /users", userhandler.GetAllUsers)

	http.Handle("GET /users/{userid}", middleware.MetricsMiddleware(http.HandlerFunc(userhandler.GetUserById)))
	//http.HandleFunc("GET /users/{userid}", userhandler.GetUserById)

	http.Handle("GET /users/{userid}/tasks", middleware.MetricsMiddleware(http.HandlerFunc(taskhandler.GetTaskByUserId)))
	//http.HandleFunc("GET /users/{userid}/tasks", taskhandler.GetTaskByUserId)

	http.Handle("PATCH /users/{userid}/tasks/{taskid}", middleware.MetricsMiddleware(http.HandlerFunc(taskhandler.UpdateTask)))
	//http.HandleFunc("PATCH /users/{userid}/tasks/{taskid}", taskhandler.UpdateTask)

	http.Handle("DELETE /users/{userid}/tasks/{taskid}", middleware.MetricsMiddleware(http.HandlerFunc(taskhandler.DeleteTask)))
	//http.HandleFunc("DELETE /users/{userid}/tasks/{taskid}", taskhandler.DeleteTask)
}
