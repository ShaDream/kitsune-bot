package repository

import (
	"errors"
	"fmt"
	"github.com/ShaDream/kitsune-bot/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepositoryPostgres struct {
	db *sqlx.DB
}

func NewUserRepositoryPostgres(db *sqlx.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (u UserRepositoryPostgres) CreateUser(userId string, username string) (*models.User, error) {
	user := new(models.User)
	query := fmt.Sprintf("INSERT"+
		" INTO %s (id, username, score, translated_pages, edited_pages, checked_pages, cleaned_pages, typed_chapters)"+
		" VALUES($1, $2, $3, $4, $5, $6, $7, $8)"+
		" RETURNING id, username, score, translated_pages, edited_pages, checked_pages, cleaned_pages, typed_chapters", userTable)
	err := u.db.QueryRow(query, userId, username, 0, 0, 0, 0, 0, 0).
		Scan(&user.Id, &user.Username, &user.Score, &user.TranslatedPages, &user.EditedPages, &user.CheckedPages, &user.CleanedPages, &user.TypedPages)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("Не удалось создать пользователя.")
	}
	return user, nil
}

func (u UserRepositoryPostgres) GetUser(userId string) (*models.User, error) {
	user := new(models.User)
	query := fmt.Sprintf("SELECT id, username, score, translated_pages, edited_pages, checked_pages, cleaned_pages, typed_chapters FROM %s WHERE id=$1", userTable)
	err := u.db.QueryRow(query, userId).
		Scan(&user.Id, &user.Username, &user.Score, &user.TranslatedPages, &user.EditedPages, &user.CheckedPages, &user.CleanedPages, &user.TypedPages)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("Не удалось получить пользователя.")
	}
	return user, nil
}

func (u UserRepositoryPostgres) HasUser(userId string) bool {
	var hasUser bool
	query := fmt.Sprintf("SELECT exists(SELECT * FROM %s WHERE id=$1)", userTable)
	err := u.db.QueryRow(query, userId).Scan(&hasUser)
	if err != nil {
		logrus.Error(err)
		return false
	}
	return hasUser
}

func (u UserRepositoryPostgres) GetTopUsers(characteristic string) ([]*models.User, error) {
	users := make([]*models.User, 0)
	query := fmt.Sprintf("SELECT id, username, score, translated_pages, edited_pages, checked_pages, cleaned_pages, typed_chapters FROM %s ORDER BY $1 DESC LIMIT 10", userTable)
	rows, err := u.db.Query(query, characteristic)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("Не удалось получить топ.")
	}
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.Id, &user.Username, &user.Score, &user.TranslatedPages, &user.EditedPages, &user.CheckedPages, &user.CleanedPages, &user.TypedPages)
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("Не удалось получить топ.")
		}
		users = append(users, user)
	}
	return users, nil
}