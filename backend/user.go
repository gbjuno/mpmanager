package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chanxuehong/rand"
	"github.com/emicklei/go-restful"
	"github.com/gbjuno/mpmanager/backend/utils"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UserList struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

func (u *User) DecryptPassword() (err error) {
	prefix := fmt.Sprintf("[%s]", "DecryptPassword")
	glog.Infof("user password decrypt as %s", u.Password)
	decryptPass, err := utils.DesDecrypt([]byte(u.Password), wxDESkey)
	if err != nil {
		errmsg := fmt.Sprintf("cannot decrypt password, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		return errors.New(errmsg)
	}
	u.Password = string(decryptPass)
	glog.Infof("user password decrypt as %s", u.Password)
	return nil
}

func (u *User) EncryptPassword() (err error) {
	prefix := fmt.Sprintf("[%s]", "EncryptPassword")
	glog.Infof("user password decrypt as %s", u.Password)
	encryptPass, err := utils.DesEncrypt([]byte(u.Password), wxDESkey)
	if err != nil {
		errmsg := fmt.Sprintf("cannot encrypt password, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		return errors.New(errmsg)
	}
	u.Password = string(encryptPass)
	glog.Infof("user password encrypt as %s", u.Password)
	return nil
}

func (u User) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/user").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("?offset={offset}&limit={limit}").Doc("get user object").To(u.findUser))
	ws.Route(ws.GET("/{user_id}").Doc("get user object").To(u.findUser))
	ws.Route(ws.POST("").To(u.createUser))
	ws.Route(ws.PUT("/{user_id}").To(u.updateUser))
	ws.Route(ws.DELETE("/{user_id}").To(u.deleteUser))
	container.Add(ws)
}

func (u User) findUser(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findUser]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	user_id := request.PathParameter("user_id")
	//phone := request.QueryParameter("phone")
	offset := request.QueryParameter("offset")
	limit := request.QueryParameter("limit")

	user := User{}
	//get user list
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
		offsetOk := true
		limitOk := true
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			offsetOk = false
			glog.Errorf("%s offset %s is not integer, ignore", prefix, offset)
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			limitOk = false
			glog.Errorf("%s limit %s is not integer, ignore", prefix, limit)
		}
		userList := UserList{}
		userList.Users = make([]User, 0)
		if offsetOk && limitOk {
			db.Offset(offsetInt).Limit(limit).Find(&userList.Users)
		} else if offsetOk {
			db.Offset(offsetInt).Find(&userList.Users)
		} else if limitOk {
			db.Limit(limitInt).Find(&userList.Users)
		} else {
			db.Find(&userList.Users)
		}

		userList.Count = len(userList.Users)
		/*
			for _, userTemp := range userList.Users {
				glog.Infof("%s decrypt user password %v", prefix, []byte(user.Password))
				userTemp.DecryptPassword()
			}
		*/
		response.WriteHeaderAndEntity(http.StatusOK, userList)
		return
		//}
	}

	id, err := strconv.Atoi(user_id)
	//fail to parse user id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get user, user_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	db.First(&user, id)
	//fail to find user
	if user.ID == 0 {
		errmsg := fmt.Sprintf("cannot find user with id %s", user_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find user
	//user.DecryptPassword()
	glog.Infof("%s return user with id %d", prefix, user.ID)
	response.WriteHeaderAndEntity(http.StatusOK, user)
	return
}

func (u User) createUser(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createUser]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	user := User{}
	err := request.ReadEntity(&user)
	if err == nil {
		if user.Password == "" {
			rawBytes := rand.New()
			user.Password = string(rawBytes[:6])
		}
		//user.EncryptPassword()
		db.Create(&user)
		if user.ID == 0 {
			//fail to create user on database
			errmsg := fmt.Sprintf("cannot create user on database")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		} else {
			//create user on database
			glog.Infof("%s create user with id %d succesfully", prefix, user.ID)
			response.WriteHeaderAndEntity(http.StatusOK, user)
			return
		}
	} else {
		//fail to parse user entity
		errmsg := fmt.Sprintf("cannot create user, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (u User) updateUser(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateCompany]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	user_id := request.PathParameter("user_id")
	user := User{}
	err := request.ReadEntity(&user)

	//fail to parse user entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update user, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(user_id)
	//fail to parse user id
	if err != nil {
		errmsg := fmt.Sprintf("cannot update user, path user_id is %s, err %s", user_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != user.ID {
		errmsg := fmt.Sprintf("cannot update user, path user_id %d is not equal to id %d in body content", id, user.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realUser := User{}
	db.First(&realUser, user.ID)
	//cannot find user
	if realUser.ID == 0 {
		errmsg := fmt.Sprintf("cannot update user, user_id %d is not exist", user.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//find user and update
	if user.Password != "" {
		//if user password is update
		//user.EncryptPassword()
	}
	db.Model(&realUser).Update(user)
	glog.Infof("%s update user with id %d successfully and return", prefix, realUser.ID)
	response.WriteHeaderAndEntity(http.StatusCreated, realUser)
	return
}

func (u User) deleteUser(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteCompany]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	user_id := request.PathParameter("user_id")
	id, err := strconv.Atoi(user_id)
	//fail to parse user id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete user, user_id %s is not integer, err %s", user_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	user := User{}
	db.First(&user, id)
	if user.ID == 0 {
		//user with id doesn't exist, return ok
		glog.Infof("%s user with id %s doesn't exist, return ok", prefix, user_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&user)

	realUser := User{}
	db.First(&realUser, id)

	if realUser.ID != 0 {
		//failed to delete user
		errmsg := fmt.Sprintf("cannot delete user,some of other object is referencing")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		//delete user successfully
		glog.Infof("%s delete user with id %s successfully", prefix, user_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
