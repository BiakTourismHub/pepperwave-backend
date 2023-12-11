package models

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bryansamperura/ticket-booking/db"
	"github.com/go-playground/validator/v10"
)

type City struct {
	Id       int    `json:"id"`
	CityName string `json:"city"`
}

type CityRequest struct {
	CityName string `json:"city"`
}

func FindAllCity() (Response, error) {
	var obj City
	var arrObj []City
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM cities"

	rows, err := con.Query(sqlStatement)

	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.CityName)
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
func FindCityById(id int) (Response, error) {
	var city City
	var res Response

	con := db.CreateConnection()

	sqlStatement := "SELECT * FROM cities WHERE id= ?"

	row := con.QueryRow(sqlStatement, id)

	err := row.Scan(&city.Id, &city.CityName)
	if err == sql.ErrNoRows {
		// Return a custom error response if the record is not found
		return Response{Status: http.StatusNotFound, Message: "Not Found"}, nil
	} else if err != nil {
		return Response{}, err
	}

	res.Status = http.StatusOK
	res.Message = "OK"
	res.Data = city

	return res, nil
}

func StoreCity(city string) (Response, error) {
	var res Response

	v := validator.New()

	cities := City{
		CityName: city,
	}

	err := v.Struct(cities)

	if err != nil {
		return res, err
	}

	con := db.CreateConnection()

	sqlStatement := "INSERT INTO cities(city_name) VALUES (?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(city)

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

func UpdateCity(id int, city string) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "UPDATE cities SET city_name=? WHERE id= ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(city, id)
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

func DeleteCity(id int) (Response, error) {
	var res Response

	con := db.CreateConnection()

	sqlStatement := "DELETE FROM cities WHERE id = ?"

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
