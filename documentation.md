# As a developer, I want to create an application.

Start your server in one terminal.

```
go run main.go serve
```

In another terminal, run the following. Follow the on-screen instructions.

```
go run main.go login -d=demo@demo.com -e=123password
```

Afterwards run the following.

```
go run main.go application_create -a="Third Party" -b="Demonstration purposes only" -c=http://demo.com -d=all -e=http://127.0.0.1:8001/appauth/code -f=https://g.com/img.png -g=1
```

You should get something like this for example:

```
"client_id":"4a9fb6c06e179980",
"client_secret":"e2d6c4d7b09796dc5254063f77b5282f0dc1218217f241f0fc098c8375613de7f8a0665ca090a27814712d10b04aefa77ef5322014d488bb5428be400bf391285c87de0e7a77aaab1eee23498a8de63b95e445540a8c3d3bca5d0c65332d2dbf9ab00597d02c98f6890db53b02bffd6d03cb134da71ebf5d585fbf7bf9bf4d1",
```

# As a user, I want my account authorized to use by a third party application

In a new terminal, run the following code:

```
go run main.go serve --client_id=4a9fb6c06e179980 \
                     --client_secret=e2d6c4d7b09796dc5254063f77b5282f0dc1218217f241f0fc098c8375613de7f8a0665ca090a27814712d10b04aefa77ef5322014d488bb5428be400bf391285c87de0e7a77aaab1eee23498a8de63b95e445540a8c3d3bca5d0c65332d2dbf9ab00597d02c98f6890db53b02bffd6d03cb134da71ebf5d585fbf7bf9bf4d1 \
                     --redirect_uri=http://127.0.0.1:8001/appauth/code \
                     --authorize_uri=http://localhost:8000/authorize \
                     --token_url=http://localhost:8000/token
```

Then load up your browser to [http://localhost:8001/](http://localhost:8001/). Follow the instructions and the code should work.
