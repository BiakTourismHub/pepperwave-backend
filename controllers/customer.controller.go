package controllers

import (
	"net/http"
	"strconv"

	"github.com/bryansamperura/ticket-booking/models"
	"github.com/labstack/echo/v4"
)

// FetchAllCustomers returns a list of all customers
// @Summary Get a list of all customers
// @Description Retrieve a list of all customers
// @Tags Customer
// @Security Bearer
// @Produce json
// @Success 200 {array} models.Customer
// @Failure 500 {object} models.HTTPError
// @Router /customers [get]
func FetchAllCustomers(c echo.Context) error {
	result, err := models.FindAllCustomer()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GetCustomerById return customer by ID
// @Summary Get customer by id
// @Description Returns the customer with the given id
// @Tags Customer
// @Security Bearer
// @Param id path int true "Customer ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Customer
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /customer/{id} [get]
func GetCustomerById(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.FindCustomerById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusOK, result)
}

// StoreCustomer stores customer data
// @Summary Create a new customer
// @Description Save a new customer to the database
// @Tags Customer
// @Security Bearer
// @Accept json
// @Consumes json
// @Param customer body models.CustomerRequest true "Customer Name"
// @Success 201 {object} models.Customer
// @Failure 500 {object} models.HTTPError
// @Router /customer [post]
func StoreCustomer(c echo.Context) error {

	request := new(models.CustomerRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	result, err := models.StoreCustomer(request.FullName, request.Email, request.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// UpdateCustomer stores customer data
// @Summary Update customer
// @Description Updates an existing customer in the database
// @Tags Customer
// @Security Bearer
// @Accept json
// @Consumes json
// @Param id path int true "Customer ID"
// @Param customer body models.CustomerRequest true "Customer Name"
// @Success 204 {object} string
// @Failure 500 {object} models.HTTPError
// @Router /customer/{id} [put]
func UpdateCustomer(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	request := new(models.CustomerRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	result, err := models.UpdateCustomer(id, request.FullName, request.Email, request.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// DeleteCustomer delete customer by id
// @Summary Delete Customer
// @Description Deletes an existing customer from the database
// @Tags Customer
// @Security Bearer
// @Produce json
// @Param id path int true "Customer ID"
// @Success 204 {object} string
// @Failure 400 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /customer/{id} [delete]
func DeleteCustomer(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.DeleteCustomer(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusNoContent, result)
}
