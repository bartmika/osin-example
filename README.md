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
                          --redirect_uri=http://127.0.0.1:8001/appauth/code
    ```

3. Setup the database.

4. Run the following environment variables in your terminal:

    ```bash
    export OSIN_DB_HOST=localhost
    export OSIN_DB_PORT=5432
    export OSIN_DB_USER=golang
    export OSIN_DB_PASSWORD=123password
    export OSIN_DB_NAME=osinexample_db
    export OSIN_APP_ADDRESS=127.0.0.1:8000
    export OSIN_APP_SECRET_KEY=pass-secret-1234566-please-change-me
    export OSIN_APP_FRONTEND_CLIENT_ID=frontend
    export OSIN APP_FRONTEND_CLIENT_SECRET=pleasechangethisnow
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
go run main.go osin_password --email=demo@demo.com --password=123password
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
go run main.go osin_refresh_token
```

Get our token from the client credentials.

```
go run main.go osin_client_credential
```

Get new refresh API

```
go run main.go refresh_token -b=xxx
```
