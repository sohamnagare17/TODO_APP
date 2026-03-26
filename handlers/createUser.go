package handlers

import (
	"encoding/json"
	db "go-sqlite/database"
	"go-sqlite/models"
	"log"
	"net/http"
)

func CreateUser(w http.ResponseWriter,r *http.Request){

	var user models.User
	err:=json.NewDecoder(r.Body).Decode(&user)
	if err!=nil{
		log.Println("error in fetching the data")
	}
	query:=`INSERT INTO Users(username,email) VALUES(?,?) `

	_,err=db.DB.Exec(query,user.UserName,user.Email)
	if err!=nil{
		log.Println("Error while inser the User data")
	}
	w.Write([]byte("User Created Successfully"))

}