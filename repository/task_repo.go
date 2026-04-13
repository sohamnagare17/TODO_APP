package repository

import (
	"database/sql"
	"fmt"
	"go-sqlite/models"
	"log"
	"time"
	"context"
	"go.opentelemetry.io/otel"
)

type TaskRepository struct {
	db *sql.DB
}

type TaskRepo interface {
	GetTaskByUserId(ctx context.Context,query string, params []interface{}) ([]models.Task, error)
	InsertTask(newtask models.Task) error
	DeleteTask(id int, userid int) (int64, error)
	UpdateTask(userid, taskid int, name, status string) (int64, error)
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetTaskByUserId(ctx context.Context,query string, params []interface{}) ([]models.Task, error) {


	 tracer := otel.Tracer("task-repo")
	 ctx , span := tracer.Start(ctx,"get repo")
	 defer span.End()

	var tasklist []models.Task
	rows, err := r.db.Query(query, params...)
	if err != nil {
		log.Println("error in execution the query", err)
		return nil, err
	}
	for rows.Next() {
		var task models.Task

		err = rows.Scan(
			&task.Id,
			&task.Name,
			&task.Status,
			&task.UserId,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			log.Println("error in scanning the  data", err)
			return nil, err
		}
		tasklist = append(tasklist, task)
	}
	return tasklist, nil
}

func (r *TaskRepository) InsertTask(newtask models.Task) error {

	query := `INSERT INTO tasks1 (name ,status,userid,createdAt,updatedAt) VALUES(?,?,?,?,?)`

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(query, newtask.Name, newtask.Status, newtask.UserId, now, now)

	if err != nil {
		log.Println("somthing went wrong to inserting the data ", err)
		return err
	}
	return nil
}

func (r *TaskRepository) DeleteTask(id int, userid int) (int64, error) {
	query := `DELETE FROM tasks1 WHERE userid=? AND id=?`

	result, err := r.db.Exec(query, userid, id)
	if err != nil {
		log.Println("error while executing the database query", err)
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Println("error in checking rows affected", err)
		return 0, err
	}

	if rowsAffected == 0 {
		log.Println("task not found ")
		return 0, fmt.Errorf("task is not found")
	}
	return rowsAffected, nil
}

func (r *TaskRepository) UpdateTask(userid, taskid int, name, status string) (int64, error) {

	var query string
	var res sql.Result
	var err error

	switch {
	case name != "" && status != "":
		query = `UPDATE tasks1 
			SET name=?, status=?, updatedAt=CURRENT_TIMESTAMP
			WHERE id=? AND userid=?`
		res, err = r.db.Exec(query, name, status, taskid, userid)

	case name != "":
		query = `UPDATE tasks1 
			SET name=?, updatedAt=CURRENT_TIMESTAMP
			WHERE id=? AND userid=?`
		res, err = r.db.Exec(query, name, taskid, userid)

	case status != "":
		query = `UPDATE tasks1 
			SET status=?, updatedAt=CURRENT_TIMESTAMP
			WHERE id=? AND userid=?`
		res, err = r.db.Exec(query, status, taskid, userid)

	default:
		return 0, fmt.Errorf("nothing to update")
	}

	if err != nil {
		log.Println("error updating task:", err)
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}
