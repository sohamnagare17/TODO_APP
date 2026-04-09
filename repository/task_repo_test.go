package repository

import (
	"go-sqlite/models"
	"go-sqlite/testutils"
	"testing"
)

func TestInsertTask_Success(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	err := repo.InsertTask(models.Task{
		Name:   "New Task",
		Status: "pending",
		UserId: 1, 
	})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
func TestInsertTask_InvalidUser(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	err := repo.InsertTask(models.Task{
		Name:   "Task",
		Status: "pending",
		UserId: 10,
	})
	if err == nil {
		t.Errorf("expected foreign key error")
	}
}
func TestInsertTask_Duplicate(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	task := models.Task{
		Name:   "Task 1",
		Status: "pending",
		UserId: 1,
	}
	err := repo.InsertTask(task)
	if err==nil{
		t.Errorf("Expected duplicate error")
	}
}
func TestInsertTask_EmptyFields(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	err := repo.InsertTask(models.Task{
		Name:"",
		Status: "",
		UserId: 1,
	})
	if err == nil {
		t.Log("no error - DB allows empty fields")
	}
}


func TestDeleteTask_Success(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	rows, err := repo.DeleteTask(1, 1) 
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rows != 1 {
		t.Errorf("expected 1 row deleted got %d", rows)
	}
}
func TestDeleteTask_NotFound(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	_, err := repo.DeleteTask(999, 1) 
	if err == nil {
		t.Errorf("expected error when task not found")
	}
}
func TestDeleteTask_WrongUser(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	_, err := repo.DeleteTask(1, 2)
	if err == nil {
		t.Errorf("expected error due to wrong user")
	}
}
func TestDeleteTask_DBClosed(t *testing.T) {
	db := testutils.SetupTestDb()
	repo := NewTaskRepository(db)
	db.Close()
	_, err := repo.DeleteTask(1, 1)
	if err == nil {
		t.Errorf("expected error when DB is closed")
	}
}
func TestDeleteTask_DeleteTwice(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	_, err := repo.DeleteTask(1, 1)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	_, err = repo.DeleteTask(1, 1)
	if err == nil {
		t.Errorf("expected error on second delete")
	}
}
func TestDeleteTask_TableNotExist(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	db.Exec("DROP TABLE tasks1")
	repo := NewTaskRepository(db)
	_, err := repo.DeleteTask(1, 1)
	if err == nil {
		t.Errorf("expected error when table does not exist")
	}
}


func TestUpdateTask_BothFields(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	rows, err := repo.UpdateTask(1, 1, "Updated Task", "done")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rows != 1 {
		t.Errorf("expected 1 row updated got %d", rows)
	}
}
func TestUpdateTask_VerifyBoth(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	_, err := repo.UpdateTask(1, 1, "Updated Task", "done")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	var name, status string
	db.QueryRow("SELECT name, status FROM tasks1 WHERE id=? AND userid=?", 1, 1).
		Scan(&name, &status)

	if name != "Updated Task" || status != "done" {
		t.Errorf("update not applied correctly")
	}
}
func TestUpdateTask_OnlyName(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	rows, err := repo.UpdateTask(1, 1, "OnlyName", "")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rows != 1 {
		t.Errorf("expected 1 row updated")
	}
}
func TestUpdateTask_OnlyStatus(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	rows, err := repo.UpdateTask(1, 1, "", "done")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rows != 1 {
		t.Errorf("expected 1 row updated")
	}
}
func TestUpdateTask_NothingToUpdate(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	_, err := repo.UpdateTask(1, 1, "", "")
	if err == nil {
		t.Errorf("expected error when nothing to update")
	}
}
func TestUpdateTask_DBClosed(t *testing.T) {
	db := testutils.SetupTestDb()
	repo := NewTaskRepository(db)
	db.Close()
	_, err := repo.UpdateTask(1, 1, "New", "done")
	if err == nil {
		t.Errorf("expected error when DB is closed")
	}
}

func TestGetTaskByUserId_Success(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	query := "SELECT * FROM tasks1 WHERE userid=?"
	params := []interface{}{1}
	tasks, err := repo.GetTaskByUserId(query, params)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(tasks) == 0 {
		t.Errorf("expected tasks but got empty")
	}
}
func TestGetTaskByUserId_DataCheck(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()

	repo := NewTaskRepository(db)

	query := "SELECT * FROM tasks1 WHERE userid=?"
	params := []interface{}{1}

	tasks, err := repo.GetTaskByUserId(query, params)

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	for _, task := range tasks {
		if task.UserId != 1 {
			t.Errorf("expected userid 1 got %d", task.UserId)
		}
	}
}
func TestRepo_GetTaskByUserId_Empty(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	query := "SELECT * FROM tasks1 WHERE userid=?"
	params := []interface{}{999}
	tasks, err := repo.GetTaskByUserId(query, params)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty result")
	}
}
func TestRepo_GetTaskByUserId_InvalidQuery(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewTaskRepository(db)
	query := "SELECT * FROM table2 WHERE userid=?"
	params := []interface{}{1}
	_, err := repo.GetTaskByUserId(query, params)

	if err == nil {
		t.Errorf("expected query error")
	}
}
func TestRepo_GetTaskByUserId_StatusFilter(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()

	repo := NewTaskRepository(db)

	query := "SELECT * FROM tasks1 WHERE userid=? AND status=?"
	params := []interface{}{1, "pending"}

	tasks, err := repo.GetTaskByUserId(query, params)

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	for _, task := range tasks {
		if task.Status != "pending" {
			t.Errorf("expected pending got %s", task.Status)
		}
	}
}
func TestRepo_GetTaskByUserId_Limit(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()

	repo := NewTaskRepository(db)

	query := "SELECT * FROM tasks1 WHERE userid=? LIMIT ?"
	params := []interface{}{1, 1}

	tasks, err := repo.GetTaskByUserId(query, params)

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("expected 1 task got %d", len(tasks))
	}
}
func TestRepo_GetTaskByUserId_DBClosed(t *testing.T) {
	db := testutils.SetupTestDb()

	repo := NewTaskRepository(db)
	db.Close()

	query := "SELECT * FROM tasks1 WHERE userid=?"
	params := []interface{}{1}

	_, err := repo.GetTaskByUserId(query, params)

	if err == nil {
		t.Errorf("expected error when DB is closed")
	}
}
