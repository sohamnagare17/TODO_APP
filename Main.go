package main

import (
	
	"net/http"
	"go-sqlite/routes"
	"log"
	"os"
	"go-sqlite/database"

	
)



func main(){
	file , err := os.OpenFile("app.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)

	if err!=nil{
		log.Fatal(err)
	}
     defer file.Close()

    log.SetOutput(file)

	dbconn := db.Dbinit()

	defer dbconn.Close()

	 routes.SetupRoutes(dbconn)

	log.Println("server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
