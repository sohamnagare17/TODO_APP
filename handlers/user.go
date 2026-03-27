package handlers

import (
	"encoding/json"
	"net/http"
	"log"
	"database/sql"
	"go-sqlite/models"
	"strconv"
)

func InsertUser(db *sql.DB) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request){

		var user models.Users;

		err := json.NewDecoder(request.Body).Decode(&user)

		if err!=nil{
			log.Println("error in decoding the data",err)
			return 
		}

		if user.Username==""|| user.Email==""{
			log.Println("username and email required ")
			return;
		}

		query := `INSERT INTO users(username,email) VALUES(?,?) `

		_,err = db.Exec(query,user.Username,user.Email)

		if err!=nil{
			log.Println("error while inserting the user ",err)
			http.Error(writer,"email already exists",409)
			return 
		}
		json.NewEncoder(writer).Encode(map[string]interface{}{
			"username":user.Username,
			"message":"This user added succesfully",
		})
	}
}
func GetAllUsers(db *sql.DB) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request){
		var userlist []models.Users;

		query := `SELECT * FROM users`;

		rows,err := db.Query(query)

		if err!=nil{
			log.Println("error in fetching the data from the database",err)
			return 
		}

		for rows.Next(){
			var user models.Users

			err = rows.Scan(&user.Userid,&user.Email,&user.Username)
			if err!=nil{
				log.Println("error in scanning the data from the rows",err)
			}
			userlist = append(userlist, user)
		}

		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":"list of the user as follows",
			"userlist":userlist,
		})

	
	}
}
func GetUserById(db *sql.DB) http.HandlerFunc{
	return func(writer http.ResponseWriter , request *http.Request){
		var user models.Users

		idstr := request.URL.Query().Get("userid")

		if idstr ==""{
			log.Println("Id required")
			return
		}

		id ,err := strconv.Atoi(idstr)
		if err!=nil{
			log.Println("id must be number")
			return 
		}

		query := `SELECT * FROM users WHERE user_id=?`

		err1:=db.QueryRow(query,id).Scan(&user.Userid,&user.Username,&user.Email)

		if err1!=nil{
			log.Println("error in the fetching the data")
			return 
		}
		json.NewEncoder(writer).Encode(map[string]interface{}{
			"message":"the user is ",
			"username":user.Username,
			"userid":user.Userid,
			"useremail":user.Email,
		})

		
	}
}