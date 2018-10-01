package models

import (
	"sync"

	"github.com/pborman/uuid"
	"golang.org/x/oauth2"
)

type User struct {
	Id          string
	OAuth2Token *oauth2.Token
}

var db = map[string]*User{}
var dbMutex sync.RWMutex

func GetUser(id string) *User {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return db[id]
}

func NewUser() *User {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	user := &User{
		Id: uuid.New(),
	}
	db[user.Id] = user
	return user
}
