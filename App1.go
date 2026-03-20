package main

import (
	
	"net/http"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"go-sqlite/db"
	"strconv"
	
)

type Task struct{
	ID int  `json:"id"`
	NAME string `json:"name"`
	STATUS  string `json:"status"`
}

func main(){
	
  dbconn := db.Dbinit()

  defer dbconn.Close()


	file , err := os.OpenFile("app.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)

	if err!=nil{
		log.Fatal(err)
	}
     defer file.Close()

    log.SetOutput(file)

	// to insert a task into database
	http.HandleFunc("/insert",InsertTask(dbconn))

	// to get the all tasks from database
	http.HandleFunc("/getall",GetAll(dbconn))

	//to get one task by id 
	http.HandleFunc("/get",gettask(dbconn))

	// to rename the task by id 
	http.HandleFunc("/rename",RenameTask(dbconn))

	//to change the status of the task
	http.HandleFunc("/ChangeStatus",ChangeStatus(dbconn))

	//to delete the task from the table 
	http.HandleFunc("/delete",DeleteTask(dbconn))

	//delete the task which are done 
	http.HandleFunc("/DeleteCompletedTask",DeleteCompletedTask(dbconn))

	//inserting the multiple task into the datbase
	http.HandleFunc("/insertmany",Insertmany(dbconn))

	//show the task by the page (pagination concept)
	http.HandleFunc("/showtask",ShowTask(dbconn))

	//pagination by cursor 
	http.HandleFunc("/viewtask",ViewTask(dbconn))

	http.ListenAndServe(":8080", nil)
}


func ViewTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request){

		  limitstr := r.URL.Query().Get("limit")
		  afterindestr := r.URL.Query().Get("afterind")

		  limit,_:= strconv.Atoi(limitstr)
		  afterind,_ := strconv.Atoi(afterindestr)

		  rows, err := db.Query("SELECT * FROM tasks WHERE id > ? LIMIT  ?" ,afterind,limit)
		  if err!=nil{
			log.Println("error in fetching the parameters")
		  }
 defer rows.Close()

		  var list []Task

		  for rows.Next(){
               var task Task
			   err := rows.Scan(&task.ID,&task.NAME,&task.STATUS)

			   if err!=nil{
				log.Println("wrong in the scanning the data from the rows ")
			   }

			   list = append(list,task)
		  }

      json.NewEncoder(w).Encode(list)

	}


}

func ShowTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

	   limitstr:= r.URL.Query().Get("limit")
	   pagenostr:= r.URL.Query().Get("pageno")

	   limit,_:=strconv.Atoi(limitstr)
	   pageno,_ := strconv.Atoi(pagenostr)
       
	    if limit < 1{
			limit=1
		}
		if pageno <1{
			pageno=1
		}
		offset := (pageno-1)*limit;

		rows ,err := db.Query("SELECT * FROM tasks LIMIT ? OFFSET ?",limit,offset)

		if err!=nil{
			log.Println("error in the data fetching")

		}

		var list []Task

		for rows.Next(){
			var task Task
			err := rows.Scan(&task.ID,&task.NAME,&task.STATUS)

			if err!=nil{
				log.Println("somthing went wrong in fetching the data")
			}
			list = append(list , task)
		}
		json.NewEncoder(w).Encode(list)

	}
}



func Insertmany(db *sql.DB) http.HandlerFunc{
	return func( w http.ResponseWriter, r *http.Request){
		tx, err := db.Begin()
		if err != nil{
			log.Fatal(err)
		}
		stmt,err := tx.Prepare("INSERT INTO tasks(name,status) VALUES(?,?)")
		if err!=nil{
			log.Fatal(err)
		}
		defer stmt.Close()

		for i:=1;i<=10000;i++{
			name:= "Task"+ strconv.Itoa(i)
			status := "pending"
			if i%2==0{
				status="Done"
			}
			_,err := stmt.Exec(name,status)
			if err!=nil{
				log.Println("error in the data inserting ")
			}
		}
  
		err = tx.Commit()
		if err!=nil{
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(map[string]string{
		"message":"the data inserted in the database succesfully",
	})

	}
}

func DeleteCompletedTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request){

		query :=`DELETE FROM tasks
		        WHERE status="DONE"    
		`

		_,err := db.Exec(query)
		
		 if err!=nil{
			log.Println("somthing went wrong  while deleting the task ")
			return 
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message":"completed task deleted succesfully ",
		})
	}
}


func DeleteTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

		id := r.URL.Query().Get("id")

		query :=`DELETE FROM tasks
		        WHERE id=?`

		_,err := db.Exec(query,id)
		
		if err!=nil{
			log.Println("error while deleting the task")
			return 
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message":"delete the task succesfully",
		})
	}
}


func ChangeStatus(db *sql.DB) http.HandlerFunc{
	return func( w http.ResponseWriter, r *http.Request){

		id := r.URL.Query().Get("id")

		query := `UPDATE tasks
		          SET status=CASE
				  when status="pending" THEN "DONE"
				  ELSE "pending"
				  END
				  WHERE id=?
				  `

	   _,err := db.Exec(query,id)
	   
	   if err!=nil{
		 log.Println("error in update  the status")
		 return 
	   }

	   json.NewEncoder(w).Encode(map[string]string{
		"message":"update the status succesfully ",
	   })
	}
}


func RenameTask(db *sql.DB) http.HandlerFunc{
	return func (w http.ResponseWriter , r *http.Request){
	

		id:= r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")

		query:=`UPDATE tasks
		        SET name=?
				WHERE  id=? `

		_,err:= db.Exec(query,name,id)

		if err!=nil{
			log.Println("error in updatting the data ")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message":"rename the task succesfully",
		})
	}
}

func gettask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request){
		var task Task;

	 id := r.URL.Query().Get("id")

	 query := `SELECT id,name,status FROM tasks where id=?`

	 err:= db.QueryRow(query,id).Scan(&task.ID,&task.NAME,&task.STATUS)

	 if err!=nil{
		log.Println("error in fetching the data");
		return 
	 }
	 json.NewEncoder(w).Encode(task)
	}
}

func GetAll(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		var list []Task;

		query := `SELECT * FROM tasks`

		res,err := db.Query(query);
		 if err!=nil{
			log.Println("wrong in fetching the data")
		 }
		 defer res.Close();

        for res.Next(){
			var task Task;

			err := res.Scan(&task.ID,&task.NAME,&task.STATUS)
			if err!=nil{
				log.Println("wrong in the scanning the data")
				return 
			}

			list = append(list,task)
		}  

		json.NewEncoder(w).Encode(list)
	}
}


func InsertTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
        var newtask Task;

		err := json.NewDecoder(r.Body).Decode(&newtask)

		if err!=nil{
		log.Println("error in fetching the data")
		}

		query := `INSERT INTO tasks (name , status) VALUES(?,?)`

		_, err = db.Exec(query,newtask.NAME,newtask.STATUS)

		if err != nil{
		log.Println("somthing went wrong to inserting the data ")
		return 
		}

		json.NewEncoder(w).Encode(map[string]string{
		"message":"the task inserted succesfully into database ",
		})

	}
}