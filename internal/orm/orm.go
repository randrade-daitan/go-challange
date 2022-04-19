package orm

import (
	"challange/internal/repository"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Orm struct {
	*gorm.DB
}

func NewOrm() repository.Repository {
	var datetimePrecision = 2

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbDSN(),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DefaultDatetimePrecision:  &datetimePrecision,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return Orm{db}
}

func dbDSN() string {
	auth := repository.DBUser() + ":" + repository.DBPass() + "@"
	url := repository.DBProtocol + "(" + repository.DBURL + ")/"
	name := repository.DBName
	return auth + url + name + "?charset=utf8&parseTime=True&loc=Local"
}

func (db Orm) GetAllTasks() (t []repository.Task, err error) {
	err = db.Find(&t).Error
	return
}

func (db Orm) GetTaskByID(id int64) (t repository.Task, err error) {
	err = db.Where("id = ?", id).First(&t).Error
	return
}

func (db Orm) GetTasksByCompletion(isCompleted bool) (t []repository.Task, err error) {
	err = db.Where("completed = ?", isCompleted).Find(&t).Error
	return
}

func (db Orm) AddTask(t repository.Task) (id int64, err error) {
	err = db.Create(&t).Error
	id = t.ID
	return
}

func (db Orm) EditTask(t repository.Task) (err error) {
	err = db.Save(&t).Error
	return
}
