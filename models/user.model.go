package models

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/bryansamperura/ticket-booking/db"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	AccountID string `json:"account_id"`
}

type UserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	AccountID string `json:"account_id"`
}

func StoreAccount(email string, password string, role string, account_id int) (Response, error) {
	var res Response

	v := validator.New()

	accounts := User{
		Email:     email,
		Password:  password,
		Role:      role,
		AccountID: strconv.Itoa(account_id),
	}

	err := v.Struct(accounts)

	if err != nil {
		return res, err
	}

	con := db.CreateConnection()

	sqlStatement := "INSERT INTO users (email, password, role, account_id) VALUES (?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(email, password, role, account_id)

	if err != nil {
		return res, err
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return res, err
	}

	res.Status = http.StatusCreated
	res.Message = "Inserted"
	res.Data = map[string]int64{
		"id": lastInsertedId,
	}

	return res, nil

}

func FindUserByEmail(email string) (Response, error) {
	var user User
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM users WHERE email = ?"

	row := con.QueryRow(sqlStatement, email)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.AccountID)

	if err == sql.ErrNoRows {
		// Return a custom error response if the record is not found
		return Response{Status: http.StatusNotFound, Message: "Not Found"}, nil
	} else if err != nil {
		return Response{}, err
	}

	res.Status = http.StatusOK
	res.Message = "OK"
	res.Data = user

	return res, nil
}
