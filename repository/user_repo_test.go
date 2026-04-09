package repository

import (
	"fmt"
	"go-sqlite/models"
	"go-sqlite/testutils"
	"testing"
)

func TestInsertUser_Success(t *testing.T){
	db:=testutils.SetupTestDb()
	defer db.Close()
	repo:=NewUserRepository(db)
	user:=models.Users{
		Username: "Sid",
		Email: "sid@gmail.com",
	}
	err:=repo.InsertUser(user)
	if err!=nil{
		t.Errorf("Unexpected Error got %v",err)
	}
}
func TestInsertUser_DuplicateEmail(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewUserRepository(db)
	err := repo.InsertUser(models.Users{
		Username: "sid",
		Email:    "sid@.gmail.com",
	})
	err = repo.InsertUser(models.Users{
		Username: "sidkk",
		Email:    "sid@.gmail.com",
	})

	if err == nil {
		t.Errorf("expected error for duplicate email")
	}
}
func TestInsertUser_MultipleUsers(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()

	repo := NewUserRepository(db)

	for i := 0; i < 5; i++ {
		err := repo.InsertUser(models.Users{
			Username: fmt.Sprintf("userr%d", i),
			Email:    fmt.Sprintf("userrr%d&gmail.com", i),
		})
		if err!=nil{
			t.Errorf("Unexpected error %v",err)
		}
	}
}
func TestInsertUser_DBClosed(t *testing.T) {
	db := testutils.SetupTestDb()

	repo := NewUserRepository(db)
	db.Close() 

	err := repo.InsertUser(models.Users{
		Username: "sid",
		Email:    "sid@gmail.com",
	})

	if err==nil{
		t.Errorf("Expected error when Db is closed")
	}
}


func TestGetAllUsers_Success(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewUserRepository(db)
	users, err := repo.GetAllUsers()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(users) < 2 {
		t.Errorf("expected at least 2 users, got %d", len(users))
	}
}
func TestGetAllUsers_Empty(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	_, err := db.Exec("DELETE FROM tasks1")
	if err != nil {
		t.Fatalf("failed to delete tasks: %v", err)
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to delete users: %v", err)
	}
	repo := NewUserRepository(db)
	users, err := repo.GetAllUsers()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected empty list")
	}
}
func TestGetAllUsers_TableNotExist(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()

	_, err := db.Exec("DELETE FROM tasks1")
	if err != nil {
		t.Fatalf("failed to delete tasks: %v", err)
	}

	_, err = db.Exec("DROP TABLE users")
	if err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}

	repo := NewUserRepository(db)

	_, err = repo.GetAllUsers()

	if err == nil {
		t.Errorf("expected error when table does not exist")
	}
}
func TestGetAllUsers_DBClosed(t *testing.T) {
	db := testutils.SetupTestDb()
	repo := NewUserRepository(db)
	db.Close() 
	_, err := repo.GetAllUsers()

	if err == nil {
		t.Errorf("expected error when DB is closed")
	}
}


func TestGetUserById_Success(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewUserRepository(db)
	user, err := repo.GetUserById(1)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if user.Username != "user1" {
		t.Errorf("expected user1 got %s", user.Username)
	}
}
func TestGetUserById_NotFound(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	repo := NewUserRepository(db)
	_, err := repo.GetUserById(999)
	if err == nil {
		t.Errorf("expected error when user not found")
	}
}
func TestGetUserById_DBClosed(t *testing.T) {
	db := testutils.SetupTestDb()

	repo := NewUserRepository(db)
	db.Close()

	_, err := repo.GetUserById(1)

	if err == nil {
		t.Errorf("expected error when DB is closed")
	}
}

func TestGetUserById_TableNotExist(t *testing.T) {
	db := testutils.SetupTestDb()
	defer db.Close()
	_, err := db.Exec("DELETE FROM tasks1")
	if err != nil {
		t.Fatalf("failed to delete tasks: %v", err)
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to delete users: %v", err)
	}
	repo := NewUserRepository(db)
	_, err = repo.GetUserById(1)
	if err == nil {
		t.Errorf("expected error when table does not exist")
	}
}