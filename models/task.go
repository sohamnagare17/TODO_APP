 package models
 

type Task struct{
	ID int  `json:"id"`
	NAME string `json:"name"`
	STATUS  string `json:"status"`
}