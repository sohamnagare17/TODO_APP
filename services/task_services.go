package services

import (
	"context"
	"fmt"
	"go-sqlite/models"
	"go-sqlite/Redis"
	"go-sqlite/repository"
	"go.opentelemetry.io/otel"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
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
	rdb  *redis.Client
}

type TaskService interface {
	InsertTask(task models.Task) error
	DeleteTask(idstr string, useridstr string) error
	UpdateTask(useridStr, taskidStr, name, status string) error
	GetTaskByUserId(ctx context.Context, useridstr, status, sortby, order, cursor, limitstr, pagenostr string) ([]models.Task, error)
}

func NewTaskServices(repo repository.TaskRepo,rdb *redis.Client) *TaskServices {
	return &TaskServices{
		repo: repo,
	    rdb : rdb ,     
	}
}

func (s *TaskServices) GetTaskByUserId(ctx context.Context, useridstr string, status string, sortby string, order string, cursor string, limitstr string, pagenostr string) ([]models.Task, error) {

	tracer := otel.Tracer("task-service")
	ctx, span := tracer.Start(ctx, "getservice")

	defer span.End()

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
		log.Println("id must be the number", err)
		return nil, err
	}

   //here add cache logic 
   
	key := fmt.Sprintf(
		"tasks:user:%d:status:%s:sort:%s:order:%s:cursor:%s:limit:%d:page:%d",
		userid, status, sortby, order, cursor, limit, pageno,
	)

	var tasks []models.Task

	
	if Redis.GetCache(s.rdb, key, &tasks) {
		log.Println("cache hit ")
		return tasks, nil
	}

	log.Println("cache miss  db hit")


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
	//return s.repo.GetTaskByUserId(ctx, query, parameters)
	tasks,err = s.repo.GetTaskByUserId(ctx, query, parameters)
	if err != nil{
		return nil,err
	}

	Redis.SetCache(s.rdb, key, tasks, time.Minute*5)
	return tasks,nil
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
	if rows == 0 {
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
