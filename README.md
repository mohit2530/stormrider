# Auth app
Authentication app built with JWT token with api gateway for authentication. Application simulates the JWT web token workflow. 

## Application use case

1. Use `.env` file for each microservice. Default `.env` file consists of `TOKEN_ENDPOINT`. Use any available port. General ports in use fall under 3xxx domain or 5xxxx domain.

2. Simulate user, password authentication from db. For testing purposes, use `admin` and `password`. `Note:` In general use, the username, password should be retrieved from the db. Since we do not have that capability, we are using default `admin` and `password`. `This is NOT IDEAL in production env`. It is up to the client to configure this however they see fit.

2. After user is succesfully logged in, `LICENSE` file should populate with the newly generated license key. All downstream api will need this api key. Client does not need to do anything at this time. Again this is not ideal. key must be safely stored in db and retrieved as a masked entity.

3. Test by using `sample.http` REST CLIENT. All `VS Code` users should install the `REST CLIENT` extension and test with the provided api. After the user is logged in from the token-generator api, update the retrieved token in the Authorization header (Some clients automatically do it for you. )

### Developer notes

1. Run `go mod tidy` to cleanup and install required deps.
2. Run `go run main.go` for `api_gateway` and `jwt-tokens`.