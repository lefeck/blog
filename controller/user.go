package controller

import (
	"blog/common"
	"blog/forms"
	"blog/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct {
	useService service.UserService
}

func (u *UserController) List(c *gin.Context) {
	userListForm := forms.UserListForm{}
	if err := c.ShouldBind(&userListForm); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	total, userList := u.useService.List(userListForm.PageSize, userListForm.PageNum)
	if total == 0 {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("no found get data"))
	}

	common.ResponseSuccess(c, gin.H{
		"total":    total,
		"userlist": userList,
	})
}

func (u *UserController) Create(c *gin.Context) {
	userForm := forms.UserForm{}
	if err := c.ShouldBind(&userForm); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	user, err := u.useService.Create(userForm.GetUser())
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, user)
}

func (u *UserController) GetUserByID(c *gin.Context) {
	user, err := u.useService.GetUserByID(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, user)
}

func (u *UserController) Update(c *gin.Context) {
	updateUserForm := forms.UpdateUserForm{}
	if err := c.ShouldBind(&updateUserForm); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	id := c.Param("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	user, err := u.useService.Update(updateUserForm.GetUser(uint(uid)))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, user)
}

func (u *UserController) Delete(c *gin.Context) {
	err := u.useService.Delete(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, nil)
}

func NewUserController(useService service.UserService) *UserController {
	return &UserController{useService: useService}
}

func (u *UserController) Name() string {
	return "users"
}

func (u *UserController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/users", u.List)
	api.GET("/users/:id", u.GetUserByID)
	api.POST("/users", u.Create)
	api.DELETE("/users/:id", u.Delete)
	api.PUT("/users/:id", u.Update)
}
