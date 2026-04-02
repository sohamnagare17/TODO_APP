package services

import(
  "go-sqlite/models"
  "log"
  "go-sqlite/repository"
  "strings"
  "net/mail"
  "strconv"
)

type UserServices struct{
	repo *repository.UserRepository
}

func NewUserServices(repo *repository.UserRepository) *UserServices{
	return &UserServices{repo:repo}
}

func (userserv *UserServices) InsertUser(newuser models.Users)error{
    	if newuser.Username == "" && newuser.Email == "" {
			
			log.Println("username and email required ")
			return  nil
		}
		if newuser.Username == "" {
			
			log.Println("username required ")
			return nil
		}
		if newuser.Email == "" {
			//http.Error(writer, "Email  Required", 400)
			log.Println("Email required ")
			return nil
		}

		if strings.TrimSpace(newuser.Username) == "" {
			//http.Error(writer, "Username Required", 400)
			log.Println("Username Required")
			return nil
		}
		if strings.TrimSpace(newuser.Email) == "" {
			
			log.Println("Email is required")
			return nil
		}
		_, err := mail.ParseAddress(newuser.Email)
		if err != nil {
			//http.Error(writer, "Invalid Email", http.StatusBadRequest)
			log.Println("Enter a valid Email")
			return err
		}
		if len(newuser.Username) < 2 {
			//http.Error(writer, "Name should greater than 2 characters", 400)
			return nil
		}
		return userserv.repo.InsertUser(newuser)
}

func (userserv *UserServices) GetUserById(idstr string) (models.Users,error){
     
       var user models.Users

		if idstr == "" {
			log.Println("Id required")
			return user,nil
		}

		id, err := strconv.Atoi(idstr)
		if err != nil {
			log.Println("id must be number")
			//http.Error(writer,"Id must be Number",400)
			return user,err
		}
		if id <=0 {

			log.Println(" Enter a positive number for the UserId")
			return user,nil
		}

		user,err = userserv.repo.GetUserById(id)
		if err!=nil{
			log.Println("error in fetching the user in service function" ,err)
			return user,err
		}
		return user,nil
}
func (userserv *UserServices) GetAllUsers()([]models.Users,error){
	return userserv.repo.GetAllUsers();
}

