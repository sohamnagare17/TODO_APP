
package handlers
 import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		
		"go-sqlite/models"
		"strconv"
  )

func ShowTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

	   limitstr:= r.URL.Query().Get("limit")
	   pagenostr:= r.URL.Query().Get("pageno")

	   limit,err:=strconv.Atoi(limitstr)
	     if err!=nil{
			 log.Println("plz provide valid limit value")
		  }

	   pageno,err1:= strconv.Atoi(pagenostr)
		  if err1!=nil{
			log.Println("plz provide valid after index value")
		  }

       
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

		var list []models.Task

		for rows.Next(){
			var task models.Task
			err := rows.Scan(&task.Id,&task.Name,&task.Status)

			if err!=nil{
				log.Println("somthing went wrong in fetching the data")
			}
			list = append(list , task)
		}

       nextpageno := pageno+1

	   response := map[string]interface{}{
		 "list of task":list,
		 "next_pageno":nextpageno,
		"has_more":  len(list) == limit,
		"currentpageno":pageno,
		"limit":limit,
	   }

		json.NewEncoder(w).Encode(response)

	}
}