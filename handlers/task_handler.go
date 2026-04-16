package handlers

import (
	"encoding/json"
	"go-sqlite/models"
	"go-sqlite/services"

	"log"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
)

type TaskHandler struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTaskByUserId(writer http.ResponseWriter, request *http.Request) {

	tracer := otel.Tracer("task-handler")
	ctx, span := tracer.Start(request.Context(), "gettaskuserbyid")
	defer span.End()

	// start := time.Now();

	useridstr := request.PathValue("userid")
	status := request.URL.Query().Get("status")
	sortby := request.URL.Query().Get("sortby")
	order := request.URL.Query().Get("order")
	cursor := request.URL.Query().Get("cursor")
	limitstr := request.URL.Query().Get("limit")
	pagenostr := request.URL.Query().Get("pageno")

	if useridstr == "" {
		http.Error(writer, "missing userid", http.StatusBadRequest)
		return
	}
	// Duration := time.Since(start).Seconds()

	tasks, err := h.service.GetTaskByUserId(ctx, useridstr, status, sortby, order, cursor, limitstr, pagenostr)
	if err != nil {
		log.Println("error in service function call", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		//metrics.HttpErrorsTotal.WithLabelValues("Get","/users/{userid}/tasks").Inc()
		return
	}

	//metrics.HttpRequestsTotal.WithLabelValues("Get","/users/{userid}/tasks","200").Inc()
	//metrics.HttpRequestDuration.WithLabelValues("Get","/users/{userid}/tasks").Observe(Duration)

	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message": "the tasks of the users are as follows",
		"tasks":   tasks,
	})
}

func (h *TaskHandler) InsertTask(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid Method type", 405)
		log.Println("Invalid Method type")
		return
	}

	tracer := otel.Tracer("task-handler")
	ctx, span := tracer.Start(request.Context(), "InsertTask")
	defer span.End()

	var newtask models.Task
	userIDStr := request.PathValue("userid")

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		log.Println("User id must be positive")
		http.Error(writer, "invalid userId", http.StatusBadRequest)
		return
	}
	if userID <= 0 {
		log.Println("User id must be positive")
		http.Error(writer, "userid must be positive", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(request.Body).Decode(&newtask)
	if err != nil {
		http.Error(writer, "Invalid body or empty body", http.StatusBadRequest)
		log.Println("error in fetching the data")
		return
	}
	if newtask.Name == "" {
		http.Error(writer, "name is required", http.StatusBadRequest)
		return
	}
	if newtask.Status == "" {
		http.Error(writer, "name is required", http.StatusBadRequest)
		return
	}
	newtask.UserId = userID
	err = h.service.InsertTask(ctx, newtask)
	if err != nil {
		log.Println("error in service function calling ")
		http.Error(writer, "Invalid body or empty body", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message":  "the task inserted succesfully into database ",
		"taskname": newtask.Name,
		"userid":   userID,
	})
}

func (h *TaskHandler) DeleteTask(writer http.ResponseWriter, request *http.Request) {
	idstr := request.PathValue("taskid")
	useridstr := request.PathValue("userid")
	// ctx := request.Context()

	tracer := otel.Tracer("task-handler")
	ctx, span := tracer.Start(request.Context(), "deleteTask")
	defer span.End()

	if idstr == "" || useridstr == "" {
		http.Error(writer, "missing id", http.StatusBadRequest)
		return
	}
	_, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(writer, "invalid task id", http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(useridstr)
	if err != nil {
		http.Error(writer, "invalid user id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTask(ctx, idstr, useridstr)
	if err != nil {
		log.Println("error in passing the data to the services", err)
		http.Error(writer, "invalid parameters", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message":        "task deleted succesfully",
		"deleted userid": useridstr,
		"deleted task":   idstr,
	})
}

func (h *TaskHandler) UpdateTask(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPatch {
		http.Error(writer, "Invalid method", 405)
		return
	}

	tracer := otel.Tracer("Task-Handler")
	ctx, span := tracer.Start(request.Context(), "updateTask")
	defer span.End()

	userid := request.PathValue("userid")
	taskid := request.PathValue("taskid")

	var reqbody struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	err := json.NewDecoder(request.Body).Decode(&reqbody)
	if err != nil {
		http.Error(writer, "Invalid body", 400)
		return
	}
	if userid == "" || taskid == "" {
		http.Error(writer, "missing id", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateTask(ctx, userid, taskid, reqbody.Name, reqbody.Status)
	if err != nil {
		http.Error(writer, err.Error(), 400)
		return
	}

	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message": "task updated successfully",
	})
}
