package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

type UserList struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

func (u User) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/user/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{user_id}").Doc("get user object").To(u.findUser))
	ws.Route(ws.POST("/{user_id}").To(u.updateUser))
	ws.Route(ws.PUT("").To(u.createUser))
	ws.Route(ws.DELETE("/{user_id}").To(u.deleteUser))
	container.Add(ws)
}

func (u User) findUser(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	user_id := request.PathParameter("user_id")
	//phone := request.QueryParameter("phone")

	user := User{}
	if user_id == "" {
		/*
			if phone != "" {
				db.Where("phone = ?", phone).First(&user)
				if user.ID == 0 {
					errmsg := fmt.Sprintf("cannot find user with phone %s", phone)
					response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
					return
				} else {
					response.WriteHeaderAndEntity(http.StatusOK, user)
					return
				}
			} else {
		*/
		userList := UserList{}
		userList.Users = make([]User, 0)
		db.Find(&userList.Users)
		userList.Count = len(userList.Users)
		response.WriteHeaderAndEntity(http.StatusOK, userList)
		return
		//}
	}

	id, err := strconv.Atoi(user_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get user, user_id is not integer, err %", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	db.First(&user, id)
	if user.ID == 0 {
		errmsg := fmt.Sprintf("cannot find user with id %s", user.ID)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, user)
		return
	}
}

func (u User) createUser(request *restful.Request, response *restful.Response) {
	user := User{}
	err := request.ReadEntity(&user)
	if err == nil {
		db.Create(&user)
		response.WriteHeaderAndEntity(http.StatusCreated, user)
		return
	} else {
		errmsg := fmt.Sprintf("cannot create user, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (u User) updateUser(request *restful.Request, response *restful.Response) {
	user_id := request.PathParameter("user_id")
	user := User{}
	err := request.ReadEntity(&user)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update user, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(user_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update user, path user_id is %s, err %s", user_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != user.ID {
		errmsg := fmt.Sprintf("cannot update user, path user_id %d is not equal to id %d in body content", id, user.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realUser := User{}
	db.First(&realUser, user.ID)
	if realUser.ID == 0 {
		errmsg := fmt.Sprintf("cannot update user, user_id %d is not exist", user.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realUser).Update(user)
	response.WriteHeaderAndEntity(http.StatusCreated, &realUser)
	return
}

func (u User) deleteUser(request *restful.Request, response *restful.Response) {
	user_id := request.PathParameter("user_id")
	id, err := strconv.Atoi(user_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete user, user_id is not integer, err %", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	user := User{}
	db.First(&user, id)
	if user.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&user)

	realUser := User{}
	db.First(&realUser, id)

	if realUser.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete user,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
