package repository

import (
	"database/sql"
	"go-sqlite/models"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) InsertUser(newuser models.Users) error {
	query := `INSERT INTO users(username,email) VALUES(?,?) `

	_, err := repo.db.Exec(query, newuser.Username, newuser.Email)

	if err != nil {
		log.Println("error while inserting the user ", err)
		return err
	}

	return nil
}

func (repo *UserRepository) GetUserById(id int) (models.Users, error) {
	var user models.Users
	query := `SELECT * FROM users WHERE user_id=?`

	err1 := repo.db.QueryRow(query, id).Scan(&user.Userid, &user.Username, &user.Email)

	if err1 != nil {
		log.Println("error in executing database query")
		return user, nil
	}

	return user, nil
}

func (repo *UserRepository) GetAllUsers() ([]models.Users, error) {
	var userlist []models.Users

	query := `SELECT * FROM users`

	rows, err := repo.db.Query(query)

	if err != nil {
		log.Println("error in fetching the data from the database", err)
		return userlist, err
	}

	for rows.Next() {
		var user models.Users

		err = rows.Scan(&user.Userid, &user.Email, &user.Username)
		if err != nil {
			log.Println("error in scanning the data from the rows", err)
		}
		userlist = append(userlist, user)
	}
	return userlist, nil
}
