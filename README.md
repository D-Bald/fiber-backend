# Backend for a selfmade CMS using [Fiber](https://github.com/gofiber/fiber)

-------------------------
- [Inspired by](#inspired-by)
- [Available Routes](#available-routes)
-------------------------

## Inspired by
- [go-fiber/recipes auth-jwt](https://github.com/gofiber/recipes/tree/master/auth-jwt)
- [Quick Start: Golang & MongoDB - Modeling Documents with Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)

## Available routes

|  Route             |  Method   |  Description    |
|  :---------------- |  :------- |  :------------------------------------------------------------------------------- |
|  `/api/sample`     |  `GET`    |  Health Check   |
|                    |  `POST`   |  Create a new Sample Entry in the `samples` collection. Specify to String Fields in the body  |
|  `/api/auth/`      |  `POST`   |  Sign in with username or email and password. If it's successful, then generates a token   |
|  `/api/user/`      |  `GET`    |  Returns all users present in the collection `users` collection |
|                    |  `POST`   |  Create a new user. You need to specify in the body the following attributes: username, email, password, names  |
|  `/api/user/:id`   |  `GET`    |  Returns user with `id`-field specified by `:id` |
|                    |  `PATCH`  | HIER WEITER  |
|                    |  `DELETE` | HIER WEITER  |

## Usage

FILL THIS OUT WHEN CONFIGURATION VIA CONFIG FILE IS AVAILABLE

WATCH OUT: MONGODB NOT SELF-HOSTET => URI FOR ATLAS IS HARDCODED EXCEPT USER, DB NAME AND CREDENTIALS

## To Do:

- Initialization of content types on start up with blogposts and events
- User Handler:
   - Update not only "names" but also username, password...
- Content Handler:
   - add `PATCH` Method for updates
- Configuration via Config file


