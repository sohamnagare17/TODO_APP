package handlers

import (
	"database/sql"
	"encoding/json"
	// "go-sqlite/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != http.MethodPatch {
			http.Error(writer, "Invalid method", 405)
			log.Println("Invalid method type")
			return
		}
		userid := request.PathValue("userid")
		taskid := request.PathValue("taskid")

		if userid == "" || taskid == "" {
			http.Error(writer, "Enter a valid task and user id", 400)
			log.Println("Id is required")
			return
		}
		uid, err := strconv.Atoi(userid)
		if uid < 0 {
			http.Error(writer, "userid shoud be positive", 400)
			log.Println("Enter a valid user id")
			return
		}
		if err != nil {
			http.Error(writer, "Invalid User id", 400)
			log.Println("Enter a valid user id")
			return
		}
		tid, err := strconv.Atoi(taskid)
		if tid < 0 {
			http.Error(writer, "Task Id Should be Positive", 400)
			log.Println("Enter a valid Task id")
			return
		}
		if err != nil {
			http.Error(writer, "Enter a valid task id", 400)
			log.Println("Enter a valid task id")
			return
		}

		var reqbody struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		}
		err = json.NewDecoder(request.Body).Decode(&reqbody)
		if err != nil {
			http.Error(writer, "Invalid body", 400)
			log.Println("Error in the request Body")
			return
		}

		name := strings.TrimSpace(reqbody.Name)
		var query string
		if reqbody.Name != "" && name == "" {
			http.Error(writer, "Name should not be empty", 400)
			return
		}

		var res sql.Result
		switch {
		case reqbody.Name != "" && reqbody.Status != "":
			query = `UPDATE tasks1 
				SET name=? ,status=?,updatedAt=CURRENT_TIMESTAMP
				WHERE id=? AND userid=?`
			res, err = db.Exec(query, reqbody.Name, reqbody.Status, tid, uid)

		case reqbody.Name != "":
			query = `UPDATE tasks1 
				SET name=?,updatedAt=CURRENT_TIMESTAMP
				WHERE id=? AND userid=?`
			res, err = db.Exec(query, reqbody.Name, tid, uid)

		case reqbody.Status != "":
			query = `UPDATE tasks1 
				SET status=?,updatedAt=CURRENT_TIMESTAMP
				WHERE id=? AND userid=?`
			res, err = db.Exec(query, reqbody.Status, tid, uid)

		default:
			http.Error(writer, "Nothing to update", 400)
			return
		}

		if err != nil {
			http.Error(writer, "Internal Server Error", 500)
			log.Println("Internal server error")
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			http.Error(writer, "Task not found", 400)
			log.Println("task not found")
			return
		}
		writer.Header().Set("Content-type", "application/json")
		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message": "task updated successfully",
		})
	}
}
