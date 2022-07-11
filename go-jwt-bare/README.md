# Go Auth (bare)

## Specifications

The API has two endpoints:

* GET /
    * an unprotected endpoints
* GET /protected
    * a protected endpoints where you need a valid JWT bearer token
* POST /login
    * a login endpoint where you can supply correct email and password and get a JWT

The user's email will be stored as an JWT claim.

The server accesses a postgres db to store and get user info.

> NOTE: Passwords are not hashed, as it is outside the scope. But this should of course be done in any production environmnet.

## How to get it working

* Specify postgres details in dbconfig.yml file.
* Install [the migrate cli](https://github.com/golang-migrate/migrate) and run the up migrations
* Run:
```
go get
go run .
```
To get the token, log into [localhost:3000/login](http://localhost:3000/login) with a json object.
It needs an email and a password. 
The database should be seeded with a user with the following credentials:
```
{
	"email": "email@email.com",
	"password": "password"
}
```