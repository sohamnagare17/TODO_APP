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
		UserId: 1, // valid user exists
	})

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
