package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simple-crud-go/internal/models"
	"github.com/simple-crud-go/internal/repository"
	"github.com/stretchr/testify/assert"
)

var preloadPostsQuery = "SELECT (.+) FROM `posts` WHERE `posts`.`user_id` = ?"

func TestUserGetById(t *testing.T) {
	_, db, mock := DB(t)

	repo := repository.NewUserRepository(db)

	user := sqlmock.NewRows([]string{
		"id", "name", "username", "password",
	}).AddRow(1, "Ibka", "ibkaanhar1", "")

	query := "SELECT (.+) FROM `users` WHERE `users`.`id` = ?"
	// preloadQuery :=
	// WithArgs 2 arguments because gorm need 2 arguments on their SQL query
	mock.ExpectQuery(query).WithArgs(1, 1).WillReturnRows(user)
	// Preload (association) query
	mock.ExpectQuery(preloadPostsQuery).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{}))
	u, err := repo.GetById(uint(1))

	assert.NoError(t, err)
	assert.Equal(t, "Ibka", u.Name)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUserGetByUsername(t *testing.T) {
	var (
		id       = 1
		username = "ibkaanhar1"
	)

	_, db, mock := DB(t)

	repo := repository.NewUserRepository(db)

	user := sqlmock.NewRows([]string{
		"id", "name", "username", "password",
	}).AddRow(id, "Ibka", username, "")

	query := "SELECT (.+) FROM `users` WHERE username = ?"
	mock.ExpectQuery(query).WithArgs(username, 1).WillReturnRows(user)
	mock.ExpectQuery(preloadPostsQuery).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{}))

	u, err := repo.GetByUsername(username)

	assert.NoError(t, err)
	assert.Equal(t, username, u.Username)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUserGetAll(t *testing.T) {
	firstUserUsername := "ibkaanhar1"

	_, db, mock := DB(t)

	repo := repository.NewUserRepository(db)

	users := sqlmock.NewRows([]string{
		"id", "name", "username", "password",
	}).AddRow(1, "Ibka", "ibkaanhar1", "").AddRow(2, "Ibka 2", "ibkaanhar2", "")

	query := "SELECT (.+) FROM `users`"
	mock.ExpectQuery(query).WillReturnRows(users)

	u, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, u[0].Username, firstUserUsername)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUserCreate(t *testing.T) {
	_, db, mock := DB(t)

	repo := repository.NewUserRepository(db)

	newUser := models.User{
		Name:     "Ibka",
		Username: "ibkaanhar",
		Password: "123",
	}

	query := "INSERT INTO `users`"

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(newUser.Name, newUser.Username, newUser.Password, AnyTime{}, AnyTime{}, nil).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(newUser)

	assert.NoError(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUserUpdate(t *testing.T) {
	_, db, mock := DB(t)

	repo := repository.NewUserRepository(db)

	updatedUser := models.User{
		ID:       1,
		Name:     "Ibka",
		Username: "ibkaanhar",
		Password: "123",
	}

	query := "UPDATE `users`"

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(updatedUser.Name, updatedUser.Username, updatedUser.Password, AnyTime{}, AnyTime{}, nil, updatedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(updatedUser)

	assert.NoError(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
