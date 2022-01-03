# OSIN Example
The purpose of this repository is to provide Golang example code of an identity provider using oAuth 2.0 to grant access to users and [third-party developers](https://github.com/bartmika/osin-thirdparty-example). This code uses the [`openshift/osin`](https://github.com/openshift/osin) Golang oAuth 2.0 library.

## Setup

1. Clone the project.

2. Install the dependencies

    ```
    go get ./...
    ```

3. Please create a default client which you can use to a frontend.

    ```
    go run main.go add_client --client_id=frontend \
                              --client_secret=pleasechangethisnow \
                              --redirect_uri=http://127.0.0.1:8002/appauth/code
    ```

3. Setup the database.

4. Run the following environment variables in your terminal:

    ```bash
    export OSIN_DB_HOST=localhost
    export OSIN_DB_PORT=5432
    export OSIN_DB_USER=golang
    export OSIN_DB_PASSWORD=123password
    export OSIN_DB_NAME=osinexample_db
    export OSIN_APP_ADDRESS=http://127.0.0.1:8000
    export OSIN_APP_SECRET_KEY=pass-secret-1234566-please-change-me
    export OSIN_APP_FRONTEND_CLIENT_ID=frontend
    export OSIN_APP_FRONTEND_CLIENT_SECRET=pleasechangethisnow
    export OSIN_APP_FRONTEND_RETURN_URL=http://127.0.0.1:8001/appauth/code
    ```

5. Run the server.

    ```bash
    go run main.go serve
    ```

6. You are ready to use your server.

## Notes

Register an identity with our system.
```
go run main.go register -b=Bart -c=Mika -d=demo@demo.com -e=123password -f=en
```

Here is how you do password based grant
```
go run main.go osin_password --email=demo@demo.com \
                             --password=123password \
                             --client_id=frontend \
                             --client_secret=pleasechangethisnow \
                             --redirect_uri=http://127.0.0.1:8001/appauth/code \
                             --authorize_uri=http://localhost:8000/authorize \
                             --token_url=http://localhost:8000/token
```


Simple login, run and then copy+paste the result export to the terminal

```
go run main.go login --email=demo@demo.com --password=123password
```

Check we are able to access our protected resource

```
go run main.go tenant_retrieve --id=1
```

Let's refresh our access token with our refresh token.

```
go run main.go osin_refresh_token --client_id=frontend \
                                  --client_secret=pleasechangethisnow \
                                  --redirect_uri=http://127.0.0.1:8001/appauth/code \
                                  --authorize_uri=http://localhost:8000/authorize \
                                  --token_url=http://localhost:8000/authorize
```

Get our token from the client credentials.

```
go run main.go osin_client_credential --client_id=frontend \
                                      --client_secret=pleasechangethisnow \
                                      --redirect_uri=http://127.0.0.1:8001/appauth/code \
                                      --authorize_uri=http://localhost:8000/authorize \
                                      --token_url=http://localhost:8000/authorize
```

Get new refresh API

```
go run main.go refresh_token --refresh_token=xxx --grant_type=refresh_token
```

## License
Made with ❤️ by [Bartlomiej Mika](https://bartlomiejmika.com).   
The project is licensed under the [Unlicense](LICENSE).

Third party libraries and resources:

* [github.com/openshift/osin](https://github.com/openshift/osin) (BSD-3-Clause License) is used for oAuth 2.0 server implementation.
* [github.com/ShaleApps/osinredis](https://github.com/ShaleApps/osinredis) (MIT) is used for the oAuth 2.0 server session storage handling.
* [go-oauth2/oauth2](https://github.com/go-oauth2/oauth2) (MIT) was not used but was a valuable learning resource.
* [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) (BSD-3-Clause License) is used for client side oAuth 2.0 library.
* [github.com/google/uuid](https://github.com/google/uuid) (BSD-3-Clause License) is used for generating UUID values.
* [github.com/rs/cors](https://github.com/rs/cors) (MIT) is used for dealing with CORS headers for every request made to this API server.
* [github.com/spf13/cobra](https://github.com/spf13/cobra) (Apache-2.0 License) is used to structure this application using a CLI.
* [github.com/go-redis/redis](https://github.com/go-redis/redis) (BSD-2-Clause License) is used for `redis` handling.
* [golang.org/x/crypto/bcrypt](https://golang.org/x/crypto/bcrypt) (BSD-3-Clause License) is used for password hashing.
* [github.com/lib/pq](https://github.com/lib/pq) (OTHER) is used for `postgres` handling.
* [gopkg.in/guregu/null.v4](https://gopkg.in/guregu/null.v4) (BSD-2-Clause License) is used for null fields in marshal/unmarshal operations.
* [rfc6749 - The OAuth 2.0 Authorization Framework](https://datatracker.ietf.org/doc/html/rfc6749) used as specs.
