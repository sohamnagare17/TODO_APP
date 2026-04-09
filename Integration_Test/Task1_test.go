package Integration_Test

import(
	"go-sqlite/models"
	"go-sqlite/repository"
	"go-sqlite/testutils"
	"go-sqlite/services"
	"go-sqlite/handlers"
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"database/sql"
	
)

func GetTaskHandler(db *sql.DB) *handlers.TaskHandler {
	//db := testutils.SetupTestDb()
	repo := repository.NewTaskRepository(db)
	service := services.NewTaskServices(repo)
	handler := handlers.NewTaskHandler(service)
	return handler

}

func GetUserHandler (db * sql.DB) *handlers.UserHandler{

	repo := repository.NewUserRepository(db)
	service := services.NewUserServices(repo)
	handler :=handlers.NewUserHandler(service)
	return handler 
}

//insertTask
func TestInsertTask(t *testing.T){
	db := testutils.SetupTestDb()
	handler := GetTaskHandler(db)
    body := `{"Name":"task 1","Status":"pending"}`

	req := httptest.NewRequest(http.MethodPost,"/users/1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	req.SetPathValue("userid","1")

	recorder := httptest.NewRecorder()

	handler.InsertTask(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 statusok but got %d",recorder.Code)
	}

	rows, _ := db.Query("SELECT name FROM tasks1 WHERE userid=1")

	var name string
	found := false

	for rows.Next() {
		rows.Scan(&name)
		if name == "Task 1" {
			found = true
		}
	}

	if !found {
		t.Errorf("task not inserted")
	}
}

// GetTaskByUserId
func TestGetTaskUserId(t *testing.T){
	db := testutils.SetupTestDb()
	handler := GetTaskHandler(db)


	request := httptest.NewRequest(http.MethodGet,"/users/1/tasks",nil)
	request.Header.Set("Content-Type", "application/json")

	request.SetPathValue("userid","1")

	recorder := httptest.NewRecorder()

	handler.GetTaskByUserId(recorder,request)

	if recorder.Code != http.StatusOK{
		t.Fatalf("expected 200 statusOk but got %d",recorder.Code)
	}
	
}

func TestGetTaskUserId_invalidUserid(t *testing.T){
	db := testutils.SetupTestDb()
	handler := GetTaskHandler(db)


	request := httptest.NewRequest(http.MethodGet,"/users/abc/tasks",nil)
	request.Header.Set("Content-Type", "application/json")

	request.SetPathValue("userid","abc")

	recorder := httptest.NewRecorder()

	handler.GetTaskByUserId(recorder,request)

	if recorder.Code != http.StatusBadRequest{
		t.Fatalf("expected 400 error for invalid userid but got %d",recorder.Code)
	}
	
}

func TestGetTaskUserId_bystatus(t *testing.T){
	db := testutils.SetupTestDb()
	handler := GetTaskHandler(db)


	request := httptest.NewRequest(http.MethodGet,"/users/1/tasks?status=pending",nil)
	request.Header.Set("Content-Type", "application/json")

	request.SetPathValue("userid","1")

	recorder := httptest.NewRecorder()

	handler.GetTaskByUserId(recorder,request)

	if recorder.Code != http.StatusOK{
		t.Fatalf("expected 200 httpstatus but  got %d",recorder.Code)
	}

  type Response struct {
	Tasks []models.Task `json:"tasks"`
   }
    var resp Response

	  err := json.Unmarshal(recorder.Body.Bytes(),&resp)
      if err!=nil{
		  t.Fatalf("invalid response %v",err)
	  }	

	  var tasks []models.Task
      tasks = resp.Tasks
	  
	  for _,task := range tasks{
		 if task.Status != "pending"{
			 t.Errorf("expected only pending task but got %s",task.Status)

		 }
	  }
}

func TestGetTaskUserId_limitandoffset(t *testing.T){
	db := testutils.SetupTestDb()
	handler := GetTaskHandler(db)


	request := httptest.NewRequest(http.MethodGet,"/users/1/tasks?limit=1&page=1",nil)
	request.Header.Set("Content-Type", "application/json")

	request.SetPathValue("userid","1")

	recorder := httptest.NewRecorder()

	handler.GetTaskByUserId(recorder,request)

	if recorder.Code != http.StatusOK{
		t.Fatalf("expected 200 httpstatus but  got %d",recorder.Code)
	}

  type Response struct {
	Tasks []models.Task `json:"tasks"`
   }
    var resp Response

	  err := json.Unmarshal(recorder.Body.Bytes(),&resp)
      if err!=nil{
		  t.Fatalf("invalid response %v",err)
	  }	

	  var tasks []models.Task
      tasks = resp.Tasks
	  
	  if len(tasks)!=1{
		 t.Errorf("expected 1 task got %d",len(tasks))
	  }
}

// deletetask handler 

func TestDeleteTask_success(t *testing.T){
	db := testutils.SetupTestDb();
	  handler := GetTaskHandler(db)
	  

	  request := httptest.NewRequest(http.MethodDelete,"/users/1/tasks/1",nil)
	  request.Header.Set("Content-Type","application/json")

	  request.SetPathValue("userid","1")
	  request.SetPathValue("taskid","1")
    
	  recorder := httptest.NewRecorder()

	  handler.DeleteTask(recorder, request)

	  if recorder.Code != http.StatusOK{
          t.Fatalf("expected 200 status code but got %d",recorder.Code)
	  }

	  rows,_ := db.Query(`SELECT id FROM tasks1 WHERE id=1 AND userid=1`)
	  
	 if rows.Next(){
		t.Errorf("task was not deleted yet ")
	 }
}

func TestDeleteTask_invaliduserid(t *testing.T){
	db := testutils.SetupTestDb();
	  handler := GetTaskHandler(db)
	  

	  request := httptest.NewRequest(http.MethodDelete,"/users/abs/tasks/1",nil)
	  request.Header.Set("Content-Type","application/json")

	  request.SetPathValue("userid","abs")
	  request.SetPathValue("taskid","1")
    
	  recorder := httptest.NewRecorder()

	  handler.DeleteTask(recorder, request)

	  if recorder.Code != http.StatusBadRequest{
          t.Fatalf("expected 400 status code but got %d",recorder.Code)
	  }

}

func TestDeleteTask_invalidtaskid(t *testing.T){
	db := testutils.SetupTestDb();
	  handler := GetTaskHandler(db)
	  

	  request := httptest.NewRequest(http.MethodDelete,"/users/1/tasks/abs",nil)
	  request.Header.Set("Content-Type","application/json")

	  request.SetPathValue("userid","1")
	  request.SetPathValue("taskid","abs")
    
	  recorder := httptest.NewRecorder()

	  handler.DeleteTask(recorder, request)

	  if recorder.Code != http.StatusBadRequest{
          t.Fatalf("expected 400 status code but got %d",recorder.Code)
	  }

}

func TestDeleteTask_Nodata(t *testing.T){
	db := testutils.SetupTestDb();
	  handler := GetTaskHandler(db)
	  

	  request := httptest.NewRequest(http.MethodDelete,"/users/999/tasks/1",nil)
	  request.Header.Set("Content-Type","application/json")

	  request.SetPathValue("userid","999")
	  request.SetPathValue("taskid","1")
    
	  recorder := httptest.NewRecorder()

	  handler.DeleteTask(recorder, request)

	  if recorder.Code != http.StatusBadRequest{
          t.Fatalf("expected 400 status code but got %d",recorder.Code)
	  }

}

//GetUserById
func TestGetUserById_Success(t *testing.T) {

	db := testutils.SetupTestDb()

	handler := GetUserHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.SetPathValue("userid", "1")

	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}

	var user models.Users
	err := json.Unmarshal(rec.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("invalid response: %v", err)
	}

	if user.Userid != 1 {
		t.Errorf("expected user id 1, got %d", user.Userid)
	}
}

func TestGetUserById_NotFound(t *testing.T) {

	db := testutils.SetupTestDb()

	handler := GetUserHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	req.SetPathValue("userid", "99")

	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)
    
	if rec.Code != http.StatusNotFound{
		t.Errorf("expected 404 user not found but we got %d",rec.Code)
	}
}

func TestGetUserById_InvalidId(t *testing.T) {

	db := testutils.SetupTestDb()

	handler := GetUserHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req.SetPathValue("userid", "abc")

	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)
    
	if rec.Code != http.StatusBadRequest{
		t.Errorf("expected 400 bad request but we got %d",rec.Code)
	}
}

func TestGetUserById_EmptyId(t *testing.T) {

	db := testutils.SetupTestDb()

	handler := GetUserHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	req.SetPathValue("userid", "")

	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)
    
	if rec.Code != http.StatusBadRequest{
		t.Errorf("expected 400 user not found but we got %d",rec.Code)
	}
}










