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

// Creates a new repository using the ORM implementation.
func NewRepository() repository.Repository {
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
	url := repository.DBProtocol + "(" + repository.DBURL() + ")/"
	name := repository.DBName()
	return auth + url + name + "?charset=utf8&parseTime=True&loc=Local"
}

func (repo Orm) GetAllTasks() (tasks []repository.Task, err error) {
	err = repo.Find(&tasks).Error
	return
}

func (repo Orm) GetTaskByID(id int64) (task repository.Task, err error) {
	err = repo.Where("id = ?", id).First(&task).Error
	return
}

func (repo Orm) GetTasksByCompletion(isCompleted bool) (tasks []repository.Task, err error) {
	err = repo.Where("completed = ?", isCompleted).Find(&tasks).Error
	return
}

func (repo Orm) AddTask(task repository.Task) (id int64, err error) {
	err = repo.Create(&task).Error
	id = task.ID
	return
}

func (repo Orm) EditTask(task repository.Task) (err error) {
	err = repo.Save(&task).Error
	return
}
