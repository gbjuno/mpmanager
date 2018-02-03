package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/gbjuno/mpmanager/backend/utils"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
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
	ws.Path(RESTAPIVERSION + "/user").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").Doc("get user object").To(u.findUser))
	ws.Route(ws.GET("/?name={name}&phone={phone}&pageNo={pageNo}&pageSize={pageSize}&order={order}").Doc("get user object").To(u.findUser))
	ws.Route(ws.GET("/{user_id}").Doc("get user object").To(u.findUser))
	ws.Route(ws.POST("").To(u.createUser))
	ws.Route(ws.POST("/").To(u.createUser))
	ws.Route(ws.PUT("/{user_id}").To(u.updateUser))
	ws.Route(ws.DELETE("/{user_id}").To(u.deleteUser))
	container.Add(ws)

	loginWs := new(restful.WebService)
	loginWs.Path(RESTAPIVERSION + "/login").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	loginWs.Route(loginWs.GET("").To(u.loginUser))
	loginWs.Route(loginWs.POST("").To(u.loginUser))
	loginWs.Route(loginWs.POST("/").To(u.loginUser))
	container.Add(loginWs)
	glog.Infof("register login service")
}

func (u User) loginUser(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [loginUser]", request.Request.RemoteAddr)
	if request.Request.Method == "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	user := User{}
	err := request.ReadEntity(&user)
	if err != nil {
		glog.Errorf("%s cannot parse user entity", prefix)
		r := Response{Status: "error", Error: "无法解析"}
		response.WriteHeaderAndEntity(http.StatusForbidden, &r)
		return
	}

	glog.Infof("%s,phone %s, password %s", prefix, user.Phone, user.Password)
	realUser := User{}
	db.Debug().Where("phone = ?", user.Phone).First(&realUser)
	if realUser.ID == 0 {
		glog.Errorf("%s user with phone %s doesn't exist", prefix, user.Phone)
		r := Response{Status: "error", Error: "手机号或密码错误"}
		response.WriteHeaderAndEntity(http.StatusUnauthorized, &r)
		return
	}
	if realUser.Admin == "F" {
		glog.Errorf("%s user with phone %s is not admin", prefix, user.Phone)
		r := Response{Status: "error", Error: "手机号或密码错误"}
		response.WriteHeaderAndEntity(http.StatusUnauthorized, &r)
		return
	}
	glog.Infof("%s find user with id %d", prefix, realUser.ID)
	/*
		encryptedPassword, err := utils.DesEncrypt([]byte(password), wxDESkey)
		if err != nil {
			glog.Errorf("%s unable to encrypt password %s", password)
			glog.Infof("%s password not match", prefix)
			return
		}*/

	// password match

	hashCode := md5.New()
	io.WriteString(hashCode, user.Password)
	md5pass := fmt.Sprintf("%x", hashCode.Sum(nil))

	if md5pass == realUser.Password {
		sessionid, err := newPasswordSession(fmt.Sprintf("%d", realUser.ID))
		if err != nil {
			glog.Errorf("%s %s", prefix, err)
			r := Response{Status: "error", Error: "请联系管理员进行处理"}
			response.WriteHeaderAndEntity(http.StatusUnauthorized, &r)
			return
		}

		cookie := http.Cookie{Name: "sessionid", Value: sessionid}
		glog.Infof("%s user %s id %d login successfully, sessionid %s", prefix, realUser.Name, realUser.ID, sessionid)
		http.SetCookie(response, &cookie)
		response.WriteHeader(http.StatusOK)
		return
	} else {
		//password not match
		glog.Errorf("%s user password not match", prefix)
		r := Response{Status: "error", Error: "手机号或密码错误"}
		response.WriteHeaderAndEntity(http.StatusUnauthorized, &r)
		return
	}
}

func (u User) findUser(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findUser]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	user_id := request.PathParameter("user_id")
	name := request.QueryParameter("name")
	phone := request.QueryParameter("phone")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchUser *gorm.DB = db.Debug()
	var noPageSearchUser *gorm.DB = db.Debug()
	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "asc"
	}

	if order == "" {
		order = "asc"
	}

	glog.Infof("%s find user with order %s", prefix, order)
	searchUser = searchUser.Order("id " + order)
	//get user list
	if user_id == "" {
		if name != "" {
			glog.Infof("%s set find user db with name %s", prefix, name)
			searchUser = searchUser.Where(fmt.Sprintf("name LIKE \"%%%s%%\"", name))
			noPageSearchUser = noPageSearchUser.Where(fmt.Sprintf("name LIKE \"%%%s%%\"", name))
		}

		if phone != "" {
			glog.Infof("%s set find user db with name %s", prefix, name)
			searchUser = searchUser.Where("phone = ?", phone)
			noPageSearchUser = noPageSearchUser.Where("phone = ?", phone)
		}

		isPageSizeOk := true
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find object with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find object with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find user db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchUser = searchUser.Offset(offset).Limit(limit)
		}

		userList := UserList{}
		userList.Users = make([]User, 0)
		searchUser.Find(&userList.Users)
		noPageSearchUser.Model(&User{}).Count(&userList.Count)
		for i, u := range userList.Users {
			company := Company{}
			db.First(&company, u.CompanyId)
			userList.Users[i].CompanyName = company.Name
		}

		response.WriteHeaderAndEntity(http.StatusOK, userList)
		glog.Infof("%s return user list", prefix)
		return
	}

	user := User{}
	id, err := strconv.Atoi(user_id)
	//fail to parse user id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get user, user_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	db.Debug().First(&user, id)
	//fail to find user
	if user.ID == 0 {
		errmsg := fmt.Sprintf("cannot find user with id %s", user_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	company := Company{}
	db.First(&company, user.ID)
	user.CompanyName = company.Name

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
		if user.Phone == "" {
			errmsg := fmt.Sprintf("please provide phone number")
			returnmsg := fmt.Sprintf("请提供手机号码")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}

		samePhoneUser := User{}
		db.Debug().Where("phone = ?", user.Phone).First(&samePhoneUser)
		if samePhoneUser.ID != 0 {
			errmsg := fmt.Sprintf("user with phone %s already exists", samePhoneUser.Phone)
			returnmsg := fmt.Sprintf("手机号码%s已存在，与现有用户%s冲突", samePhoneUser.Phone, samePhoneUser.Name)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}

		if user.CompanyId == 0 {
			errmsg := fmt.Sprintf("please provide company id")
			returnmsg := fmt.Sprintf("请提供/选择公司")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}

		company := Company{}
		db.First(&company, user.CompanyId)
		if company.ID == 0 {
			errmsg := fmt.Sprintf("company id %d not found", company.ID)
			returnmsg := fmt.Sprintf("公司%s已被删除", company.Name)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}

		var rawPass string
		if user.Password == "" {
			rawPass = user.Phone[5:]
		} else {
			rawPass = user.Password
		}

		glog.Infof("%s set user password %s", rawPass)

		hashCode := md5.New()
		io.WriteString(hashCode, rawPass)
		user.Password = fmt.Sprintf("%x", hashCode.Sum(nil))

		//user.EncryptPassword()
		db.Debug().Create(&user)
		if user.ID == 0 {
			//fail to create user on database
			errmsg := fmt.Sprintf("cannot create user on database")
			returnmsg := fmt.Sprintf("无法创建用户，请联系管理员")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
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
		returnmsg := fmt.Sprintf("创建用户失败，提供的信息错误")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
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
		returnmsg := fmt.Sprintf("无法更新用户信息，用户信息解析失败")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	id, err := strconv.Atoi(user_id)
	//fail to parse user id
	if err != nil {
		errmsg := fmt.Sprintf("cannot update user, path user_id is %s, err %s", user_id, err)
		returnmsg := fmt.Sprintf("无法更新用户信息，用户id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != user.ID {
		errmsg := fmt.Sprintf("cannot update user, path user_id %d is not equal to id %d in body content", id, user.ID)
		returnmsg := fmt.Sprintf("无法更新用户信息，用户id与URL中的用户id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realUser := User{}
	db.Debug().First(&realUser, user.ID)
	//cannot find user
	if realUser.ID == 0 {
		errmsg := fmt.Sprintf("cannot update user, user_id %d is not exist", user.ID)
		returnmsg := fmt.Sprintf("无法更新用户信息，用户已被删除")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//find user and update
	if user.Password != "" {
		rawPass := user.Password
		hashCode := md5.New()
		io.WriteString(hashCode, rawPass)
		glog.Infof("%s set user password %s", rawPass)
		user.Password = fmt.Sprintf("%x", hashCode.Sum(nil))
	}
	db.Debug().Model(&realUser).Update(user)
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
		returnmsg := fmt.Sprintf("删除用户失败，提供的用户id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	user := User{}
	db.Debug().First(&user, id)
	if user.ID == 0 {
		//user with id doesn't exist, return ok
		glog.Infof("%s user with id %s doesn't exist, return ok", prefix, user_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Debug().Delete(&user)

	realUser := User{}
	db.Debug().First(&realUser, id)

	if realUser.ID != 0 {
		//failed to delete user
		errmsg := fmt.Sprintf("cannot delete user,some of other object is referencing")
		returnmsg := fmt.Sprintf("无法删除用户，用户仍被引用")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	} else {
		//delete user successfully
		glog.Infof("%s delete user with id %s successfully", prefix, user_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
