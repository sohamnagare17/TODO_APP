package models

import "time"

 

type Task struct{
	UserId int `json:"user_id"`
	Id int  `json:"id"`
	Name string `json:"name"`
	Status  string `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

