package repository

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Orm struct {
	*gorm.DB
}

func newOrmRepository() Repository {
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
	auth := DBUser() + ":" + DBPass() + "@"
	url := DBProtocol + "(" + DBURL() + ")/"
	name := DBName()
	return auth + url + name + "?charset=utf8&parseTime=True&loc=Local"
}

func (repo Orm) GetAllTasks() (tasks []Task, err error) {
	err = repo.Find(&tasks).Error
	return
}

func (repo Orm) GetTaskByID(id int64) (task Task, err error) {
	err = repo.Where("id = ?", id).First(&task).Error
	return
}

func (repo Orm) GetTasksByCompletion(isCompleted bool) (tasks []Task, err error) {
	err = repo.Where("completed = ?", isCompleted).Find(&tasks).Error
	return
}

func (repo Orm) AddTask(task Task) (id int64, err error) {
	err = repo.Create(&task).Error
	id = task.ID
	return
}

func (repo Orm) EditTask(task Task) (err error) {
	err = repo.Save(&task).Error
	return
}
