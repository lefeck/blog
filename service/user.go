package service

import (
	"blog/model"
	"blog/repository"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	List(pageSize int, pageNum int) (int, []interface{})
	Create(user *model.User) (*model.User, error)
	GetUserByID(string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(string) error
	GetAll(userlist []model.User) ([]model.User, error)
	//Validate(*model.User) error
	//Auth(*AuthUser) (*model.User, error)
	//Login(param request.Login) (*model.User, error)
	//Export(data *[]model.User, headerName []string, filename string, c *gin.Context) error
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (u *userService) List(pageSize int, pageNum int) (int, []interface{}) {
	return u.userRepository.List(pageSize, pageNum)
}

func (u *userService) Create(user *model.User) (*model.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PassWord = string(password)
	return u.userRepository.Create(user)
}

func (u *userService) Update(user *model.User) (*model.User, error) {
	if len(user.PassWord) > 0 {
		password, err := bcrypt.GenerateFromPassword([]byte(user.PassWord), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PassWord = string(password)
	}
	return u.userRepository.Update(user)
}

func (u *userService) GetUserByID(id string) (*model.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return u.userRepository.GetUserByID(uid)
}

func (u *userService) Delete(id string) error {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return u.userRepository.Delete(uid)
}

func (u *userService) GetAll(userlist []model.User) ([]model.User, error) {
	return u.userRepository.GetAll(userlist)
}
