go mod init github.com/bartmika/osin-example

go get github.com/google/uuid
go get github.com/rs/cors
go get github.com/spf13/cobra
go get github.com/go-redis/redis  
go get github.com/dgrijalva/jwt-go
go get golang.org/x/crypto/bcrypt
go get github.com/lib/pq
go get gopkg.in/guregu/null.v4
go get golang.org/x/oauth2
go get github.com/ShaleApps/osinredis
go get github.com/openshift/osin


```bash
export OSIN_DB_HOST=localhost
export OSIN_DB_PORT=5432
export OSIN_DB_USER=golang
export OSIN_DB_PASSWORD=123password
export OSIN_DB_NAME=osinexample_db
export OSIN_APP_ADDRESS=127.0.0.1:8000
export OSIN_APP_SECRET_KEY=pass-secret-1234566-please-change-me
```

Register an identity with our system.
```
go run main.go register -b=Bart -c=Mika -d=bart@mikasoftware.com -e=123password -f=en
```

Here is how you do password based grant
```
go run main.go password --email=bart@mikasoftware.com --password=123password
```


Simple login, run and then copy+paste the result export to the terminal

```
go run main.go login --email=bart@mikasoftware.com --password=123password
```

Check we are able to access our protected resource

```
go run main.go tenant_retrieve --id=1
```

Let's refresh our access token with our refresh token.

```
go run main.go refresh_token
```

Get our token from the client credentials.

```
go run main.go client_credential
```
