package models

import "time"
 

type Task struct{
	Userid int `json:"userId"`
	ID int  `json:"id"`
	NAME string `json:"name"`
	STATUS  string `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct{
	Userid int `json:"userId"`
	UserName string `json:"username"`
	Email string `json:"email"`
}