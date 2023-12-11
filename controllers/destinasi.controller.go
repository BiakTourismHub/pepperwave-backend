package controllers

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bryansamperura/ticket-booking/models"
	"github.com/labstack/echo/v4"
)

// FetchAllDestination returns a list of all destination
// @Summary Get a list of all destination
// @Description Retrieve a list of all destination
// @Tags Destinations
// @Security Bearer
// @Produce json
// @Success 200 {array} models.Destination
// @Failure 500 {object} models.HTTPError
// @Router /destination [get]
func FetchAllDestination(c echo.Context) error {
	result, err := models.FindAllDestination()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GetDestinationById return destination by ID
// @Summary Get destination by id
// @Description Returns the destination with the given id
// @Tags Destinations
// @Security Bearer
// @Param id path int true "Destination ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Destination
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /destination/{id} [get]
func GetDestinationById(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.FindDestinationById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusOK, result)
}

// StoreDestination stores destination data
// @Summary Create a new destination
// @Description Save a new destination to the database
// @Tags Destinations
// @Security Bearer
// @Accept multipart/form-data
// @Consumes json
// @Param destination_name formData string true "Destination Name"
// @Param city_id formData int true "City ID"
// @Param description formData string true "Description"
// @Param price formData int true "Price"
// @Param image formData file true "Image"
// @Success 201 {object} models.Destination
// @Failure 500 {object} models.HTTPError
// @Router /destination [post]
func StoreDestination(c echo.Context) error {
	var fileType, fileName string

	destinationName := c.FormValue("destination_name")
	cityID := c.FormValue("city_id")
	description := c.FormValue("description")
	price := c.FormValue("price")
	file, err := c.FormFile("image")

	priceInt, err := strconv.Atoi(price)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing image field"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusOK, Response{Message: "Upload Failed", Data: nil})
	}
	defer src.Close()

	fileByte, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read file"})
	}

	fileType = http.DetectContentType(fileByte)
	fileName = generateFileName(fileType)

	err = os.WriteFile(fileName, fileByte, 0777)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	result, err := models.StoreDestination(destinationName, fileName, cityID, description, priceInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// UpdateDestination stores destination data
// @Summary Update destination
// @Description Updates an existing destination in the database
// @Tags Destinations
// @Security Bearer
// @Accept multipart/form-data
// @Consumes json
// @Param id path int true "destination ID"
// @Param destination_name formData string true "Destination Name"
// @Param city_id formData int true "City ID"
// @Param description formData string true "Description"
// @Param price formData int true "Price"
// @Param image formData file true "Image"
// @Success 204 {object} string
// @Failure 500 {object} models.HTTPError
// @Router /destination/{id} [put]
func UpdateDestination(c echo.Context) error {
	var fileType, fileName string

	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	destinationName := c.FormValue("destination_name")
	cityID := c.FormValue("city_id")
	description := c.FormValue("description")
	price := c.FormValue("price")
	priceInt, err := strconv.Atoi(price)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	file, err := c.FormFile("image")

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing image field"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusOK, Response{Message: "Upload Failed", Data: nil})
	}
	defer src.Close()

	fileByte, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read file"})
	}

	fileType = http.DetectContentType(fileByte)
	fileName = generateFileName(fileType)

	err = os.WriteFile(fileName, fileByte, 0777)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	result, err := models.UpdateDestination(id, destinationName, fileName, cityID, description, priceInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// DeleteDestination delete destination by id
// @Summary Delete destination
// @Description Deletes an existing destination from the database
// @Tags Destinations
// @Security Bearer
// @Produce json
// @Param id path int true "Destination ID"
// @Success 204 {object} string
// @Failure 400 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /destination/{id} [delete]
func DeleteDestination(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.DeleteDestination(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusNoContent, result)
}

func generateFileName(fileType string) string {
	return "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + getFileExtension(fileType)
}

func getFileExtension(fileType string) string {
	switch fileType {
	case "image/jpg":
		return ".jpg"
	case "image/jpeg":
		return ".jpeg"
	case "image/png":
		return ".png"
	default:
		return ".dat" // Provide a default extension or handle unsupported types accordingly
	}
}
