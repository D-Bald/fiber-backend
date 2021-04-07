# Backend for a mini selfmade headless CMS using [Fiber](https://github.com/gofiber/fiber)

-------------------------

## Content

- [API](#api)
- [Workflows](#workflows)
   - [Create content and content types](#create-content-and-content-types)
   - [Create users and manage roles](#create-users-and-manage-roles)
- [Database setup](#database-setup)
- [TO DO](#to-do)
- [Thanks to...](#thanks-to...)

-------------------------

## API

| Endpoint              | Method    | Authentification required   | Description  |
| :-------------------- | :-------- | :-------------------------- | :------------------------------------------- |
| `/api/sample`         | `GET`     | &cross;                     | Health Check |
|                       | `POST`    | &cross;                     | Create a new Sample Entry in the `samples` collection. Specify two string fields in the request body.   |
| `/api/auth/`          | `POST`    | &cross;                     | Sign in with username or email (`identity`) and `password`. If it's successful, then generates a token. |
| `/api/user/`          | `GET`     | &cross;                     | Return all users present in the `users` collection.  |
|                       | `POST`    | &cross;                     | Create a new user.</br> Specify the following attributes in the request body: `username`, `email`, `password`, `names`.   |
| `/api/user/:id`       | `GET`     | &cross;                     | Return user with id `id`.   |
|                       | `PATCH`   | &check;                     | Update user with id `id`. </br> If you want to update `role`, you have to be authenticated with a admin-user.  |
|                       | `DELETE`  | &check;                     | Delete user with id `id`.</br> Specify user´s password in the request body.   |
| `/api/contenttypes`       | `GET`     | &cross;                     | Return all content types present in the `contenttypes` collection. |
|                       | `POST`    | &check;                     | Create a new content type.</br> Specify the following attributes in the request body: `typename`, `collection`, `field_schema`.</br> By convention the collection should be plural of the typename. The last attribute is a list of key-value pairs specifying name and type of fields, that an content entry of this content type should have.</br> Example: ```{"typename": "Event", "collection": "events", "field_schema": {"date": "time.Time", "place": "string"}}```  |
| `/api/contenttypes/:id`   | `GET`     | &cross;                     | Return content type with id `:id`.   |
|                       | `DELETE`  | &check;                     | Delete content type with id `:id`.   |
| `/api/:content`           | `GET`     | &cross;                     | Return all content entries of the content type, where `content` is the corresponding collection. By convention this should be plural of the `typename`.</br> For the previous example: `content` has to be set to `events`.   |
|                       | `POST`    | &check;                     | Create a new content entry of the content type, where `content` is the corresponding collection.</br> Specify the following attributes in the request body: `title` (string), `published`(bool), `fields`(key-value pairs: field name - field value). |
| `/api/:content/*`       | `GET`     | &cross;                     | Return content entries filtered by parameters in URL mathing the following regular expression: `[a-z\.\_]+=[a-zA-Z0-9\%]+`</br>
The first group represents the search key and the second the search value. Custom fields can be queried directly so **don't** use dot-notation or similar (see Example 2). Only the whole field value is matched, so no submatches supported. Queries for single Tags are possible. </br> Examples: 1. `/api/events/title=Title%20Test&id=606886f352caea1f9aa86471`</br> 2. `/api/blogposts/text=Hello%20World` |
|                       | `PATCH`   | &check;                     | Update content entry with id `id` of the content type, where `content` is the corresponding collection.  |
|                       | `DELETE`  | &check;                     | Delete content entry with id `id` of the content type, where `content` is the corresponding collection.   |


## Workflows
### Create content and content types
The Content Types *event* and *blogpost* are preset and you can start adding entries on those routes (`/api/events` or `/api/blogposts`). Events have custom fields *description* and *date* whereas blogposts come with *description* and *text*.
If you want to create a custom Content Type, first use the `/api/contenttypes`endpoint, because the `/api/:content` route is validated by a lookup in the `contenttypes` collection. The mongoDB collections for new types are created automatically on first content insertion.
When you update a content entry custom fields can be reached by specifying it as nested object in the request body. JSON example on route: `http://localhost:3000/api/events/606da00b3b00808b74c418c8`
```
{ "fields": {"description": "foo bar"} }
```
Preset fields can be reached directly. JSON example on route `http://localhost:3000/api/events/606da00b3b00808b74c418c8`:
```
{ "tags": ["foo", "bar"] }
```


### Create users and manage roles
The admin user *adminUser* is preset with the password `ADMIN_PASSWORD` from the [.env](https://github.com/D-Bald/fiber-backend/blob/master/.env.sample) file.
Anybody can create a new user. The role is automatically set to *user*. A user can edit the own data i.e. *username*, *email*, *password*, *names*.
Every user with role *admin* can edit any other user and particularly can edit the field *role* of any user. Roles must be updated as array containing all roles as single strings:
```
{ "roles": ["user", "admin"] }
```


## Database setup

FILL THIS OUT WHEN CONFIGURATION VIA CONFIG FILE IS AVAILABLE

WATCH OUT: MONGODB NOT SELF-HOSTET => URI FOR ATLAS IS HARDCODED EXCEPT USER, DB NAME AND CREDENTIALS.

For self-hosted DB adjust [mongoURI in this line](https://github.com/D-Bald/fiber-backend/blob/0f15612d722b1bbc8c7a5356fff78ae308da2c71/database/connect.go#L24)


## TO DO

- User Handler:
   - change the GET endpoint to parse URL parameters to filter users similar to the content handler
- Add Auth Middleware: Distinguish between user roles: just admins can reach routes, that are now protected, but anyone can create a user (By handing the role to the jwt claims). Create a default admin user on start.
- Validation schemas for Input Data (https://docs.gofiber.io/guide/validation)
- Implement file upload
- Configuration via external config file

## Thanks to...

[go-fiber/recipes/auth-jwt](https://github.com/gofiber/recipes/tree/master/auth-jwt)