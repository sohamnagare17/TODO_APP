package handlers

import (
	"database/sql"
	"encoding/json"
	"go-sqlite/models"
	"log"
	"net/http"
	"strconv"
)

func CompletedTasks(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.URL.Query().Get("userId")
		uid, err := strconv.Atoi(userId)
		if err != nil {
			log.Println("Enter valid id")
		}

		query := `SELECT * FROM tasks
				WHERE userId=?
				AND status='DONE'`

		tasks := []models.Task{}
		rows, err := db.Query(query, uid)
		for rows.Next() {
		    var task models.Task
			rows.Scan(&task.ID,&task.NAME,&task.STATUS,&task.CreatedAt, &task.UpdatedAt,&task.Userid)
			tasks = append(tasks, task)
		}
		rows.Close()

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}
