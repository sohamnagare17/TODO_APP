package handlers

import (
	"database/sql"
	"encoding/json"
	"go-sqlite/models"
	"log"
	"net/http"
	"strconv"
)

func GetTaskByCreatedAt(db *sql.DB) http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("userid")
	if id == "" {
		log.Fatal("Id is missing")
	}
	iid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Enter a valid id")
	}
	query := `SELECT * FROM tasks 
			WHERE userId = ?
			ORDER BY createdAt DESC`

	rows, err := db.Query(query, iid)
	if err != nil {
		log.Fatal("Internal server Error")
	}
	var task models.Task
	tasks := []models.Task{}

	for rows.Next() {
		rows.Scan(&task.ID, &task.NAME, &task.STATUS, &task.CreatedAt, &task.UpdatedAt, &task.USERID)
		tasks = append(tasks, task)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

}
