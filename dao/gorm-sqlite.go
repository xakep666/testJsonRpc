package dao

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testJsonRpc/model"
	"testJsonRpc/helpers"
	"time"
	"errors"
)

type gormSqliteImpl struct {
	db *gorm.DB
}

func (g *gormSqliteImpl) Register(login string) error {
	uuid, err:=helpers.NewUUID()
	if err!=nil {
		return err
	}
	return transactionalOperation(g.db, func(tx *gorm.DB) (*gorm.DB, error){
		var user model.User
		if resp:=tx.Where("login = ?", login).Find(&user); !resp.RecordNotFound() {
			return resp, errors.New("user " + login + " is already registered")
		}
		if resp:=tx.Where("uuid = ?", uuid.String()).Find(&user); !resp.RecordNotFound() {
			return resp, errors.New("uuid " + uuid.String() + " is already present in base")
		}
		return tx.Create(&model.User{
			Login:            login,
			Uuid:             uuid.String(),
			RegistrationDate: time.Now(),
		}), nil
	})
}

func (g *gormSqliteImpl) GetByLogin(login string) (user model.User,err error) {
	err=transactionalOperation(g.db, func (tx *gorm.DB) (*gorm.DB, error) {
		resp:=tx.Where("login = ?", login).Find(&user)
		if resp.RecordNotFound() {
			return resp, errors.New("user " + login + " was not found")
		}
		return resp, nil
	})
	return
}

func (g *gormSqliteImpl) Edit(login string, newData model.User) error {
	return transactionalOperation(g.db, func(tx *gorm.DB) (*gorm.DB, error) {
		var userRec model.User
		if resp := tx.Where("login = ?", login).Find(&userRec); resp.RecordNotFound() {
			return resp, errors.New("user " + login + " is not registered")
		}
		var userRec1 model.User
		if resp := tx.Where("uuid = ?", newData.Uuid).Find(&userRec1); !resp.RecordNotFound() {
			return resp, errors.New("uuid " + newData.Uuid + " is already present in base")
		}
		var userRec2 model.User
		if resp := tx.Where("login = ?", newData.Login).First(&userRec2); !resp.RecordNotFound() {
			return resp, errors.New("user " + newData.Login + " is already registered")
		}
		return tx.Model(&userRec2).Update(newData), nil
	})
}

func (g *gormSqliteImpl) Close() {
	g.db.Close()
}

func SetupDb(path string) (err error) {
	var impl gormSqliteImpl
	impl.db, err = gorm.Open("sqlite3",path)
	if err!=nil {
		return
	}
	impl.db.AutoMigrate(&model.User{})
	dbImpl = &impl
	return
}

