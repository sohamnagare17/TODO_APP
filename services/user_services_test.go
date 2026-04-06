package services

import (
	"errors"
	"go-sqlite/models"
	"testing"
)

type FakeUserRepo struct {
	users []models.Users
	err   error
}

func (f *FakeUserRepo) InsertUser(user models.Users) error {
	if f.err != nil {
		return f.err
	}
	f.users = append(f.users, user)
	return nil
}

func (f *FakeUserRepo) GetAllUsers() ([]models.Users, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.users, nil
}

func (f *FakeUserRepo) GetUserById(id int) (models.Users, error) {
	if f.err != nil {
		return models.Users{}, f.err
	}
	if id < 0 || id >= len(f.users) {
		return models.Users{}, errors.New("User Not Found")
	}
	return f.users[id], nil
}

func setup() (*UserServices, *FakeUserRepo) {
	fakerepo := &FakeUserRepo{}
	service := NewUserServices(fakerepo)
	return service, fakerepo
}

func TestInsertUser_Success(t *testing.T) {
	service, fakerepo := setup()
	user := models.Users{
		Username: "sid",
		Email:    "sid@gmail.com",
	}
	err := service.InsertUser(user)
	if err != nil {
		t.Errorf("Expected No error got %v", err)
	}
	if len(fakerepo.users) != 1 {
		t.Errorf("Expected 1 user got %d", len(fakerepo.users))
	}
}

func TestInsertUser_EmptyFields(t *testing.T) {
	service, _ := setup()
	user := models.Users{}
	err := service.InsertUser(user)
	if err == nil {
		t.Errorf("Expected error for empty fields")
	}
}

func TestInsertUser_InvalidEmail(t *testing.T) {
	service, _ := setup()

	user := models.Users{
		Username: "sid",
		Email:    "invalidemail",
	}

	err := service.InsertUser(user)
	if err == nil {
		t.Errorf("Expected error invalid mail got %v", err)
	}
}

func TestInsertUser_DBError(t *testing.T) {

	fakerepo := &FakeUserRepo{
		err: errors.New("DB ERROR"),
	}
	service := NewUserServices(fakerepo)

	user := models.Users{
		Username: "siddharth",
		Email:    "sid@gmail.com",
	}

	err := service.InsertUser(user)
	if err == nil {
		t.Errorf("Expected DB error")
	}
}

func TestGetAllUsers_Success(t *testing.T) {
	fakerepo := &FakeUserRepo{
		users: []models.Users{
			{Username: "SID", Email: "sid@gmail.com"},
			{Username: "SOHAM", Email: "soham@gmail.com"},
		},
	}
	service := NewUserServices(fakerepo)
	users, err := service.GetAllUsers()
	if err != nil {
		t.Errorf("Unexpected Error got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users got %d", len(users))
	}
}

func TestGetAllUsers_DBError(t *testing.T) {
	fakerepo := &FakeUserRepo{
		err: errors.New("DB error"),
	}
	service := NewUserServices(fakerepo)

	users, err := service.GetAllUsers()
	if err == nil {
		t.Errorf("Expected error got nil")
	}
	if users != nil {
		t.Errorf("expected nil users got %v", users)
	}
}

func TestGetAllUsers_EmptyList(t *testing.T) {
	fakerepo := &FakeUserRepo{
		users: []models.Users{},
	}
	service := NewUserServices(fakerepo)
	users, err := service.GetAllUsers()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Expected 0 users got %d", len(users))
	}
}

func TestGetAllUsers_NilUsers(t *testing.T) {
	fakerepo := &FakeUserRepo{
		users: nil,
	}
	service := NewUserServices(fakerepo)
	users, err := service.GetAllUsers()

	if err != nil {
		t.Errorf("Unexpected Error got %v", err)
	}
	if users != nil {
		t.Errorf("expected nil errors got %v", users)
	}
}

func TestGetUserById_Success(t *testing.T) {
	fakerepo := &FakeUserRepo{
		users: []models.Users{
			{Username: "sid", Email: "sid@gmail.com"},
			{Username: "soham", Email: "soham@gmail.com"},
		},
	}
	service := NewUserServices(fakerepo)
	user, err := service.GetUserById("1")
	if err != nil {
		t.Errorf("Unexpected Error got %v", err)
	}
	if user.Username != "soham" {
		t.Errorf("Expected soham got %v", user.Username)
	}
}

func TestGetUserById_NotPresent(t *testing.T) {
	fakerepo := &FakeUserRepo{
		users: []models.Users{
			{Username: "sid", Email: "sid@gmail.com"},
		},
	}
	service := NewUserServices(fakerepo)
	_, err := service.GetUserById("1")
	if err == nil {
		t.Errorf("expected user not found error")
	}
}

func TestGetUserById_DBError(t *testing.T) {
	fakerepo := &FakeUserRepo{
		err: errors.New("DB error"),
	}
	service := NewUserServices(fakerepo)
	_, err := service.GetUserById("1")
	if err == nil {
		t.Errorf("Expected db error")
	}
}

func TestGetUserById_InvalidId(t *testing.T) {
	service, _ := setup()
	_, err := service.GetUserById("abc")
	if err == nil {
		t.Errorf("Expected Invalid Id error")
	}
}

func TestGetUserById_EmptyId(t *testing.T) {
	service, _ := setup()
	_, err := service.GetUserById("")
	if err == nil {
		t.Errorf("Expected id required error")
	}
}
