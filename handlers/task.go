package handlers

import(
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		"go-sqlite/models"
		"strconv"
)

func GetTaskByUserId(db *sql.DB) http.HandlerFunc{
	return func( writer http.ResponseWriter ,request *http.Request){

		var tasklist []models.Task;
		useridstr := request.URL.Query().Get("user_id")

		if useridstr ==""{
			log.Println("id required plz!")
			return
		}

		userid,err := strconv.Atoi(useridstr);
		if err!=nil{
			log.Println("id must be the number")
			return 
		}
		query := `SELECT id,name,status,userid FROM tasks1 WHERE userid=?`

		rows,err1:=db.Query(query,userid)
		for rows.Next(){
			var task models.Task

			err = rows.Scan(&task.ID,&task.NAME,&task.STATUS,&task.USERID)
			if err1!=nil{
				log.Println("error in scanning the data from the rows",err)
			}
			tasklist = append(tasklist, task)
		}
		if err!=nil{
			log.Println("error in fetching the data")
			return 
		}

		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":"the task of the user are",
			"tasklist":tasklist,
		})
	}
}

func InsertTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
        var newtask models.Task;

		err := json.NewDecoder(r.Body).Decode(&newtask)

		if err!=nil{
		log.Println("error in fetching the data")
		}

		query := `INSERT INTO tasks1 (name , status,userid) VALUES(?,?,?)`

		_, err = db.Exec(query,newtask.NAME,newtask.STATUS,newtask.USERID)

		if err != nil{
		log.Println("somthing went wrong to inserting the data ")
		return 
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
		"message":"the task inserted succesfully into database ",
		"taskname":newtask.NAME,
		"userid":newtask.USERID,
		})

	}
}

func GetTask(db *sql.DB) http.HandlerFunc{
	return func( writer http.ResponseWriter ,request *http.Request){

		var tasklist []models.Task;
		useridstr := request.URL.Query().Get("user_id")
		status := request.URL.Query().Get("status")

		if useridstr ==""{
			log.Println("id required plz!")
			return
		}

		userid,err := strconv.Atoi(useridstr);
		if err!=nil{
			log.Println("id must be the number")
			return 
		}
		query := `SELECT id,name,status,userid FROM tasks1 WHERE userid=? AND status=?`

		rows,err1:=db.Query(query,userid,status)
		for rows.Next(){
			var task models.Task

			err = rows.Scan(&task.ID,&task.NAME,&task.STATUS,&task.USERID)
			if err1!=nil{
				log.Println("error in scanning the data from the rows",err)
			}
			tasklist = append(tasklist, task)
		}
		if err!=nil{
			log.Println("error in fetching the data")
			return 
		}

		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":"the task of the user are",
			"tasklist":tasklist,
		})
	}
}

func UpdateStatusOfTask(db *sql.DB) http.HandlerFunc{
	return func ( writer http.ResponseWriter, request *http.Request){

		var req struct{
			TaskId int `json:"taskid"`
			UserId int `json:"user_id"`
			Status string `json:"status"`
		}
      err:= json.NewDecoder(request.Body).Decode(&req);
	  if err!=nil{
		log.Println("error in decoding the data")
		return 
	  }
       res, err := db.Exec(
			"UPDATE tasks1 SET status=? WHERE id=? AND userid=?",
			req.Status, req.TaskId, req.UserId,
		)

		if err!=nil{
			log.Println("error in the db execution",err)
			return
		}

		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":"status updated succesfully",
			"result":res,
		})

	
	}
}