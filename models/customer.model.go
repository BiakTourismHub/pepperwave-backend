package models

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bryansamperura/ticket-booking/db"
	"github.com/go-playground/validator/v10"
)

type Customer struct {
	Id       int    `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type CustomerRequest struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func FindAllCustomer() (Response, error) {
	var obj Customer
	var arrObj []Customer
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM customers"

	rows, err := con.Query(sqlStatement)

	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.FullName, &obj.Email, &obj.Phone)

		if err != nil {
			return res, err
		}

		arrObj = append(arrObj, obj)
	}

	res.Status = http.StatusOK
	res.Message = "OK"
	res.Data = arrObj

	return res, nil
}

func FindCustomerById(id int) (Response, error) {
	var customer Customer
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM customers WHERE id= ?"

	row := con.QueryRow(sqlStatement, id)

	err := row.Scan(&customer.Id, &customer.FullName, &customer.Email, &customer.Phone)
	if err == sql.ErrNoRows {
		// Return a custom error response if the record is not found
		return Response{Status: http.StatusNotFound, Message: "Not Found"}, nil
	} else if err != nil {
		return Response{}, err
	}

	res.Status = http.StatusOK
	res.Message = "OK"
	res.Data = customer

	return res, nil
}

func StoreCustomer(fullname string, email string, phone string) (Response, error) {
	var res Response

	v := validator.New()

	customers := Customer{
		FullName: fullname,
		Email:    email,
		Phone:    phone,
	}

	err := v.Struct(customers)

	if err != nil {
		return res, err
	}

	con := db.CreateConnection()

	sqlStatement := "INSERT INTO customers(fullname, email, phone) VALUES (?, ? , ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(fullname, email, phone)

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

func UpdateCustomer(id int, fullname string, email string, phone string) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "UPDATE customers SET fullname=?, email=?, phone=? WHERE id= ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(fullname, email, phone, id)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Status = http.StatusCreated
	res.Message = "Updated"
	res.Data = map[string]int64{
		"rows_affected": rowsAffected,
	}

	return res, nil
}

func DeleteCustomer(id int) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "DELETE FROM customers WHERE id = ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		err = errors.New("City not found")
		return res, err
	}

	if err != nil {
		return res, err
	}

	res.Status = http.StatusNoContent
	res.Message = "Deleted"
	res.Data = map[string]int64{
		"rows_affected": rowsAffected,
	}

	return res, nil
}
