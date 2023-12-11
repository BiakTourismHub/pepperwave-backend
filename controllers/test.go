package controllers

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string
	Data    interface{}
}

func Test(c echo.Context) error {
	var response Response
	var isSuccess bool
	var fileType, fileName string
	var fileSize int64

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusOK, Response{Message: "Upload Failed", Data: nil})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusOK, Response{Message: "Upload Failed", Data: nil})
	}
	defer src.Close()

	fileByte, _ := io.ReadAll(src)
	fileType = http.DetectContentType(fileByte)

	if fileType == "image/jpg" {
		fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
	} else if fileType == "image/jpeg" {
		fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpeg"
	} else if fileType == "image/png" {
		fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	}

	err = os.WriteFile(fileName, fileByte, 0777)
	if err != nil {
		isSuccess = false
	} else {
		fileSize = file.Size
		isSuccess = true
	}

	if isSuccess {
		response = Response{
			Message: "Upload Success",
			Data: struct {
				Filename string
				Filetype string
				Filesize int64
			}{
				Filename: fileName,
				Filetype: fileType,
				Filesize: fileSize,
			},
		}
	}

	return c.JSON(http.StatusOK, response)
}
