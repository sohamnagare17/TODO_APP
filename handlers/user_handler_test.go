package handlers

import (
	"fmt"
	"go-sqlite/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FakeService struct{
	err error
	users []models.Users
}

func (m *FakeService) InsertUser(newuser models.Users) error {
	return m.err
}
func(m *FakeService)GetUserById(idstr string) (models.Users, error){
	if m.err!=nil{
		return models.Users{},m.err
	}
	return models.Users{},nil
}
func(m *FakeService)GetAllUsers() ([]models.Users, error){
	if m.err!=nil{
		return nil,m.err
	}
	return m.users,nil
}

func TestInsertUser_Success(t *testing.T){
	service:=&FakeService{}
	handler:=NewUserHandler(service)
	body := `{"name":"Siddharth","email":"test@gmail.com"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w:=httptest.NewRecorder()
	handler.InsertUser(w,req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}
func TestInsertUser_InvalidMethod(t *testing.T) {
	service := &FakeService{}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	handler.InsertUser(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 got %d", w.Code)
	}
}
func TestInsertUser_InvalidJSON(t *testing.T) {
	service := &FakeService{}
	handler := NewUserHandler(service)
	body := `invalid json data`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.InsertUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestInsertUser_EmptyBody(t *testing.T) {
	service := &FakeService{}
	handler := NewUserHandler(service)

	req := httptest.NewRequest("POST", "/user", nil)

	w := httptest.NewRecorder()

	handler.InsertUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}
func TestInsertUser_ServiceError(t *testing.T) {
	service := &FakeService{
		err: fmt.Errorf("insert failed"),
	}
	handler := NewUserHandler(service)
	body := `{"name":"Siddharth","email":"sid@gmail.com"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.InsertUser(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 got %d", w.Code)
	}
}
func TestInsertUser_EmptyFields(t *testing.T) {
	service := &FakeService{}
	handler := NewUserHandler(service)
	body := `{"name":"","email":""}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler.InsertUser(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}


func TestGetAllUsers_Success(t *testing.T) {
	service := &FakeService{
		users: []models.Users{
			{Username: "sid"},
			{Username: "soham"},
		},
	}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	handler.GetAllUsers(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}
func TestGetAllUsers_Empty(t *testing.T) {
	service := &FakeService{
		users: []models.Users{},
	}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	handler.GetAllUsers(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}
func TestGetAllUsers_ServiceError(t *testing.T) {
	service := &FakeService{
		err: fmt.Errorf("failed to fetch users"),
	}
	handler := NewUserHandler(service)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 got %d", w.Code)
	}
}
func TestGetAllUsers_InvalidMethod(t *testing.T){
	service := &FakeService{}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("POST", "/users", nil)
	w := httptest.NewRecorder()
	handler.GetAllUsers(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 got %d", w.Code)
	}
}


func TestGetUserById_Success(t *testing.T) {
	service := &FakeService{
		users: []models.Users{
			{
				Username: "sidd",
				Userid: 1,
				Email:    "test@gmail.com",
			},
		},
	}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("GET", "/users/1", nil)
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.GetUserById(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}

func TestGetUserById_EmptySlice(t *testing.T) {
	service := &FakeService{
		users: []models.Users{}, 
	}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("GET", "/users/1", nil)
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.GetUserById(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}

func TestGetUserById_ServiceError(t *testing.T) {
	service := &FakeService{
		err: fmt.Errorf("failed"),
	}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("GET", "/users/1", nil)
	req.SetPathValue("userid", "1")
	w := httptest.NewRecorder()
	handler.GetUserById(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 got %d", w.Code)
	}
}

func TestGetUserById_InvalidMethod(t *testing.T) {
	service := &FakeService{}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("POST", "/users/1", nil)
	w := httptest.NewRecorder()
	handler.GetUserById(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 got %d", w.Code)
	}
}

func TestGetUserById_MissingID(t *testing.T) {
	service := &FakeService{}
	handler := NewUserHandler(service)
	req := httptest.NewRequest("GET", "/users/", nil)
	w := httptest.NewRecorder()
	handler.GetUserById(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}