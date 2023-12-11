package models

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bryansamperura/ticket-booking/db"
	"github.com/go-playground/validator/v10"
)

type Destination struct {
	Id              int    `json:"id"`
	DestinationName string `json:"destination_name"`
	Image           string `json:"image" form:"image"`
	City            string `json:"city_id"`
	CityName        int    `json:"city_name"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
}

type DestinationRequest struct {
	DestinationName string `json:"destination_name" form:"destination_name"`
	Image           string `json:"-" form:"image"`
	City            string `json:"city_id" form:"city_id"`
	Description     string `json:"description" from:"description"`
	Price           int    `json:"price" from:"price"`
}

type DestinationRequestBody struct {
	DestinationName string `json:"destination_name"`
	City            string `json:"city_id"`
	Description     string `json:"description"`
}

func FindAllDestination() (Response, error) {
	var obj Destination
	var arrObj []Destination
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM destination"

	rows, err := con.Query(sqlStatement)

	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.DestinationName, &obj.Image, &obj.City, &obj.Description, &obj.Price)
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

// FindById retrieves a city by its ID
func FindDestinationById(id int) (Response, error) {
	var destination Destination
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM destination WHERE id= ?"

	row := con.QueryRow(sqlStatement, id)

	err := row.Scan(&destination.Id, &destination.DestinationName, &destination.Image, &destination.City, &destination.Description, &destination.Price)
	if err == sql.ErrNoRows {
		// Return a custom error response if the record is not found
		return Response{Status: http.StatusNotFound, Message: "Not Found"}, nil
	} else if err != nil {
		return Response{}, err
	}

	res.Status = http.StatusOK
	res.Message = "OK"
	res.Data = destination

	return res, nil
}

func StoreDestination(destination_name string, image string, city string, description string, price int) (Response, error) {
	var res Response

	v := validator.New()

	destination := Destination{
		DestinationName: destination_name,
		Image:           image,
		City:            city,
		Description:     description,
	}

	err := v.Struct(destination)

	if err != nil {
		return res, err
	}

	con := db.CreateConnection()

	sqlStatement := "INSERT INTO destination(destination_name, image, city_id, description, price) VALUES (?, ?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(destination_name, image, city, description, price)

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

func UpdateDestination(id int, destination_name string, image string, city string, description string, price int) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "UPDATE destination SET destination_name = ?, image = ?, city_id= ?, description = ?, price = ? WHERE id= ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(destination_name, image, city, description, price, id)
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

func DeleteDestination(id int) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "DELETE FROM destination WHERE id = ?"

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
