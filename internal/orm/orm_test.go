package orm

import (
	"challange/internal/repository"
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestOrmDSN(t *testing.T) {
	os.Setenv("DBUSER", "aaa")
	os.Setenv("DBPASS", "bbb")

	dsn := dbDSN()
	expectedDSN := "aaa:bbb@tcp(127.0.0.1:3306)/challange?charset=utf8&parseTime=True&loc=Local"

	if dsn != expectedDSN {
		t.Errorf("orm dsn is incorrect")
	}
}

func TestGetAllTasks(t *testing.T) {
	repo := repository.RepositoryForTesting(testableOrm, t)
	repo.TestRepositoryGetTasks("SELECT * FROM `task`", t)
}

func TestGetTaskByID(t *testing.T) {
	repo := repository.RepositoryForTesting(testableOrm, t)
	repo.TestRepositoryGetTaskById("SELECT * FROM `task` WHERE id = ? ORDER BY `task`.`id` LIMIT 1", t)
}

func TestAddTask(t *testing.T) {
	task := repository.Task{
		ID:        30,
		Name:      "ORM Test",
		Completed: true,
	}
	repo := repository.RepositoryForTesting(testableOrm, t)

	repo.Mock.ExpectBegin()
	repo.Mock.
		ExpectExec("INSERT INTO `task` (`name`,`completed`,`id`) VALUES (?,?,?)").
		WithArgs(task.Name, task.Completed, task.ID).
		WillReturnResult(sqlmock.NewResult(task.ID, 1))
	repo.Mock.ExpectCommit()

	repo.TestRepositoryAddTask(task, t)
}

func TestEditTask(t *testing.T) {
	task := repository.Task{
		ID:        2,
		Name:      "edited",
		Completed: true,
	}
	repo := repository.RepositoryForTesting(testableOrm, t)

	repo.Mock.ExpectBegin()
	repo.Mock.
		ExpectExec("UPDATE `task` SET `name`=?,`completed`=? WHERE `id` = ?").
		WithArgs(task.Name, task.Completed, task.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	repo.Mock.ExpectCommit()

	repo.TestRepositoryEditTask(task, t)
}

func TestGetTasksByCompletion(t *testing.T) {
	q := "SELECT * FROM `task` WHERE completed = ?"

	t.Run("fetch completed tasks", func(t *testing.T) {
		repo := repository.RepositoryForTesting(testableOrm, t)
		repo.TestRepositoryGetTasksByCompletion(q, t, true)
	})

	t.Run("fetch incompleted tasks", func(t *testing.T) {
		repo := repository.RepositoryForTesting(testableOrm, t)
		repo.TestRepositoryGetTasksByCompletion(q, t, false)
	})
}

func testableOrm(db *sql.DB, t *testing.T) repository.Repository {
	t.Helper()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})
	ormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("error creating new mock: %v", err)
	}

	return Orm{ormDB}
}
