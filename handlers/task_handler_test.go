package handlers

import (
	"fmt"
	"go-sqlite/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	err   error
	tasks []models.Task
}

func (m *MockService) InsertTask(newtask models.Task) error {
	return m.err
}
func (m *MockService) DeleteTask(idstr string, useridstr string) error {
	return m.err
}
func (m *MockService) UpdateTask(useridStr, taskidStr, name, status string) error {
	return m.err
}
func (m *MockService) GetTaskByUserId(useridstr string, status string, sortby string, order string, cursor string, limitstr string, pagenostr string) ([]models.Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tasks, nil
}

func TestInsertTask_Success(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	body := `{"name":"learn c","status":"pending"}`
	req := httptest.NewRequest("POST", "/users/{userid}/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.InsertTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", w.Code)
	}
}
func TestInsertTask_InvalidJson(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	body := `invalid json`
	req := httptest.NewRequest("POST", "/users/{userid}/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.InsertTask(w, req)

}
func TestInsertTask_EmptyBody(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("POST", "/users/1/tasks", nil)
	req.SetPathValue("userid", "1")

	w := httptest.NewRecorder()

	handler.InsertTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}
func TestInsertTask_MissingUserID(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	body := `{"name":"learn c","status":"pending"}`
	req := httptest.NewRequest("POST", "/users//tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.InsertTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}
func TestInsertTask_ServiceError(t *testing.T) {
	service := &MockService{
		err: fmt.Errorf("failed to insert"),
	}
	handler := NewTaskHandler(service)
	body := `{"name":"learn c","status":"pending"}`
	req := httptest.NewRequest("POST", "/users/1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.InsertTask(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 got %d", w.Code)
	}
}
func TestInsertTask_EmptyName(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	body := `{"name":"","status":"pending"}`
	req := httptest.NewRequest("POST", "/users/1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("userid", "1")

	w := httptest.NewRecorder()

	handler.InsertTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}
func TestInsertTask_EmptyStatus(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	body := `{"name":"learn c","status":""}`
	req := httptest.NewRequest("POST", "/users/1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.InsertTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	body := ``
	req := httptest.NewRequest("DELETE", "/users/1/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-type", "application/json")
	req.SetPathValue("userid", "1")
	req.SetPathValue("taskid", "1")
	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}
func TestDeleteTask_ServiceError(t *testing.T) {
	service := &MockService{
		err: fmt.Errorf("delete failed"),
	}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("DELETE", "/users/1/tasks/1", nil)
	req.SetPathValue("userid", "1")
	req.SetPathValue("taskid", "1")

	w := httptest.NewRecorder()

	handler.DeleteTask(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 got %d", w.Code)
	}
}
func TestDeleteTask_MissingUserID(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("DELETE", "/users//tasks/1", nil)
	req.SetPathValue("taskid", "1")

	w := httptest.NewRecorder()

	handler.DeleteTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestDeleteTask_MissingTaskID(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	req := httptest.NewRequest("DELETE", "/users/1/tasks/", nil)
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestDeleteTask_BothIDsMissing(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("DELETE", "/users//tasks/", nil)

	w := httptest.NewRecorder()

	handler.DeleteTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestDeleteTask_InvalidIDs(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("DELETE", "/users/abc/tasks/xyz", nil)
	req.SetPathValue("userid", "abc")
	req.SetPathValue("taskid", "xyz")

	w := httptest.NewRecorder()

	handler.DeleteTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	body := `{"name":"new task","status":"done"}`

	req := httptest.NewRequest("PATCH", "/users/1/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("userid", "1")
	req.SetPathValue("taskid", "1")

	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}
func TestUpdateTask_InvalidMethod(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("GET", "/users/1/tasks/1", nil)

	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 got %d", w.Code)
	}
}
func TestUpdateTask_InvalidJSON(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	body := `invalid json`

	req := httptest.NewRequest("PATCH", "/users/1/tasks/1", strings.NewReader(body))
	req.SetPathValue("userid", "1")
	req.SetPathValue("taskid", "1")

	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestUpdateTask_EmptyBody(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("PATCH", "/users/1/tasks/1", nil)
	req.SetPathValue("userid", "1")
	req.SetPathValue("taskid", "1")

	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestUpdateTask_ServiceError(t *testing.T) {
	service := &MockService{
		err: fmt.Errorf("update failed"),
	}
	handler := NewTaskHandler(service)

	body := `{"name":"new task","status":"done"}`

	req := httptest.NewRequest("PATCH", "/users/1/tasks/1", strings.NewReader(body))
	req.SetPathValue("userid", "1")
	req.SetPathValue("taskid", "1")

	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestUpdateTask_MissingIDs(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	body := `{"name":"new task","status":"done"}`
	req := httptest.NewRequest("PATCH", "/users//tasks/", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}

func TestGetTaskByUserId_Success(t *testing.T) {
	service := &MockService{
		tasks: []models.Task{
			{Name: "task1", Status: "pending"},
		},
	}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("GET", "/users/1/tasks?status=pending", nil)
	req.SetPathValue("userid", "1")

	w := httptest.NewRecorder()

	handler.GetTaskByUserId(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}
func TestGetTaskByUserId_ServiceError(t *testing.T) {
	service := &MockService{
		err: fmt.Errorf("failed to fetch"),
	}
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("GET", "/users/1/tasks", nil)
	req.SetPathValue("userid", "1")

	w := httptest.NewRecorder()

	handler.GetTaskByUserId(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}

func TestGetTaskByUserId_WithQueryParams(t *testing.T) {
	service := &MockService{
		tasks: []models.Task{
			{Name: "task1"},
		},
	}
	handler := NewTaskHandler(service)
	req := httptest.NewRequest(
		"GET",
		"/users/1/tasks?status=done&sortby=createdAt&order=desc&limit=10&pageno=2",
		nil,
	)
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.GetTaskByUserId(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}

func TestGetTaskByUserId_EmptyResult(t *testing.T) {
	service := &MockService{
		tasks: []models.Task{},
	}
	handler := NewTaskHandler(service)
	req := httptest.NewRequest("GET", "/users/1/tasks", nil)
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.GetTaskByUserId(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}
func TestGetTaskByUserId_MissingUserID(t *testing.T) {
	service := &MockService{}
	handler := NewTaskHandler(service)
	req := httptest.NewRequest("GET", "/users//tasks", nil)
	w := httptest.NewRecorder()
	handler.GetTaskByUserId(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
