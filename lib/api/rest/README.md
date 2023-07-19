# Routes

## PubSub

```[GET] /pubsub/topics/```

Gets a list of topics

```[GET] /pubsub/topics/:topic/subscribers```

Gets a list of subscribers


```[DELETE] /pubsub/topics/:topic```

Deletes a topic

``` [POST] /pubsub/publisher/publish```

Sends a new publication

Body of the request should be in the following format
```json
{
    "source": "gmail",
    "template": "greeting",  
    "users": ["user"],
    "meta": {
        "from": "Mailer Service"
    },
    "at": 1667331860 // a unix timestamp
}
```

```[GET] /pubsub/publisher/:id```

Gets the publication information

```[GET] /pubsub/publisher/:id/report```

Gets the publication report

```[GET] /pubsub/publisher/:id/report/:user_id```

Get personal report

```[GET] /pubsub/publisher/reports```

Gets a list of reports

```[GET] /pubsub/publisher/```

Gets a list of publications

```[GET] /pubsub/subscriptions/:id```

Gets subscriptions of the specified user

```[GET] /pubsub/subscriptions/:id/:topic```

Gets a user subscription

```[POST] /pubsub/subscriptions/:id/:topic```

Subscribes the user to the topic

```[PUT] /pubsub/subscriptions/:id/:topic```

Updates the specified subscription

```[DELETE] /pubsub/subscriptions/:id/:topic```

Deletes the specified subscription 

```[DELETE] /pubsub/subscriptions/:id```

Deletes all user's subscriptions

```[GET] /templates/```

Gets templates

``` [GET] /templates/:id```

Gets the specified template

```[DELETE] /templates/:id```

Deletes the specified template

```[POST] /templates/```

Creates a new template

Body of the request should be in the following format:
```json
{
    "name": "greeting",
    "data": "Hello, <b>{{ .user.Name}} {{ .user.Surname}}</b>.",
    "sub": {
        "subject": {
            "data": "Welcome to Mailer!"
        },
        "from": {
            "data": "{{.meta.from}}"
        }
    },
    "meta": {
        "headers": {
            "Content-Type": "text/html"
        }
    }
}
```

```[PUT] /templates/```

Updates the template

```[GET] /users/```

Gets users

```[GET] /users/:id```

Gets the specified user

```[DELETE] /users/:id```

Deletes the specified user

```[PUT] /users/:id```

Updates the specified user

```[POST] /users/```

Creates a new user

Body of the request should be in the following format:

```json
{
    "name": "Roman",
    "surname": "Ischenko",
    "email": "mail@gmail.com"
}
```

```[GET] /sources/```

Gets sources

```[GET] /sources/:name```

Gets the source by the specified name

```[DELETE] /sources/:name```

Deletes the specified source

```[POST] /sources/```

Creates a new source

Body of the request should be in the following format:
```json
{
    "name": "gmail",
    "address": "mail@gmail.com",
    "password": "password",
    "host": "smtp.gmail.com",
    "port": 587
}
```

```[PUT] /sources/```

Updates the source