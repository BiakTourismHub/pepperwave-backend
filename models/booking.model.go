package models

import (
	"net/http"

	"github.com/bryansamperura/ticket-booking/db"
	"github.com/go-playground/validator/v10"
)

type Booking struct {
	Id              int    `json:"id"`
	CustomerName    string `json:"customer_name"`
	Qty             int    `json:"qty"`
	DestinationName string `json:"destination_name"`
	Price           int    `json:"price"`
	TanggalBooking  string `json:"booking_date"`
}

type BookingRequest struct {
	CustomerID     int    `json:"customer_id"`
	Qty            int    `json:"qty"`
	DestinationID  int    `json:"destination_id"`
	TanggalBooking string `json:"booking_date"`
}

func FindAllBooking() (Response, error) {
	var obj Booking
	var arrObj []Booking
	var res Response

	con := db.CreateConnection()

	sqlStatement := `SELECT 
						booking.id,
						customers.fullname, 
						booking.qty, 
						destination.destination_name, 
						destination.price, 
						booking.booking_date 
					FROM booking 
					JOIN 
						destination ON destination.id = booking.destination_id 
					JOIN 
						customers ON customers.id = booking.customer_id`

	rows, err := con.Query(sqlStatement)

	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.CustomerName, &obj.Qty, &obj.DestinationName, &obj.Price, &obj.TanggalBooking)
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

func FindBookingById(customer_id int) (Response, error) {
	var obj Booking
	var arrObj []Booking
	var res Response

	con := db.CreateConnection()

	sqlStatement := `SELECT 
						booking.id,
						customers.fullname, 
						booking.qty, 
						destination.destination_name, 
						destination.price, 
						booking.booking_date 
					FROM booking 
					JOIN 
						destination ON destination.id = booking.destination_id 
					JOIN 
						customers ON customers.id = booking.customer_id
					WHERE booking.customer_id = ?`

	rows, err := con.Query(sqlStatement, customer_id)

	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.CustomerName, &obj.Qty, &obj.DestinationName, &obj.Price, &obj.TanggalBooking)
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

func StoreBooking(customer_id int, qty int, destination_id int, booking_date string) (Response, error) {
	var res Response

	v := validator.New()

	booking := BookingRequest{
		CustomerID:     customer_id,
		Qty:            qty,
		DestinationID:  destination_id,
		TanggalBooking: booking_date,
	}

	err := v.Struct(booking)

	if err != nil {
		return res, err
	}

	con := db.CreateConnection()

	sqlStatement := "INSERT INTO booking(customer_id, qty, destination_id,booking_date) VALUES (?, ?, ?,?)"

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(customer_id, qty, destination_id, booking_date)

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
