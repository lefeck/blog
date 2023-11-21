package repository

import (
	"blog/global"
	"blog/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (u *userRepository) List(pageSize int, pageNum int) (int, []interface{}) {
	var users []model.User
	userList := make([]interface{}, 0, len(users))

	if err := global.DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error; err != nil {
		return 0, nil
	}
	total := len(users)
	for _, user := range users {
		userItemMap := map[string]interface{}{
			"id":       user.ID,
			"name":     user.UserName,
			"password": user.PassWord,
			"email   ": user.Email,
			"avatar  ": user.Avatar,
		}
		userList = append(userList, userItemMap)
	}
	return total, userList

}

func (u *userRepository) Create(user *model.User) (*model.User, error) {
	//var user *model.User

	userdata := []string{"username", "password", "email"}

	if err := global.DB.Select(userdata).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) Delete(id int) error {
	var user *model.User
	if err := global.DB.Delete(&user, id).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepository) Update(user *model.User) (*model.User, error) {

	//var userdata = make(map[string]interface{})
	if err := global.DB.Model(&model.User{}).Where(&user, "id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) GetAll(userlist []model.User) ([]model.User, error) {
	if err := global.DB.Find(&userlist).Error; err != nil {
		return nil, err
	}
	return userlist, nil
}

func (u *userRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := global.DB.Where(&user, "id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetUserByName(name string) (*model.User, error) {
	var user model.User
	if err := global.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) Migrate() error {
	return global.DB.AutoMigrate()
}
