# Backend for a mini selfmade headless CMS using [Fiber](https://github.com/gofiber/fiber)

-------------------------

## Content

* [Inspired by...](#inspired-by...)
* [API](#api)
* [Workflows](#workflows)
   * [Create content and content types](#create-content-and-content-types)
   * [Create users and manage roles](#create-users-and-manage-roles)
* [Database setup](#database-setup)
* [TO DO](#to-do)

-------------------------

## Inspired by...

* [go-fiber/recipes auth-jwt](https://github.com/gofiber/recipes/tree/master/auth-jwt)
* [Quick Start: Golang & MongoDB - Modeling Documents with Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)

## API

| Endpoint              | Method    | Authentification required   | Description  |
| :-------------------- | :-------- | :-------------------------- | :------------------------------------------- |
| `/api/sample`         | `GET`     | &cross;                     | Health Check |
|                       | `POST`    | &cross;                     | Create a new Sample Entry in the `samples` collection. Specify two string fields in the request body.   |
| `/api/auth/`          | `POST`    | &cross;                     | Sign in with username or email (`identity`) and `password`. If it's successful, then generates a token. |
| `/api/user/`          | `GET`     | &cross;                     | Return all users present in the `users` collection.  |
|                       | `POST`    | &cross;                     | Create a new user.</br> Specify the following attributes in the request body: `username`, `email`, `password`, `names`.   |
| `/api/user/:id`       | `GET`     | &cross;                     | Return user with id `id`.   |
|                       | `PATCH`   | &check;                     | Update user with id `id`. </br> Currently only `names` can be updated. Flexible Update of all fields is planned.   |
|                       | `DELETE`  | &check;                     | Delete user with id `id`.</br> Specify user´s password in the request body.   |
| `/api/contenttypes`       | `GET`     | &cross;                     | Return all content types present in the `contenttypes` collection. |
|                       | `POST`    | &check;                     | Create a new content type.</br> Specify the following attributes in the request body: `typename`, `collection`, `field_schema`.</br> By convention the collection should be plural of the typename. The last attribute is a list of key-value pairs specifying name and type of fields, that an content entry of this content type should have.</br> Example: ```{"typename": "Event", "collection": "events", "field_schema": {"date": "time.Time", "place": "string"}}```  |
| `/api/contenttypes/:id`   | `GET`     | &cross;                     | Return content type with id `:id`.   |
|                       | `DELETE`  | &check;                     | Delete content type with id `:id`.   |
| `/api/:content`           | `GET`     | &cross;                     | Return all content entries of the content type, where `content` is the corresponding collection. By convention this should be plural of the `typename`.</br> For the previous example: `content` has to be set to `events`.   |
|                       | `POST`    | &check;                     | Create a new content entry of the content type, where `content` is the corresponding collection.</br> Specify the following attributes in the request body: `title` (string), `published`(bool), `fields`(key-value pairs: field name - field value). |
| `/api/:content/:id`       | `GET`     | &cross;                     | Return content entry with id `id` of the content type, where `content` is the corresponding collection.   |
|                       | `DELETE`  | &check;                     | Delete content entry with id `id` of the content type, where `content` is the corresponding collection.   |


## Workflows
### Create content and content types
The Content Types *event* and *blogpost* are preset and you can start adding entries on those routes (`/api/events` or `/api/blogposts`).
If you want to create a custom Content Type, first use the `/api/contenttypes`endpoint, because the `/api/:content` route is validated by a lookup in the `contenttypes` collection. The mongoDB collections for new types are created automatically on first content insertion.

### Create users and manage roles
NOT YET IMPLEMENTED see [TO DO](#to-do)

## Database setup

FILL THIS OUT WHEN CONFIGURATION VIA CONFIG FILE IS AVAILABLE

WATCH OUT: MONGODB NOT SELF-HOSTET => URI FOR ATLAS IS HARDCODED EXCEPT USER, DB NAME AND CREDENTIALS.

For self-hosted DB adjust [mongoURI in this line](https://github.com/D-Bald/fiber-backend/blob/0f15612d722b1bbc8c7a5356fff78ae308da2c71/database/connect.go#L24)

## TO DO

* User Handler:
   * Update not only "names" but also username, password... (by using Query-Params in URL instead of specification in request body?)
   * Manage User Roles
* Content Handler:
   * add `PATCH` Method to update of any content entry
   * Query content and types by title/Name, tags, (and field-values?)
   * Distinguish between user roles: just admins can reach routes, that are now protected, but anyone can create a user. Create a default admin user on start.
* API Endpoints for File/Media upload
* Configuration via Config file
* Validation schemas for Input Data (https://docs.gofiber.io/guide/validation)
