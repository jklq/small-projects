package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"

	"github.com/go-playground/validator"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jklq/small-projects/go-jwt-bare/user"
)

type LoginParams struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func isValidParams(c *fiber.Ctx, params interface{}) bool {
	err := validate.Struct(params)

	if err != nil {
		return false
	}
	return true
}

var userq *user.DBQuerier

func init() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://postgres:123@localhost:5432/sandbox")

	if err != nil {
		panic(err)
	}

	userq = user.NewQuerier(conn)
}

func main() {
	app := fiber.New()

	// Login route
	app.Post("/login", login)

	// Unauthenticated route
	app.Get("/", accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	// Restricted Routes
	app.Get("/protected", restricted)

	app.Listen(":3000")
}

var validate = validator.New()

func login(c *fiber.Ctx) error {

	loginParams := new(LoginParams)

	if err := c.BodyParser(loginParams); err != nil {
		return err
	}

	// Validate input
	if !isValidParams(c, loginParams) {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	results, err := userq.GetUserByEmail(context.Background(), loginParams.Email)

	if err != nil {
		panic(err)
	}

	if len(results) == 0 || results[0].Password != loginParams.Password {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"email": results[0].Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	return c.SendString("Welcome " + email)
}
