package services

import (
	"fmt"
	"go-sqlite/models"
	"go-sqlite/repository"
	"log"
	"strconv"
	"strings"
)

var validfields = map[string]bool{
	"name":      true,
	"createdAt": true,
	"updatedAt": true,
}

var validstatus = map[string]bool{
	"pending": true,
	"done":    true,
}

type TaskServices struct {
	repo repository.TaskRepo
}

type TaskService interface {
	InsertTask(task models.Task) error
	DeleteTask(idstr string, useridstr string) error
	UpdateTask(useridStr, taskidStr, name, status string) error
	GetTaskByUserId(useridstr, status, sortby, order, cursor, limitstr, pagenostr string) ([]models.Task, error)
}

func NewTaskServices(repo repository.TaskRepo) *TaskServices {
	return &TaskServices{repo: repo}
}

func (s *TaskServices) GetTaskByUserId(useridstr string, status string, sortby string, order string, cursor string, limitstr string, pagenostr string) ([]models.Task, error) {

	var err error
	limit := 5
	pageno := 1
	if limitstr != "" {
		parsedlimit, err := strconv.Atoi(limitstr)
		if err != nil {
			log.Println("plz provide valid limit", err)
			return nil, err
		}
		limit = parsedlimit
	}
	if pagenostr != "" {
		parsedpage, err := strconv.Atoi(pagenostr)
		if err != nil {
			log.Println("plz provide valid pageno", err)
			return nil, err
		}
		pageno = parsedpage
	}

	if useridstr == "" {
		log.Println("id required plz!")
		return nil, fmt.Errorf("id required plz")
	}

	userid, err := strconv.Atoi(useridstr)
	if err != nil {
		log.Println("id must be the number",err)
		return nil, err
	}
	query := `SELECT * FROM tasks1 WHERE userid=?`
	parameters := []interface{}{userid}

	if status != "" {
		query = query + " AND status=? "
		parameters = append(parameters, status)
	}

	if cursor != "" {
		query += " AND createdAt > ?"
		parameters = append(parameters, cursor)
	}

	if validfields[sortby] {
		query = query + " ORDER BY " + sortby

		if order == "ASC" || order == " asc " {
			query += " ASC "
		} else if order == " desc" || order == " DESC" {
			query += " DESC"
		} else {
			query += " DESC"
		}
	} else {
		query += " ORDER BY createdAt DESC"
	}
	//pagination
	if cursor != "" {
		query += " LIMIT ? "
		parameters = append(parameters, limit)
	} else {
		offset := (pageno - 1) * limit
		query += " LIMIT ? OFFSET ?"
		parameters = append(parameters, limit, offset)
	}
	log.Println("Query:", query)
	log.Println("Values:", parameters)
	return s.repo.GetTaskByUserId(query, parameters)
}

func (s *TaskServices) InsertTask(newtask models.Task) error {

	log.Println(newtask.UserId)
	newtask.Name = strings.TrimSpace(newtask.Name)
	if newtask.Name == "" {

		log.Println("Enter a task")
		return fmt.Errorf("enter task name!")
	}

	newtask.Status = strings.ToLower(strings.TrimSpace(newtask.Status))
	if newtask.Status == "" {
		newtask.Status = "pending"
	} else if !validstatus[newtask.Status] {

		log.Println("Invalid status(done/pending only allowed)")
		return fmt.Errorf("invalid status ")
	}
	return s.repo.InsertTask(newtask)
}

func (s *TaskServices) DeleteTask(idstr string, useridstr string) error {
	if idstr == "" || useridstr == "" {
		log.Println("userid and taskid  required plz provide ids")
		return fmt.Errorf("userid and taskid requried")
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("id must be integer", err)
		return err
	}

	userid, err1 := strconv.Atoi(useridstr)
	if err1 != nil {
		log.Println("userid must be integer", err1)
		return err1
	}

	rows, err := s.repo.DeleteTask(id, userid)
	if err != nil {
		log.Println("error while executing the database query", err)
		return err
	}
	if rows==0{
		return fmt.Errorf("Task not found")
	}
	return nil

}

func (s *TaskServices) UpdateTask(useridStr, taskidStr, name, status string) error {

	if useridStr == "" || taskidStr == "" {
		return fmt.Errorf("userid and taskid required")
	}

	uid, err := strconv.Atoi(useridStr)
	if err != nil || uid <= 0 {
		return fmt.Errorf("invalid userid")
	}

	tid, err := strconv.Atoi(taskidStr)
	if err != nil || tid <= 0 {
		return fmt.Errorf("invalid taskid")
	}

	if strings.TrimSpace(name) == "" && name != "" {
		return fmt.Errorf("name should not be empty")
	}

	if name == "" && status == "" {
		return fmt.Errorf("nothing to update")
	}

	rows, err := s.repo.UpdateTask(uid, tid, name, status)
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}



