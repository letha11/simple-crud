package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simple-crud-go/internal/models"
	"github.com/simple-crud-go/internal/repository"
	"github.com/stretchr/testify/assert"
)

var preloadUserQuery = "SELECT (.+) FROM `users` WHERE `users`.`id`"

func TestPostGetById(t *testing.T) {
	var (
		id    = 1
		title = "My First Post"
		body  = "My Body"
	)

	_, db, mock := DB(t)

	repo := repository.NewPostRepository(db)

	post := sqlmock.NewRows([]string{
		"id", "title", "body", "user_id",
	}).AddRow(id, title, body, 1)

	query := "SELECT (.+) FROM `posts` WHERE `posts`.`id` = ?"
	// WithArgs 2 arguments because gorm need 2 arguments on their SQL query
	mock.ExpectQuery(query).WithArgs(id, 1).WillReturnRows(post)
	// Preload (association) query
	mock.ExpectQuery(preloadUserQuery).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{}))
	p, err := repo.GetById(1)

	assert.NoError(t, err)
	assert.Equal(t, title, p.Title)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPostGetAll(t *testing.T) {
	var (
		id         = 1
		firstTitle = "My First Post"
		body       = "My Body"
	)

	_, db, mock := DB(t)

	repo := repository.NewPostRepository(db)

	post := sqlmock.NewRows([]string{
		"id", "title", "body", "user_id",
	}).AddRow(id, firstTitle, body, 1).AddRow(2, "Second Post!", "Second post body", 2)

	query := "SELECT (.+) FROM `posts`"
	// WithArgs 2 arguments because gorm need 2 arguments on their SQL query
	mock.ExpectQuery(query).WillReturnRows(post)
	// Preload (association) query
	mock.ExpectQuery(preloadUserQuery).WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{}))
	p, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, firstTitle, p[0].Title)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPostUpdate(t *testing.T) {
	_, db, mock := DB(t)

	newPost := models.Post{
		Title:  "First Post",
		Body:   "Brother",
		UserID: 1,
	}

	repo := repository.NewPostRepository(db)

	query := "INSERT INTO `posts`"
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(newPost.Title, newPost.Body, newPost.UserID, AnyTime{}, AnyTime{}, nil).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(&newPost)

	assert.NoError(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPostCreate(t *testing.T) {
	_, db, mock := DB(t)

	updatedPost := models.Post{
		ID:     1,
		Title:  "First Post",
		Body:   "Brother",
		UserID: 1,
	}

	repo := repository.NewPostRepository(db)

	query := "UPDATE `posts` SET"
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(updatedPost.Title, updatedPost.Body, updatedPost.UserID, AnyTime{}, AnyTime{}, nil, updatedPost.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(&updatedPost)

	assert.NoError(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPostDelete(t *testing.T) {
	_, db, mock := DB(t)

	repo := repository.NewPostRepository(db)

	postIdToBeDeleted := 1

	query := "UPDATE `posts` SET `deleted_at`=\\? WHERE `posts`.`id` = \\?" // query will be update because the model have `deleted_at` field when this field exist, gorm will automatically use update instead of DELETE (https://gorm.io/docs/delete.html#Soft-Delete)

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(AnyTime{}, postIdToBeDeleted).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(uint(postIdToBeDeleted))

	assert.NoError(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
