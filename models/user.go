// user
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/shxsun/redis"
	"strconv"
	"strings"
)

var client redis.Client

func init() {
	client.Addr = "127.0.0.1:6379"
	client.Db = 0
}

type User struct {
	Id    int
	Name  string
	Email string
}

var (
	ErrorInvalid    = errors.New("Invalid argument")
	ErrorRepeated   = errors.New("Repeated")
	ErrorUnexpected = errors.New("Unexpected")
)

func AddUser(user User) (err error) {
	if user.Name == "" {
		beego.Warn("user name in models add is empty")
		return ErrorInvalid
	}
	if user.Email == "" {
		beego.Warn("email in add user is empty")
		return ErrorInvalid
	}
	userdata, err := json.Marshal(user)
	if err != nil {
		beego.Error("json marshal user error")
		return ErrorUnexpected
	}
	// TODO: finished user data add
	beego.Debug(string(userdata))

	keys, err := client.Keys("user:*:name")
	if err != nil {
		return
	}
	exists := false
	for _, k := range keys {
		name, _ := client.Get(k)
		if string(name) == user.Name {
			exists = true
			break
		}
	}
	if exists {
		beego.Warn("user already registed")
		return ErrorRepeated
	}
	id, err := client.Incr("user:count")
	if err != nil {
		return
	}
	beego.Debug("user id:", id)
	client.Set(fmt.Sprintf("user:%d:name", id), []byte(user.Name))
	client.Set(fmt.Sprintf("user:%d:email", id), []byte(user.Email))
	return
}

func ListUser() (users []User, err error) {
	keys, err := client.Keys("user:*:name")
	if err != nil {
		return
	}
	users = make([]User, 0, 100)
	for _, k := range keys {
		vv := strings.Split(k, ":")
		var user User
		id, er := strconv.Atoi(vv[1])
		if er != nil {
			beego.Warn("Invalid key:", k)
			continue
		}
		user.Id = id
		name, er := client.Get(k)
		if er != nil {
			continue
		}
		user.Name = string(name)
		email, er := client.Get(fmt.Sprintf("user:%d:email", id))
		if er != nil {
			continue
		}
		user.Email = string(email)
		users = append(users, user)
	}
	return users, nil
}
