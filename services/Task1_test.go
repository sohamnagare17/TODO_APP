package services

import (
	"testing"
     "go-sqlite/models"
	"go-sqlite/repository"
	"go-sqlite/testutils"
)

func Getservice() *TaskServices {
	db := testutils.SetupTestDb()
	repo := repository.NewTaskRepository(db)
	service := NewTaskServices(repo)
	return service

}

// gettaskuserbyid function
func TestGetTaskByUserId_Success(t *testing.T) {

	service := Getservice()
	tasks, err := service.GetTaskByUserId(
		"1", // userid
		"",  // status
		"",  // sortby
		"",  // order
		"",  // cursor
		"",  // limit
		"",  // page
	)
	// 5. Assertions
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) == 0 {
		t.Fatalf("expected tasks, got empty")
	}
}

func TestGetTaskByUserId_Withstatus(t *testing.T) {
	service := Getservice()
	tasks, err := service.GetTaskByUserId(
		"1",
		"pending", "", "", "", "", "",
	)
	if err != nil {
		t.Fatalf("unexpectd error:%v", err)
	}
	if len(tasks) == 0 {
		t.Fatalf("expected task, got empty:")
	}
	for _, task := range tasks {
		if task.Status != "pending" {
			t.Errorf("status is not pending: %v", task.Status)
		}
	}
}

func TestGetTaskByUserId_Pagination(t *testing.T) {
	tasks, err := Getservice().GetTaskByUserId(
		"1",
		"", "", "", "",
		"1", "1",
	)
	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if len(tasks) == 0 {
		t.Fatalf("expected tasks,got empty")
	}
}

func TestGetTaskByUserId_invaliduser(t *testing.T) {
	_, err := Getservice().GetTaskByUserId(
		"abc",
		"", "", "", "",
		"", "",
	)
	if err == nil {
		t.Fatalf("expected error for invalid user:%v", err)
	}
}

func TestGetTaskByUserId_emptyuser(t *testing.T) {
	_, err := Getservice().GetTaskByUserId(
		"",
		"", "", "", "",
		"", "",
	)
	if err == nil {
		t.Fatalf("expected error for invalid user:%v", err)
	}

}

func TestGetTaskByUserId_Nodata(t *testing.T) {
	tasks, err := Getservice().GetTaskByUserId(
		"999",
		"", "", "", "", "", "",
	)
	if err != nil {
		t.Fatalf("unexpected error :%v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expecting 0 tasks but got :%d", len(tasks))
	}
}

// inserttask function
func TestInsertTask_Succes(t *testing.T) {
	task := models.Task{
		UserId: 1,
		Name:   "new task",
		Status: "pending",
	}

	 err := Getservice().InsertTask(task)
	if err != nil {
		t.Errorf("unexpected error occur:%v", err)
	}
}

func TestInsertTask_EmptyName(t *testing.T){
	task := models.Task{
		UserId: 1,
		Name:   "",
		Status: "pending",
	}

	 err := Getservice().InsertTask(task)
	if err == nil {
		t.Errorf("expected error for empty task name :%v", err)
	}

}

func TestInsertTask_Invalidstatus(t *testing.T){
	task := models.Task{
		UserId: 1,
		Name:   "play football",
		Status: "invalid",
	}

	 err := Getservice().InsertTask(task)
	if err == nil {
		t.Errorf("expected error for empty task name :%v", err)
	}

}

func TestInsertTask_Emptystatus(t *testing.T){
	task := models.Task{
		UserId: 1,
		Name:   "football",
		Status: "",
	}

	 err := Getservice().InsertTask(task)
	if err != nil {
		t.Errorf("unexpected error occur :%v", err)
	}

}


