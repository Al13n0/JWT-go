# JWT-go
This repo contains a baseline application in go that allows the signup and signin of new users. The authenticated users received a JWT token that allow them to access the /protected endpoint


## Setup
Create a file called .env with the following information

```
POSTGRESS_URL=""
SECRET=""
```
For testing the webapp you can use a service like https://www.elephantsql.com/ to manage the database part.
The Secret is used to sign the JWT token. Only the one that knows the secret can verify the signature of the JWT token

## Testing 
Use the postman collection to test the API endpoints.

```
POST /signup ----> endpoint to create a new user

POST /login ----> endpoint for the login

GET /protected ----> protected endpoint you can only access this endpoint with a valid JWT token
```
