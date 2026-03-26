package handlers

import (
	"net/http"
	"strconv"
)

func GetTask(w http.ResponseWriter, r *http.Request) {

	id,:= r.URL.Query().Get("userid")
	iid := strconv.Atoi(id)

	query := `SELECT * FROM tasks 
			WHERE userId=?
			ORDER BY createdAt`
		
	rows,err:=db.DB.Query(query)
	var task models.Task

	rows,err:=d


}
