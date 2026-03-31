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

		// userIDStr := r.PathValue("userId")

		// userID, err := strconv.Atoi(userIDStr)
		// if err != nil {
		//     http.Error(w, "invalid userId", http.StatusBadRequest)
		//     return
		// }

		useridstr := request.PathValue("userid")

		if useridstr == "" {
			log.Println("id required plz!")
			return
		}

		userid, err := strconv.Atoi(useridstr)
		if err != nil {
			log.Println("id must be the number")
			return
		}
		query := `SELECT id,name,status,userid FROM tasks1 WHERE userid=?`

		rows, err1 := db.Query(query, userid)
		for rows.Next() {
			var task models.Task

			err = rows.Scan(&task.ID, &task.NAME, &task.STATUS, &task.USERID)
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
		//time.Now()
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

func GetTaskByStatus(db *sql.DB) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var tasklist []models.Task
		useridstr := request.PathValue("userid")
		status := request.PathValue("status")

		if useridstr == "" {
			log.Println("id required plz!")
			return
		}

		userid, err := strconv.Atoi(useridstr)
		if err != nil {
			log.Println("id must be the number")
			return
		}
		query := `SELECT id,name,status,userid FROM tasks1 WHERE userid=? AND status=?`

		rows, err1 := db.Query(query, userid, status)
		for rows.Next() {
			var task models.Task

			err = rows.Scan(&task.ID, &task.NAME, &task.STATUS, &task.USERID)
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

func GetTasksBySorted(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		useridstr := r.PathValue("userid")
		if useridstr == "" {
			log.Println("enter a valid user id")
		}
		userid, err := strconv.Atoi(useridstr)
		if err != nil {
			log.Println("Error")
		}

		sort := r.PathValue("sort")

		query := `SELECT * FROM tasks1
			WHERE userId = ?`

		switch sort {
		case "name":
			query += "ORDER BY name ASC"
		case "createdAt":
			query += "ORDER BY createdAt DESC"
		case "updatedAt":
			query += "ORDER BY updatedAt DESC"
		default:
			query += "ORDER BY id DESC"

		}
		log.Println(query)
		rows, err := db.Query(query, userid)
		if err != nil {
			log.Fatal("Internal server Error")
		}
		var task models.Task
		tasks := []models.Task{}

		for rows.Next() {
			rows.Scan(&task.ID, &task.NAME, &task.STATUS, &task.USERID, &task.CreatedAt, &task.UpdatedAt)
			tasks = append(tasks, task)
		}
		rows.Close()
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tasks": tasks,
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


func UpdateTask(db *sql.DB) http.HandlerFunc{
	return  func(writer http.ResponseWriter, request *http.Request) {

		userid:=request.PathValue("userid")
		taskid:=request.PathValue("taskid")

		if userid==""||taskid==""{
			http.Error(writer,"Enter a valid task and user id",400)
			log.Println("Id is required")
			return
		}
		uid,err:=strconv.Atoi(userid)
		if err!=nil{
			http.Error(writer,"Invalid User id",400)
			log.Println("Enter a valid user id")
			return
		}
		tid,err:=strconv.Atoi(taskid)
		if err!=nil{
			http.Error(writer,"Enter a valid task id",400)
			log.Println("Enter a valid task id")
			return
		}
		
		var reqbody struct{
			Name string `json:"name"`
			Status string `json:"status"`
		}
		err=json.NewDecoder(request.Body).Decode(&reqbody)
		if err!=nil{
			http.Error(writer,"Invalid body",400)
			log.Println("Error in the request Body")
			return
		}

		var query string
		var res sql.Result
		switch{
		case reqbody.Name!=""&&reqbody.Status!="":
			query=`UPDATE tasks1 
				SET name=? ,status=?,updatedAt=CURRENT_TIMESTAMP
				WHERE id=? AND userid=?`
			res,err=db.Exec(query,reqbody.Name,reqbody.Status,tid,uid)
		
		case reqbody.Name!="":
			query=`UPDATE tasks1 
				SET name=?,updatedAt=CURRENT_TIMESTAMP
				WHERE id=? AND userid=?`
			res,err=db.Exec(query,reqbody.Name,tid,uid)
		
		case reqbody.Status!="":
			query=`UPDATE tasks1 
				SET status=?,updatedAt=CURRENT_TIMESTAMP
				WHERE id=? AND userid=?`
			res,err=db.Exec(query,reqbody.Status,tid,uid)

		default:
			http.Error(writer,"Nothing to update",400)
			return
		}

		if err!=nil{
			http.Error(writer,"Internal Server Error",500)
			log.Println("Internal server error")
			return 
		}

		rows,_:=res.RowsAffected()
		if rows==0{
			http.Error(writer,"Task not found",400)
			log.Println("task not found")
			return 
		}
		writer.Header().Set("Content-type","application/json")
		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message": "task updated successfully",
			"rows":    rows,
		})
	}
}