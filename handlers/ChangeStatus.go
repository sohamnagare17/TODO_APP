package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func ChangeStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.URL.Query().Get("userId")
		if userId == "" {
			log.Println("Enter valid userId")
		}
		uid, err := strconv.Atoi(userId)
		if err != nil {
			log.Println("user must be number ")
			return
		}

		idstr := r.URL.Query().Get("id")
		if idstr == "" {
			log.Println("Id is required plz provide id")
			return
		}
		id, err := strconv.Atoi(idstr)
		if err != nil {
			log.Println("id must be number ")
			return
		}

		query := `UPDATE tasks
		        SET status=CASE
				when status="pending" THEN "DONE"
				ELSE "pending"
				END,
				updatedAt=CURRENT_TIMESTAMP
				WHERE id=? AND userId=? `

		_, err1 := db.Exec(query, id, uid)
		if err1 != nil {
			log.Println("error in update  the status")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "update the status succesfully ",
		})
	}
}
