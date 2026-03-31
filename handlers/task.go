package handlers

import (
	"database/sql"
	"encoding/json"
	"go-sqlite/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetTaskByUserId(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var tasklist []models.Task
		useridstr := request.PathValue("userid")
		status := request.URL.Query().Get("status")
		sortby := request.URL.Query().Get("sortby")
		order := request.URL.Query().Get("order")

		if useridstr == "" {
			log.Println("id required plz!")
			return
		}

		userid, err := strconv.Atoi(useridstr)
		if err != nil {
			log.Println("id must be the number")
			return
		}

		validfields := map[string]bool{
			"name":true,
			"createdAt":true,
			"updatedAt":true,
		}
		query := `SELECT * FROM tasks1 WHERE userid=?`
		parameters := []interface{}{userid}

		if status!=""{
			query=query+" AND status=? "
			parameters = append(parameters,status)
		}
		

		if validfields[sortby]{
			query = query + " ORDER BY " + sortby
			
			if order=="ASC" || order=="asc"{
				query+=" ASC "
			}else{
				query+="DESC"
			}
		}else{
			query+= " ORDER BY createdAt DESC"
		}
		log.Println("Query:", query)
        log.Println("Values:", parameters)

		rows, err1 := db.Query(query, parameters...)
		if err1!=nil{
			log.Println("something went wrong in the execution of the database query")
		}
		for rows.Next() {
			var task models.Task

			err = rows.Scan(&task.ID, &task.NAME, &task.STATUS, &task.USERID,&task.CreatedAt,&task.UpdatedAt)
			if err1 != nil {
				log.Println("error in scanning the data from the rows", err)
			}
			tasklist = append(tasklist, task)
		}
		if err != nil {
			log.Println("error in fetching the data")
			return
		}

		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":  "the task of the user are",
			"tasklist": tasklist,
		})
	}
}

func InsertTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newtask models.Task

		userIDStr := r.PathValue("userId")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "invalid userId", http.StatusBadRequest)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&newtask)

		if err != nil {
			log.Println("error in fetching the data")
		}

		query := `INSERT INTO tasks1 (name ,status,userid,createdAt,updatedAt) VALUES(?,?,?,?,?)`
		now := time.Now().UTC().Format(time.RFC3339)

		_, err = db.Exec(query, newtask.NAME, newtask.STATUS, userID, now, now)

		if err != nil {
			log.Println("somthing went wrong to inserting the data ", err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "the task inserted succesfully into database ",
			"taskname": newtask.NAME,
			"userid":   userID,
		})

	}
}



func UpdateStatusOfTask(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var req struct {
			TaskId int    `json:"taskid"`
			UserId int    `json:"user_id"`
			Status string `json:"status"`
		}
		err := json.NewDecoder(request.Body).Decode(&req)

		if err != nil {
			log.Println("error in decoding the data")
			return
		}

		res, err := db.Exec(
			"UPDATE tasks1 SET status=?, updatedAt=CURRENT_TIMESTAMP WHERE id=? AND userid=? ",
			req.Status, req.TaskId, req.UserId,
		)

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Println("RowsAffected error:", err)
			http.Error(writer, "Error checking result", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(writer, "Task not found", http.StatusNotFound)
			return
		}

		writer.Header().Set("Content-type", "application/json")
		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message": "status updated succesfully",
			"result":  res,
		})

	}
}


func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idstr := request.PathValue("taskid")
		useridstr := request.PathValue("userid")

		if idstr == "" || useridstr == "" {
			log.Println("userid and taskid  required plz provide ids")
			return
		}

		id, err := strconv.Atoi(idstr)
		if err != nil {
			log.Println("id must be integer", err)
			return
		}

		userid, err1 := strconv.Atoi(useridstr)
		if err1 != nil {
			log.Println("userid must be integer", err1)
			return
		}

		query := `DELETE FROM tasks1 WHERE userid=? AND id=?`

		result, err := db.Exec(query, userid, id)

		if err != nil {
			log.Println("error while executing the database query", err)
			return
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			log.Println("error in checking rows affected", err)
			return
		}

		if rowsAffected == 0 {
			json.NewEncoder(writer).Encode(map[string]interface{}{
				"error": "task not found",
			})
			return
		}

		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":        "task deleted succesfully",
			"deleted userid": userid,
			"deleted task":   id,
		})
	}
}

func RenameTask(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		userid := request.PathValue("userid")
		taskid := request.PathValue("taskid")

		if userid == "" || taskid == "" {
			log.Println("userid and taskid  required plz provide ids")
			return
		}

		tid, err := strconv.Atoi(taskid)
		if err != nil {
			log.Println("Taskid must be integer", err)
			return
		}

		uid, err1 := strconv.Atoi(userid)
		if err1 != nil {
			log.Println("userid must be integer", err1)
			return
		}

		var task models.Task
		json.NewDecoder(request.Body).Decode(&task)

		query := `UPDATE tasks1 SET NAME=?, updatedAt=CURRENT_TIMESTAMP 
				WHERE id=? AND userid=?`

		res, err := db.Exec(query, task.NAME, tid, uid)

		if err != nil {
			log.Println("error while executing the database query", err)
			return
		}

		rowsAffected, err := res.RowsAffected()

		if err != nil {
			log.Println("error in checking rows affected", err)
			return
		}

		if rowsAffected == 0 {
			json.NewEncoder(writer).Encode(map[string]interface{}{
				"error": "task not found",
			})
			return
		}

		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":        "task rename succesfully",
		})

	}
}
