package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-sqlite/Redis"
	"go-sqlite/models"
	"go-sqlite/repository"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

type UserServices struct {
	repo repository.UserRepo
	rdb  *redis.Client
}

type UserService interface {
	InsertUser(newuser models.Users) error
	GetUserById(idstr string) (models.Users, error)
	GetAllUsers(context.Context) ([]models.Users, error)
}

func NewUserServices(repo repository.UserRepo, rdb *redis.Client) *UserServices {
	return &UserServices{repo: repo,
		rdb: rdb}
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

	key := fmt.Sprintf("user:%d", id)

	if id <= 0 {

		log.Println(" Enter a positive number for the UserId")
		return user, fmt.Errorf("invalid user id ")
	}

	if Redis.GetCache(userserv.rdb, key, &user) {
		log.Println("cache hit")
		return user, nil
	}

	log.Println("cache miss")
	user, err = userserv.repo.GetUserById(id)
	log.Println(user.Userid)
	if err != nil {
		log.Println("error in fetching the user in service function", err)
		return user, err
	}

	Redis.SetCache(userserv.rdb, key, user, time.Minute*5)
	return user, nil
}

func (userserv *UserServices) GetAllUsers(ctx context.Context) ([]models.Users, error) {

	tracer := otel.Tracer("user-services")
	ctx, span := tracer.Start(ctx, "GetAllUsersService")
	defer span.End()

	key := "users:All"

	var users []models.Users

	if Redis.GetCache(userserv.rdb, key, &users) {
		log.Println("cache hit")
		return users, nil
	}

	log.Println("cache missed then hit db")
	users, err := userserv.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	Redis.SetCache(userserv.rdb, key, users, time.Minute*5)
	return users, nil
}
