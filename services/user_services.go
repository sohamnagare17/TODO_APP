package services

import (
	"errors"
	"go-sqlite/models"
	"go-sqlite/repository"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"fmt"
)

type UserServices struct {
	repo repository.UserRepo
}

func NewUserServices(repo repository.UserRepo) *UserServices {
	return &UserServices{repo: repo}
}

func (userserv *UserServices) InsertUser(newuser models.Users) error {
	if newuser.Username == "" && newuser.Email == "" {

		log.Println("username and email required ")
		return errors.New("Username and Email required")
	}
	if newuser.Username == "" {

		log.Println("username required ")
		return errors.New("Username is required")
	}
	if newuser.Email == "" {
		//http.Error(writer, "Email  Required", 400)
		log.Println("Email required ")
		return errors.New("Email is required")
	}

	if strings.TrimSpace(newuser.Username) == "" {
		//http.Error(writer, "Username Required", 400)
		log.Println("Username Required")
		return errors.New("Username is required")
	}
	if strings.TrimSpace(newuser.Email) == "" {

		log.Println("Email is required")
		return errors.New("Email is Required")
	}
	_, err := mail.ParseAddress(newuser.Email)
	if err != nil {
		//http.Error(writer, "Invalid Email", http.StatusBadRequest)
		log.Println("Enter a valid Email")
		return errors.New("Enter a valid email")
	}
	if len(newuser.Username) < 2 {
		//http.Error(writer, "Name should greater than 2 characters", 400)
		return errors.New("Username is too short")
	}
	return userserv.repo.InsertUser(newuser)
}

func (userserv *UserServices) GetUserById(idstr string) (models.Users, error) {

	var user models.Users

	if idstr == "" {
		log.Println("Id required")
		return user, errors.New("Id is required")
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		log.Println("id must be number")
		//http.Error(writer,"Id must be Number",400)
		return user, err
	}
	if id <= 0 {

		log.Println(" Enter a positive number for the UserId")
		return user, fmt.Errorf("invalid user id ")
	}

	user, err = userserv.repo.GetUserById(id)
	log.Println(user.Userid)
	if err != nil {
		log.Println("error in fetching the user in service function", err)
		return user, err
	}
	return user, nil
}

func (userserv *UserServices) GetAllUsers() ([]models.Users, error) {
	return userserv.repo.GetAllUsers()
}
