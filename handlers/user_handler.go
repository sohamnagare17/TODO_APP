package handlers

import (
	"encoding/json"
	"go-sqlite/models"
	"go-sqlite/services"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) InsertUser(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid Method type", 405)
		log.Println("Invalid Method type")
		return
	}
	var newuser models.Users

	err := json.NewDecoder(request.Body).Decode(&newuser)

	if err != nil {
		http.Error(writer, "error in fetching the data", 400)
		log.Println("error in fetching the data")
		return
	}

	if newuser.Username == "" || newuser.Email == "" {

		http.Error(writer, "empty fields", 400)
		return
	}
	err = handler.service.InsertUser(newuser)
	log.Println("error", err)
	if err != nil {
		log.Println("error in service function calling ", err)
		http.Error(writer, "empty username or the email or may be both", 400)
		return
	}
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message": "the user inserted succesfully into database ",
	})
}

func (handler *UserHandler) GetUserById(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idstr := request.PathValue("userid")
	if idstr == "" {
		http.Error(writer, "missing userid", http.StatusBadRequest)
		return
	}
	user, err := handler.service.GetUserById(idstr)
	//log.Println("here ===",err)

	if err != nil {
		log.Println("error in fetching data in handler function", err)
		if err.Error() == "sql: no rows in result set" {
			http.Error(writer, "user not found", http.StatusNotFound)
		} else if err.Error() == "failed" {
			http.Error(writer, "error in service ", 500)
			return
		} else {

			http.Error(writer, "bad request", http.StatusBadRequest)
		}

	}

	if user.Userid == 0 {
		log.Println("bad request")
		http.Error(writer, "empty user", http.StatusBadRequest)
		return
	}

	log.Println(user.Userid)
	json.NewEncoder(writer).Encode(user)
}

func (handler *UserHandler) GetAllUsers(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Invalid method type", http.StatusMethodNotAllowed)
		return
	}

	tracer := otel.Tracer("user-handler")
	ctx, span := tracer.Start(request.Context(), "GET /user")
	defer span.End()

	users, err := handler.service.GetAllUsers(ctx)
	if err != nil {
		log.Println("error in handler function while calling service function")
		http.Error(writer, "Failed to fetch Users", 500)
	}
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"message":       "the users are ",
		"list of users": users,
	})
}
