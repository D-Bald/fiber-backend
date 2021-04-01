# Backend for a selfmade CMS using [Fiber](https://github.com/gofiber/fiber)

-------------------------

* [Inspired by...](#inspired-by...)
* [API](#api)
* [Database setup](#database-setup)
* [TO DO](#to-do)

-------------------------

## Inspired by...

- [go-fiber/recipes auth-jwt](https://github.com/gofiber/recipes/tree/master/auth-jwt)
- [Quick Start: Golang & MongoDB - Modeling Documents with Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)

## API

| Endpoint              | Method    | Description  |
| :-------------------- | :-------- | :------------------------------------------------------------------------- |
| `/api/sample`         | `GET`     | Health Check |
|                       | `POST`    | Create a new Sample Entry in the `samples` collection. Specify two string fields in the request body.   |
| `/api/auth/`          | `POST`    | Sign in with username or email and password. If it's successful, then generates a token. |
| `/api/user/`          | `GET`     | Return all users present in the `users` collection.  |
|                       | `POST`    | Create a new user.</br> Specify the following attributes in the request body: `username`, `email`, `password`, `names`.   |
| `/api/user/:id`       | `GET`     | Return user with id `:id`.   |
|                       | `PATCH`   | Update user with id `:id`. </br> Currently only `names` can be updated. Flexible Update of all fields is planned.   |
|                       | `DELETE`  | Delete user with id `:id`.</br> Specify user´s password in the request body.   |
| `/contenttypes`       | `GET`     | Return all content types present in the `contenttypes` collection. |
|                       | `POST`    | Create a new content type.</br> Specify the following attributes in the request body: `typename`, `collection`, `field_schema`.</br> By convention the collection should be plural of the typename. The last attribute is a list of key-value pairs specifying name and type of fields, that an content entry of this content type should have.  |
| `/contenttypes/:id`   | `GET`     | Return content type with id `:id`.   |
|                       | `DELETE`  | Delete content type with id `:id`.   |
| `/:type`              | `GET`     | Return all content entries of type `:type` from the corresponding collection.  |
|                       | `POST`    | Create a new content entry of type `:type`.</br> Specify the following attributes in the request body: `title` (string), `published`(bool), `fields`(key-value pairs: field name - field value). |
| `/:type/:id`          | `GET`     | Return content entry with type `:type` and id `:id`.   |
|                       | `DELETE`  | Delete content entry with type `:type`and id `:id`. |

## Database Setup

FILL THIS OUT WHEN CONFIGURATION VIA CONFIG FILE IS AVAILABLE

WATCH OUT: MONGODB NOT SELF-HOSTET => URI FOR ATLAS IS HARDCODED EXCEPT USER, DB NAME AND CREDENTIALS
For self-hosted DB. Adjust MongoURI [here](https://github.com/D-Bald/fiber-backend/blob/0f15612d722b1bbc8c7a5356fff78ae308da2c71/database/connect.go#L24)

## TO DO

* Initialization of content types on start up with blogposts and events
* User Handler:
   * Update not only "names" but also username, password...
* Content Handler:
   * add `PATCH` Method for updates
* Configuration via Config file
* API Endpoints for File/Media Uploading
* Add in API-Table which Routes are protected

