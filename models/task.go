package models

// import "time"
 

type Task struct{
	USERID int `json:"user_id"`
	ID int  `json:"id"`
	NAME string `json:"name"`
	STATUS  string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

