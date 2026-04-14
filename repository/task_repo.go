package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-sqlite/metrics"
	"go-sqlite/models"
	"go.opentelemetry.io/otel"
	"log"
	"time"
)

type TaskRepository struct {
	db *sql.DB
}

type TaskRepo interface {
	GetTaskByUserId(ctx context.Context, query string, params []interface{}) ([]models.Task, error)
	InsertTask(newtask models.Task) error
	DeleteTask(id int, userid int) (int64, error)
	UpdateTask(userid, taskid int, name, status string) (int64, error)
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetTaskByUserId(ctx context.Context, query string, params []interface{}) ([]models.Task, error) {

	tracer := otel.Tracer("task-repo")
	ctx, span := tracer.Start(ctx, "get repo")
	defer span.End()

	start := time.Now()
	var tasklist []models.Task
	rows, err := r.db.Query(query, params...)

	duration := time.Since(start).Seconds()

	if err != nil {
		log.Println("error in execution the query", err)
		metrics.DBErrorsTotal.WithLabelValues("select", "tasks1").Inc()
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
	metrics.DBQueryDuration.WithLabelValues("select", "tasks1").Observe(duration)
	return tasklist, nil
}

func (r *TaskRepository) InsertTask(newtask models.Task) error {

	query := `INSERT INTO tasks1 (name ,status,userid,createdAt,updatedAt) VALUES(?,?,?,?,?)`

	now := time.Now().UTC().Format(time.RFC3339)
	start := time.Now()
	_, err := r.db.Exec(query, newtask.Name, newtask.Status, newtask.UserId, now, now)

	duration := time.Since(start).Seconds()
	if err != nil {
		log.Println("somthing went wrong to inserting the data ", err)
		metrics.DBErrorsTotal.WithLabelValues("insert", "tasks1").Inc()
		return err
	}
	metrics.DBQueryDuration.WithLabelValues("insert", "tasks1").Observe(duration)
	return nil
}

func (r *TaskRepository) DeleteTask(id int, userid int) (int64, error) {
	query := `DELETE FROM tasks1 WHERE userid=? AND id=?`
	start := time.Now()
	result, err := r.db.Exec(query, userid, id)
	duration := time.Since(start).Seconds()
	if err != nil {
		log.Println("error while executing the database query", err)
		metrics.DBErrorsTotal.WithLabelValues("delete", "tasks1").Inc()
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
	metrics.DBQueryDuration.WithLabelValues("delete", "tasks1").Observe(duration)
	return rowsAffected, nil
}

func (r *TaskRepository) UpdateTask(userid, taskid int, name, status string) (int64, error) {

	var query string
	var res sql.Result
	var err error
	start := time.Now()
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

	duration := time.Since(start).Seconds()
	if err != nil {
		log.Println("error updating task:", err)
		metrics.DBErrorsTotal.WithLabelValues("update", "tasks1").Inc()
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	metrics.DBQueryDuration.WithLabelValues("update", "tasks1").Observe(duration)
	return rows, nil
}
