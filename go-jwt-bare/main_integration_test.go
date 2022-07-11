package main

import (
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestHttpBin(t *testing.T) {
	var jwt string
	Test(t,
		Description("GET accessible endpoint"),
		Get("http://localhost:3000/"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Accessible"),
	)
	Test(t,
		Description("GET protected endpoint is a bad request"),
		Get("http://localhost:3000/protected"),
		Expect().Status().Equal(http.StatusBadRequest),
	)
	Test(t,
		Description("POST /login with missing password"),
		Post("http://localhost:3000/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"email": "wrong"}),
		Expect().Status().Equal(http.StatusUnprocessableEntity),
	)
	Test(t,
		Description("POST /login with missing email"),
		Post("http://localhost:3000/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"password": "wrong"}),
		Expect().Status().Equal(http.StatusUnprocessableEntity),
	)
	Test(t,
		Description("POST /login with wrong details"),
		Post("http://localhost:3000/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]string{"email": "wrong", "password": "wrong"}),
		Expect().Status().Equal(http.StatusNotFound),
	)
	Test(t,
		Description("POST /login with correct details"),
		Post("http://localhost:3000/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]string{"email": "email@email.com", "password": "password"}),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("token"),
		Store().Response().Body().JSON().JQ(".token").In(&jwt),
	)
	Test(t,
		Description("POST /protected with valid JWT"),
		Get("http://localhost:3000/protected"),
		Send().Headers("Authorization").Add("Bearer "+jwt),
		Expect().Status().Equal(http.StatusOK),
	)
}
