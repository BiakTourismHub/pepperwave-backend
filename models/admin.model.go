package models

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bryansamperura/ticket-booking/db"
	"github.com/go-playground/validator/v10"
)

type Admin struct {
	Id       int    `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type AdminRequest struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func FindAllAdmin() (Response, error) {
	var obj Admin
	var arrObj []Admin
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM admin"

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

func FindAdminById(id int) (Response, error) {
	var admin Admin
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM admin WHERE id= ?"

	row := con.QueryRow(sqlStatement, id)

	err := row.Scan(&admin.Id, &admin.FullName, &admin.Email, &admin.Phone)
	if err == sql.ErrNoRows {
		// Return a custom error response if the record is not found
		return Response{Status: http.StatusNotFound, Message: "Not Found"}, nil
	} else if err != nil {
		return Response{}, err
	}

	res.Status = http.StatusOK
	res.Message = "OK"
	res.Data = admin

	return res, nil
}

func StoreAdmin(fullname string, email string, phone string) (Response, error) {
	var res Response

	v := validator.New()

	Admin := Admin{
		FullName: fullname,
		Email:    email,
		Phone:    phone,
	}

	err := v.Struct(Admin)

	if err != nil {
		return res, err
	}

	con := db.CreateConnection()

	sqlStatement := "INSERT INTO admin(fullname, email, phone) VALUES (?, ? , ?)"

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

func UpdateAdmin(id int, fullname string, email string, phone string) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "UPDATE admin SET fullname=?, email=?, phone=? WHERE id= ?"

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

func DeleteAdmin(id int) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "DELETE FROM admin WHERE id = ?"

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
