package services

import (
	"errors"
	"go-sqlite/models"
	"testing"
)

type FakeTaskRepo struct {
	tasks []models.Task
	err   error
	rows  int64
}

func (f *FakeTaskRepo) InsertTask(newtask models.Task) error {
	if f.err != nil {
		return f.err
	}
	f.tasks = append(f.tasks, newtask)
	return nil
}
func (f *FakeTaskRepo) DeleteTask(taskid int, userid int) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.rows, nil
}
func (f *FakeTaskRepo) UpdateTask(userid, taskid int, name, status string) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.rows, nil
}
func (f *FakeTaskRepo) GetTaskByUserId(query string, params []interface{}) ([]models.Task, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.tasks, nil
}

func TestInsertTask_Success(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	task := models.Task{
		Name:   "Learn Go-lang",
		Status: "pending",
	}
	err := service.InsertTask(task)
	if err != nil {
		t.Errorf("Unexpected error got %v", err)
	}
	if len(mockrepo.tasks) != 1 {
		t.Errorf("Expected 1 task got %d", len(mockrepo.tasks))
	}
}
func TestInsertTask_Emptyname(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	task := models.Task{
		Name:   "",
		Status: "pending",
	}
	err := service.InsertTask(task)
	if err == nil {
		t.Errorf("Expected error")
	}
}
func TestInsertTask_EmptyFields(t *testing.T) {

	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	task := models.Task{}
	err := service.InsertTask(task)
	if err == nil {
		t.Errorf("Expected error")
	}
}
func TestInsertTask_EmptyStatus(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	task := models.Task{
		Name:   "learning java",
		Status: "",
	}
	err := service.InsertTask(task)
	if err != nil {
		t.Errorf("Unexpected Error got %v", err)
	}
	if len(mockrepo.tasks) != 1 {
		t.Errorf("Expected 1 task got %d", len(mockrepo.tasks))
	}
}
func TestInsertTask_RepoError(t *testing.T) {

	mockrepo := &FakeTaskRepo{
		err: errors.New("DB Error"),
	}
	service := NewTaskServices(mockrepo)
	task := models.Task{
		Name:   "Learning cpp",
		Status: "pending",
	}
	err := service.InsertTask(task)
	if err == nil {
		t.Errorf("Expected error from repo got nil")
	}

}
func TestInsertTask_SpacesInNames(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	task := models.Task{
		Name:   "       ",
		Status: "pending",
	}
	err := service.InsertTask(task)
	if err == nil {
		t.Errorf("expected Error got nil")
	}
}

func TestDeleteTask_Success(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		rows: 1,
	}
	service := NewTaskServices(mockrepo)
	err := service.DeleteTask("10", "1")
	if err != nil {
		t.Errorf("Unexpected error got %v", err)
	}
}
func TestDeleteTask_EmptyTaskId(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.DeleteTask("", "1")
	if err == nil {
		t.Errorf("Expected error from empty TaskId")
	}
}
func TestDeleteTask_EmptyUserId(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.DeleteTask("10", "")
	if err == nil {
		t.Errorf("Expected Error from empty UserId")
	}
}
func TestDeleteTask_BothEmpty(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)

	err := service.DeleteTask("", "")

	if err == nil {
		t.Errorf("Expected error when both ids are empty")
	}
}
func TestDeleteTask_No_Task(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		rows: 0,
	}
	service := NewTaskServices(mockrepo)
	err := service.DeleteTask("1", "10")
	if err == nil {
		t.Errorf("Expected error when no task found, got nil")
	}
}
func TestDeleteTask_InvalidTaskid(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.DeleteTask("abc", "1")
	if err == nil {
		t.Errorf("Expected error got nil")
	}
}
func TestDeleteTask_InvalidUserId(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.DeleteTask("5", "xyz")
	if err == nil {
		t.Errorf("Expected error got nil")
	}
}

func TestUpdateTask_Success(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		rows: 1,
	}
	service := NewTaskServices(mockrepo)
	err := service.UpdateTask("1", "1", "newtask", "done")
	if err != nil {
		t.Errorf("Unexpected error got %v", err)
	}
}
func TestUpdateTask_NoRows(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		rows: 0,
	}
	service := NewTaskServices(mockrepo)
	err := service.UpdateTask("1", "1", "newtask", "done")
	if err == nil {
		t.Errorf("Expected error when no rows updated")
	}
}
func TestUpdateTask_EmptyName(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.UpdateTask("1", "1", "", "done")
	if err == nil {
		t.Errorf("Expected error got nil")
	}
}
func TestUpdateTask_InvalidUserId(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.UpdateTask("abc", "1", "task", "done")
	if err == nil {
		t.Errorf("Expected error for invalid user id")
	}
}
func TestUpdateTask_InvalidTaskId(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.UpdateTask("1", "xyz", "task", "done")
	if err == nil {
		t.Errorf("Expected error for invalid task id")
	}
}
func TestUpdateTask_EmptyIds(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	err := service.UpdateTask("", "", "task", "done")
	if err == nil {
		t.Errorf("Expected error for empty ids")
	}
}
func TestUpdateTask_EmptyFields(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)

	err := service.UpdateTask("1", "1", "", "")

	if err == nil {
		t.Errorf("Expected error when nothing to update")
	}
}

func TestGetTaskByUserID_Success(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		tasks: []models.Task{
			{Name: "Task1", Status: "pending"},
		},
	}
	service := NewTaskServices(mockrepo)
	result, err := service.GetTaskByUserId("1", "", "", "", "", "", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 task, got %d", len(result))
	}
}
func TestGetTaskByUserId_EmptyUserId(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)

	_, err := service.GetTaskByUserId("", "", "", "", "", "", "")

	if err == nil {
		t.Errorf("Expected error for empty user id")
	}
}

func TestGetTaskByUserId_InvalidUserId(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)

	_, err := service.GetTaskByUserId("abc", "", "", "", "", "", "")

	if err == nil {
		t.Errorf("Expected error for invalid user id")
	}
}
func TestGetTaskByUserId_InvalidLimit(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)

	_, err := service.GetTaskByUserId("1", "", "", "", "", "abc", "")

	if err == nil {
		t.Errorf("Expected error for invalid limit")
	}
}
func TestGetTaskByUserId_InvalidPageNo(t *testing.T) {
	mockrepo := &FakeTaskRepo{}
	service := NewTaskServices(mockrepo)
	_, err := service.GetTaskByUserId("1", "", "", "", "", "", "xyz")
	if err == nil {
		t.Errorf("Expected error for invalid page number")
	}
}
func TestGetTaskByUserId_RepoError(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		err: errors.New("DB Error"),
	}
	service := NewTaskServices(mockrepo)
	_, err := service.GetTaskByUserId("1", "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Expected repo error")
	}
}
func TestGetTaskByUserId_WithStatus(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		tasks: []models.Task{
			{Name: "Task1", Status: "done"},
		},
	}
	service := NewTaskServices(mockrepo)
	result, err := service.GetTaskByUserId("1", "done", "", "", "", "", "")
	if err != nil {
		t.Errorf("Unexpected error")
	}
	if result[0].Status != "done" {
		t.Errorf("Expected status done")
	}
}
func TestGetTaskByUserId_WithPagination(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		tasks: []models.Task{
			{Name: "Task1"},
		},
	}
	service := NewTaskServices(mockrepo)
	_, err := service.GetTaskByUserId("1", "", "", "", "", "10", "2")
	if err != nil {
		t.Errorf("Unexpected error")
	}
}
func TestGetTaskByUserId_WithCursor(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		tasks: []models.Task{
			{Name: "Task1"},
		},
	}
	service := NewTaskServices(mockrepo)
	_, err := service.GetTaskByUserId("1", "", "", "", "2024-01-01", "", "")
	if err != nil {
		t.Errorf("Unexpected error")
	}
}
func TestGetTaskByUserId_Sorting(t *testing.T) {
	mockrepo := &FakeTaskRepo{
		tasks: []models.Task{
			{Name: "Task1"},
		},
	}
	service := NewTaskServices(mockrepo)
	_, err := service.GetTaskByUserId("1", "", "createdAt", "ASC", "", "", "")
	if err != nil {
		t.Errorf("Unexpected error")
	}
}
