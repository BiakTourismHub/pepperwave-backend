package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bryansamperura/ticket-booking/models"
	"github.com/labstack/echo/v4"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	AccountID int    `json:"uid"`
	Role      string `json:"role"`
}

type SignUpRequest struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Register Register new customer
// @Summary Regster new customer
// @Description register a new user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param data body SignUpRequest true "Register Data"
// @Success 200 {object} models.Response
// @Router /register [post]
func Register(c echo.Context) error {
	request := new(models.AuthRegisterRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	cust, err := models.StoreCustomer(request.Fullname, request.Email, request.Phone)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var accountID int64

	if data, ok := cust.Data.(map[string]int64); ok {
		id := data["id"]
		accountID = id
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	request.Password = string(hashPassword)
	request.Role = "customer"

	account, err := models.StoreAccount(request.Email, request.Password, request.Role, int(accountID))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, account)
}

// Login Make authentication
// @Summary Login customer
// @Description make authentication for the users
// @Tags Auth
// @Accept  json
// @Param data body models.AuthRequest true "Login Data"
// @Success 200 {object} models.Response
// @Router /login [post]
func Login(c echo.Context) error {
	request := new(models.AuthRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	user, err := models.FindUserByEmail(request.Email)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if user.Status == 404 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Bad Credentials"})
	}

	userData := user.Data.(models.User)

	AccID, err := strconv.Atoi(userData.AccountID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(request.Password))

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	// Generate a PASETO token
	token, err := generateToken(AccID, userData.Role)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": "Bearer " + token})
}

func Protected(c echo.Context) error {
	userID := c.Get("uid").(int)
	role := c.Get("role").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "You accessed a protected route!",
		"user_id":  userID,
		"username": role,
	})
}

func generateToken(userID int, role string) (string, error) {

	token, err := paseto.NewV2().Encrypt([]byte("YELLOW SUBMARINE, BLACK WIZARDRY"), CustomClaims{
		AccountID: userID,
		Role:      role,
	}, nil)

	if err != nil {
		return "", err
	}

	return token, nil
}

// GetAccountInfo Get account info
// @Summary Show account info
// @Description Show account info
// @Security Bearer
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} CustomClaims
// @Router /account-info [get]
func GetAccountInfo(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token is missing")
	}

	// Split the token string by space to get the actual token value
	tokenParts := strings.Fields(token)
	if len(tokenParts) != 2 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
	}

	// Extract the actual token value
	actualToken := tokenParts[1]

	var claims CustomClaims

	if err := paseto.NewV2().Decrypt(actualToken, []byte("YELLOW SUBMARINE, BLACK WIZARDRY"), &claims, nil); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	return c.JSON(http.StatusOK, claims)
}
