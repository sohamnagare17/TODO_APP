package handlers

import (
	"encoding/json"
	"go-sqlite/models"
	"go-sqlite/services"
	"log"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	service *services.TaskServices
}

func NewTaskHandler(service *services.TaskServices) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTaskByUserId(writer http.ResponseWriter, request *http.Request) {

	useridstr := request.PathValue("userid")
	status := request.URL.Query().Get("status")
	sortby := request.URL.Query().Get("sortby")
	order := request.URL.Query().Get("order")
	cursor := request.URL.Query().Get("cursor")
	limitstr := request.URL.Query().Get("limit")
	pagenostr := request.URL.Query().Get("pageno")

	tasks, err := h.service.GetTaskByUserId(useridstr, status, sortby, order, cursor, limitstr, pagenostr)
	if err != nil {
		log.Println("error in service function call", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

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
		http.Error(writer, "userid must be positive", 400)
		return
	}

	err = json.NewDecoder(request.Body).Decode(&newtask)

	if err != nil {
		http.Error(writer, "Invalid body or empty body", 400)
		log.Println("error in fetching the data")
		return
	}
	newtask.UserId = userID

	err = h.service.InsertTask(newtask)
	if err != nil {
		log.Println("error in service function calling ")
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

	err := h.service.DeleteTask(idstr, useridstr)
	if err != nil {
		log.Println("error in passing the data to the services",err)
		http.Error(writer,"invalid parameters",400)
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

	err = h.service.UpdateTask(userid, taskid, reqbody.Name, reqbody.Status)
	if err != nil {
		http.Error(writer, err.Error(), 400)
		return
	}

	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message": "task updated successfully",
	})
}
