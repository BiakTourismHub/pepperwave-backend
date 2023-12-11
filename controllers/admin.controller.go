package controllers

import (
	"net/http"
	"strconv"

	"github.com/bryansamperura/ticket-booking/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// FetchAllAdmin returns a list of all admin
// @Summary Get a list of all admin
// @Description Retrieve a list of all admin
// @Tags Admin
// @Produce json
// @Success 200 {array} models.Admin
// @Failure 500 {object} models.HTTPError
// @Router /admin [get]
func FetchAllAdmin(c echo.Context) error {
	result, err := models.FindAllAdmin()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GetAdminById return admin by ID
// @Summary Get admin by id
// @Description Returns the admin with the given id
// @Tags Admin
// @Param id path int true "Admin ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Admin
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /admin/{id} [get]
func GetAdminById(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.FindAdminById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusOK, result)
}

// StoreAdmin stores admin data
// @Summary Create a new admin
// @Description Save a new admin to the database
// @Tags Admin
// @Accept json
// @Consumes json
// @Param admin body models.AdminRequest true "Admin Name"
// @Success 201 {object} models.Admin
// @Failure 500 {object} models.HTTPError
// @Router /admin [post]
func StoreAdmin(c echo.Context) error {

	request := new(models.AdminRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	result, err := models.StoreAdmin(request.FullName, request.Email, request.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var accountID int64

	if data, ok := result.Data.(map[string]int64); ok {
		id := data["id"]
		accountID = id
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	request.Password = string(hashPassword)
	request.Role = "admin"

	account, err := models.StoreAccount(request.Email, request.Password, request.Role, int(accountID))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, account)

}

// UpdateAdmin stores admin data
// @Summary Update admin
// @Description Updates an existing admin in the database
// @Tags Admin
// @Accept json
// @Consumes json
// @Param id path int true "Admin ID"
// @Param admin body models.AdminRequest true "Admin Name"
// @Success 204 {object} string
// @Failure 500 {object} models.HTTPError
// @Router /admin/{id} [put]
func UpdateAdmin(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	request := new(models.AdminRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	result, err := models.UpdateCustomer(id, request.FullName, request.Email, request.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// DeleteAdmin delete admin by id
// @Summary Delete admin
// @Description Deletes an existing admin from the database
// @Tags Admin
// @Produce json
// @Param id path int true "Admin ID"
// @Success 204 {object} string
// @Failure 400 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /admin/{id} [delete]
func DeleteAdmin(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.DeleteAdmin(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusNoContent, result)
}
