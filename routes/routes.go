package routes

import (
	"database/sql"
	"go-sqlite/handlers"
	"net/http"
)

func SetupRoutes(taskhandler *handlers.TaskHandler, userhandler *handlers.UserHandler, db *sql.DB) {

    //done
	http.HandleFunc("POST /users/{userid}/tasks", taskhandler.InsertTask)
    
	//done
	http.HandleFunc("POST /user", userhandler.InsertUser)

	//done
	http.HandleFunc("GET /users", userhandler.GetAllUsers)
    
	//done
	http.HandleFunc("GET /users/{userid}", userhandler.GetUserById)

	//done
	http.HandleFunc("GET /users/{userid}/tasks", taskhandler.GetTaskByUserId)

	
	http.HandleFunc("PATCH /users/{userid}/tasks/{taskid}", taskhandler.UpdateTask)

    //done
	http.HandleFunc("DELETE /users/{userid}/tasks/{taskid}", taskhandler.DeleteTask)
}
