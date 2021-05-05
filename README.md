# Backend for a mini selfmade headless CMS using [Fiber](https://github.com/gofiber/fiber)

-------------------------

## Content

- [API](#api)
- [Workflows](#workflows)
   - [Create content and content types](#create-content-and-content-types)
   - [Update content entries](#update-content-entries)
   - [Create users](#create-users)
   - [Update users](#update-users)
   - [Query users and content entries by route parameters](#query-users-and-content-entries-by-route-parameters)
- [Database setup](#database-setup)
- [TO DO](#to-do)
- [Thanks to...](#thanks-to...)

-------------------------

## API

| Endpoint                 | Method    | Authentification required   | Response Fields<sup>*</sup>  | Description  |
| :----------------------- | :-------- | :-------------------------- | :--------------------------- | :----------- |
| `/api/sample`            | `GET`     | &cross;                     | `sample`                     | Health Check |
|                          | `POST`    | &cross;                     | `sample`                     | Create a new Sample Entry in the `samples` collection. Specify two string fields in the request body.   |
| `/api/auth/login`        | `POST`    | &cross;                     | `token`, `user`              | Sign in with username or email (`identity`) and `password`. On success returns token and user. |
| `/api/user/`             | `GET`     | &check;                     | `user`                       | Return all users present in the `users` collection.  |
|                          | `POST`    | &cross;                     | `token`, `user`              | Create a new user.</br> Specify the following attributes in the request body: `username`, `email`, `password`, `names`. On success returns token and user.  |
| `/api/user/*`            | `GET`     | &check;                     | `user`                       | Return users filtered by parameters in URL mathing the following regular expression: `[a-z]+=[a-zA-Z0-9\%]+`</br> The first group represents the search key and the second the search value.  |
|                          | `PATCH`   | &check;                     | `result`                     | Update user with id `id`. </br> If you want to update `role`, you have to be authenticated with a admin-user.  |
|                          | `DELETE`  | &check;                     | `result`                     | Delete user with id `id`.</br> Specify user´s password in the request body.   |
| `/api/contenttypes`      | `GET`     | &cross;                     | `contenttype`                | Return all content types present in the `contenttypes` collection. |
|                          | `POST`    | &check; (admin)             | `contenttype`                | Create a new content type.</br> Specify the following attributes in the request body: `typename`, `collection`, `field_schema`.   |
| `/api/contenttypes/:id`  | `GET`     | &cross;                     | `contenttype`                | Return content type with id `:id`.   |
|                          | `DELETE`  | &check; (admin)             | `result`                     | Delete content type with id `:id`.   |
| `/api/:content`          | `GET`     | &cross;                     | `content`                    | Return all content entries of the content type, where `content` is the corresponding collection. By convention this should be plural of the `typename`.</br> For the previous example: `content` has to be set to `events`.   |
|                          | `POST`    | &check; (admin)             | `content`                    | Create a new content entry of the content type, where `content` is the corresponding collection.</br> Specify the following attributes in the request body: `title` (string), `published`(bool), `fields`(key-value pairs: field name - field value). |
| `/api/:content/*`        | `GET`     | &cross;                     | `content`                    | Return content entries filtered by parameters in URL mathing the following regular expression: `[a-z\_]+=[a-zA-Z0-9\%]+`</br> The first group represents the search key and the second the search value.  |
|                          | `PATCH`   | &check; (admin)             | `result`                     | Update content entry with id `id` of the content type, where `content` is the corresponding collection.  |
|                          | `DELETE`  | &check; (admin)             | `result`                     | Delete content entry with id `id` of the content type, where `content` is the corresponding collection.   |

<sup>*</sup> `status` and `message` are returned on every request.

## Workflows
### Create content and content types
The content types *event* and *blogpost* are preset and you can start adding entries on those routes (`/api/events` or `/api/blogposts`). Events have custom fields *description* and *date* whereas blogposts come with *description* and *text*. By convention the collection should be plural of the typename.
If you want to create a custom content type, first use the `/api/contenttypes` endpoint, because the `/api/:content` route is validated by a lookup in the `contenttypes` collection. The mongoDB collections for new types are created automatically on first content insertion.</br>
The last attribute for a new content type, *field_schema*, is a list of key-value pairs specifying name and type of fields, that an content entry of this content type should have.</br>
Exapmle JSON request body:
```
{
    "typename": "protected-admin-test",
    "collection": "protected-admin-test-entries",
    "field_schema": {
        "text_field": "string"
    }
}
```

The last attribute for a new content entry, *fields* is a list of key-value pairs specifying name and value of fields, that should match the *field_schema* of the corresponding content type. **A validation is not yet implemented.** </br>
Example JSON request body:
```
{
    "title": "Blogpost Test",
    "published": false,
    "tags": [
        "foo",
        "bar"
    ],
    "fields": {
        "description": "Hello world",
        "date": "2021-04-08T12:00:00+02:00"
    }
}
```

### Update content entries
Every `PATCH` request updates the field `published`, so it has to be set to `true` in any request if this state is wanted after. This is due to poor handling of boolean values when using bson-flag `omitempty` in structs as update schema: *false* is interpreted as *not updated*. Therefore this flag is not used for this field and it can't be omitted or the omitted field is automatically set to false.</br>
To update custom fields you have to specify it as nested object in the request body.</br>
Example JSON request body:
```
{ "fields": {"description": "foo bar"} }
```

Preset fields can be reached directly. Example JSON request body:
```
{
    "tags": ["foo", "bar"],
    "published": true
}
```


### Create users
The admin user *adminUser* is preset with the password `ADMIN_PASSWORD` from the [.env](https://github.com/D-Bald/fiber-backend/blob/master/.env.sample) file in the root direcory of the executable.
Anybody can create a new user. The role is automatically set to *user*.</br>
Example JSON request body:
```
{
    "username": "TestUser",
    "email": "unique@mail.com",
    "password":"123",
    "names":"names",
}
```

### Update users
A user can edit the own data i.e. *username*, *email*, *password*, *names*.
Every user with role *admin* can edit any other user and particularly can edit the field *role* of any user. Roles must be updated as array containing all roles as single strings.</br>
Example JSON request body:
```
{ "roles": ["user", "admin"] }
```

### Query users and content entries by route parameters
A search parameter has the structure `key=value`. Multiple parameters are seperated by `&` (example 1). Custom fields of content entries can be queried directly so **don't** use dot-notation or similar (example 2). Only the whole field value is matched, so submatches are not supported. Queries for single user roles or single Tags are possible (example 3). To query multiple tags or roles add a new parameter for each (example 4).</br>
Examples:
   1. `/api/events/title=Title&id=606886f352caea1f9aa86471`
   2. `/api/blogposts/text=Hello%20World`
   3. `/api/users/roles=admin`
   4. `/api/blogposts/tags=foo&tags=bar`

## Database setup

In this repo the mongoDB is not self-hostet, so URI for Atlas ist hardcoded except DB_USER, DB_NAME and credentials which are read from an .env file in the root directory of the executable.

NEW: `HOSTED` variable in .env file to indicate either `ATLAS` or `DOCKER` as database host-setup

for putting up the docker database run the [*docker-compose.yaml*](https://github.com/D-Bald/fiber-backend/blob/master/docker-compose.yaml)

For self-hosted DB adjust [mongoURI in this line](https://github.com/D-Bald/fiber-backend/blob/a919ea78383a1d3fa7c30cd6498f68f72fb88620/database/connect.go#L16)


## TO DO

- Containerize backend and mongoDB
- Implement file upload
- Validate field_schema on content entry creation (https://docs.mongodb.com/manual/core/schema-validation/)

## Thanks to...

[go-fiber/recipes/auth-jwt](https://github.com/gofiber/recipes/tree/master/auth-jwt)