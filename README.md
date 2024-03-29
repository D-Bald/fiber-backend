# Backend for a mini selfmade headless CMS using [Fiber](https://github.com/gofiber/fiber)

## Content

- [Usage](#usage)
- [API](#api)
- [Workflows](#workflows)
    - [Roles](#roles)
    - [Create content and content types](#create-content-and-content-types)
    - [Update content entries](#update-content-entries)
    - [Create users](#create-users)
    - [Update users](#update-users)
    - [Query users and content entries by route parameters](#query-users-and-content-entries-by-route-parameters)
- [TODO](#to-do)
- [Thanks to...](#thanks-to...)

-------------------------

## Usage

You can run this package on its own by setting the [.env](https://github.com/D-Bald/fiber-backend/blob/master/.env.sample) accordingly, for exapmle with a Atlas hosted MongoDB Cluster (the *.env* file variables are used by both the docker-compose and the fiber-backend), but using the [docker-compose.yaml](https://github.com/D-Bald/fiber-backend/blob/master/docker-compose.yaml) is the easiest way to deploy all dependencies on a server.

Follow these steps:
1. Install [docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/)
2. Download [.env file](https://github.com/D-Bald/fiber-backend/blob/master/.env.sample)
    ```shell
    $ sudo wget -O .env https://raw.githubusercontent.com/D-Bald/fiber-backend/main/.env.sample
    ```
3. Set the `DB_HOST` variable in the *.env* file to the name of the docker service (in this [docker-compose.yaml](https://github.com/D-Bald/fiber-backend/blob/master/docker-compose.yaml) the service is named `mongodb`). If you use a Atlas hosted MongoDB database, set this variable to `ATLAS`. Also check environment variables like ports, database name, user and passwor and PLEASE change `SECRET` and `ADMIN_PASSWORD`.
4. Download [docker-compose.yaml file](https://github.com/D-Bald/fiber-backend/blob/master/docker-compose.yaml)
    ```shell
    $ sudo wget -O docker-compose.yaml https://raw.githubusercontent.com/D-Bald/fiber-backend/main/docker-compose.yaml
    ```
5. Execute the following commands in the root directory of [docker-compose.yaml](https://github.com/D-Bald/fiber-backend/blob/master/docker-compose.yaml):
    - To get the containers up and running execute:
        ```shell
        $ docker-compose up -d
        ```
    - To stop the containers execute:
        ```shell
        $ docker-compose down -v
        ```

This setup will create and start three docker containers:
 - [mongo](https://hub.docker.com/_/mongo/)
 - [mongo-express](https://hub.docker.com/_/mongo-express)
 - [fiber-backend](https://github.com/D-Bald/fiber-backend)

The data is persistent over multiple `up` and `down` cycles using [docker volumes](https://docs.docker.com/compose/#preserve-volume-data-when-containers-are-created).<br>
Check the database setup with [mongo-express](https://hub.docker.com/_/mongo-express) on `http://localhost:8081`.

## API

| Endpoint                 | Method    | Authentification required                     | Response Fields<sup>*</sup>  | Description  |
| :----------------------- | :-------: | :-------------------------------------------- | :--------------------------: | :----------- |
| `/api`                   | `GET`     | &cross;                                       |                              | Health-Check |
| `/api/auth/login`        | `POST`    | &cross;                                       | `token`, `user`              | Sign in with username or email (`identity`) and `password`. On success returns token and user. |
| `/api/role`              | `GET`     | &check;                                       | `role`                       | Returns all existing roles. |
|                          | `POST`    | &check; (admin)                               | `role`                       | Creates a new Role. |
| `/api/role/:id`          | `PATCH`   | &check; (admin)                               | `result`                     | Updates role with id `id`. |
|                          | `DELETE`  | &check; (admin)                               | `result`                     | Deletes role with id `id`. Also removes references to this role in user and content type documents.  |
| `/api/user`              | `GET`     | &check;                                       | `user`                       | Return users present in the `users` collection. |
|                          | `POST`    | &cross;                                       | `token`, `user`              | Creates a new user.<br> Specify the following attributes in the request body: `username`, `email`, `password`, `names`. On success returns token and user. |
| `/api/user/:id`          | `PATCH`   | &check;                                       | `result`                     | Updates user with id `id`. <br> If you want to update `role`, you have to be authenticated with a admin-user. |
|                          | `DELETE`  | &check;                                       | `result`                     | Deletes user with id `id`.<br> Specify user´s password in the request body. |
| `/api/contenttypes`      | `GET`     | &cross;                                       | `contenttype`                | Returns all content types present in the `contenttypes` collection. |
|                          | `POST`    | &check; (admin)                               | `contenttype`                | Creates a new content type.<br> Specify the following attributes in the request body: `typename`, `collection`, `field_schema`. |
| `/api/contenttypes/:id`  | `GET`     | &cross;                                       | `contenttype`                | Returns content type with id `:id`. |
|                          | `PATCH`   | &check; (admin)                               | `result`                     | Updates content type with id `:id`. |
|                          | `DELETE`  | &check; (admin)                               | `result`                     | Deletes content type with id `:id`. **Watch out: Also deletes all content entries with this content type.** |
| `/api/:content`          | `GET`     | &cross;                                       | `content`                    | Returns content entries of the content type, where `content` is the corresponding collection. By convention this should be plural of the `typename`.<br> For the previous example: `content` has to be set to `events`. |
|                          | `POST`    | &check; (depends on content type permissions) | `content`                    | Creates a new content entry of the content type, where `content` is the corresponding collection.<br> Specify the following attributes in the request body: `title` (string), `published`(bool), `fields`(key-value pairs: field name - field value). |
| `/api/:content/:id`      | `PATCH`   | &check; (depends on content type permissions) | `result`                     | Updates content entry with id `id` of the content type, where `content` is the corresponding collection. |
|                          | `DELETE`  | &check; (depends on content type permissions) | `result`                     | Deletes content entry with id `id` of the content type, where `content` is the corresponding collection. |

<sup>*</sup> `status` and `message` are returned on every request.

## Workflows
### Roles

Two roles are initiated out of the box:
```json
{
    "tag":"default",
    "name":"User"
}
{
    "tag":"admin",
    "name":"Administrator"
}
```
The *default* role is given any new user. The *admin* role is used as general access role and on start a new *adminUser* is created, if no other user with role tag *admin* is found. By changing the `name` you can decide how an admin is called and which default role is given any new user.<br>
**Warning:** Removing these roles or changing the tag causes trouble because user creation will fail due to missing default role and you can loose your last admin access user. On the next start a new *admin* role and *adminUser* is created, but this leads to a redundant *adminUser* and you can not reliably login with real a admin access. This issue can be solved by deleting the *adminUser* that has not the *admin* role, but it can be hard to debug.

Just one `GET` endpoint exists, which returns all roles. There is no use for the data of a singe role.<br>
Example JSON request body:
```json
// POST or PATCH
{
    "tag":"moderator",
    "Name":"Moderator"
}
```

### Create content and content types

The content types *event* and *blogpost* are preset and you can start adding entries on those routes (`/api/events` or `/api/blogposts`). Events have custom fields *description* and *date* whereas blogposts come with *description* and *text*. By convention the collection should be plural of the typename.
If you want to create a custom content type, first use the `/api/contenttypes` endpoint, because the `/api/:content` route is validated by a lookup in the `contenttypes` collection. The mongoDB collections for new types are created automatically on first content insertion.

The `GET` Endpoint is not protected by any middleware. `POST`, `PATCH` and `DELETE` endoints for any content are protected and you have to specify the roles that users have to have to perform each method (see example below) in the content types `Permissions` object. Users with *admin* role tag can perform any method on any content. Both default contenttypes (*event* and *blogpost*) set all method permissions to the [*default* role](#roles).

The last attribute for a new **content type**, *field_schema*, is a list of key-value pairs specifying name and type of fields, that an content entry of this content type should have.<br>
Exapmle JSON request body:
```json
{
    "typename": "protected-admin-test",
    "collection": "protected-admin-test-entries",
    "permissions": {
        "POST": [
            "Moderator"
        ],
        "PATCH": [
            "Moderator"
        ],
        "DELETE": [
            "Moderator"
        ]
    },
    "field_schema": {
        "text_field": "string"
    }
}
```

The last attribute for a new **content entry**, *fields* is a list of key-value pairs specifying name and value of fields, that should match the *field_schema* of the corresponding content type. **A schema validation is not yet implemented.** <br>
Example JSON request body:
```json
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

### Update content and content types

To update an array, like tags of content entries or permissions of content types, the whole array has to be sent. On this structure the update behaves more like a `PUT` method.

To update custom fields of content entries you have to specify it as an object in the request body.<br>
Example JSON request body:
```json
{ "fields": {"description": "foo bar"} }
```

Preset fields can be reached directly. Example JSON request body:
```json
{
    "tags": ["foo", "bar"],
    "published": true
}
```


### Create users

The admin user *adminUser* is preset with the password `ADMIN_PASSWORD` from the [.env](https://github.com/D-Bald/fiber-backend/blob/master/.env.sample) file in the root direcory of the executable.
Anybody can create a new user. The role is automatically set to *user*.<br>
Example JSON request body:
```json
{
    "username": "TestUser",
    "email": "unique@mail.com",
    "password":"123",
    "names":"names",
}
```

### Update users

A user can edit the own data i.e. *username*, *email*, *password*, *names*.
Every user with role *admin* can edit any other user and particularly can edit the field *role* of any user. Roles must be updated as array containing all roles as single strings.<br>
Example JSON request body:
```json
{
    "username": "John Doe",
    "roles": ["User","Moderator"]
}
```

### Query users and content entries by route parameters

To get all Users or all content entries of one content type, just use the bare API `GET` endpoint (example 1). To search for Users and content entries with certain properties, a query string can be added to the API endpoint. The query string begins with `?`. Each search parameter has the structure `key=value` and is case sensitive. Each document has a unique ID, that can be used to query for a single result (example 2). Multiple parameters are seperated by `&` (example 3). Custom fields of content entries can be queried directly so **don't** use dot-notation or similar (example 4). Only the whole field value is matched, so submatches are not supported. Queries for single elements of array fields like 'tags' are possible (example 5). In queries with multiple array elements, the query values currently have to have the same order as in the database and can not be a subset of the stored ones(example 6). Therefore queries with single values are recommended.<br>
Examples:
```markdown
# 1
/api/events
# 2
/api/user?_id=609273e9f17aa49bcd126418
# 3.
/api/blogposts?title=Title&published=false
# 4.
/api/events?place=home
# 5.
/api/blogposts?tags=foo
# 6.
/api/events?tags=foo,bar
```
Example 6 currently only returns documents with a full match on tags like:
```json
{ "tags": ["foo", "bar"] }
```

## TODO
- Fix Issue: Dates cannot be queried, because the `+` sign in a query string is treated as empty space. Maybe by escaping + if the parameter value is send in `""`
- Add idiomatic Endpoints for common getters and setters like: Set title, set username set names, set password...
- Implement permission control over user, roles and content types endpoints
- Issue: *standard_init_linux.go:219: exec user process caused: no such file or directory* on `docker-compose up` when using the :latest image created by workflow [CI](https://github.com/D-Bald/fiber-backend/blob/main/.github/workflows/dockerhub.yml) on GitHub Actions => workflow currently disabled and a locally on an ubuntu server built image is used in the [docker-compose.yaml](https://github.com/D-Bald/fiber-backend/blob/ee64c31317c3ccdd0b75b9ed90117d2b09207efe/docker-compose.yaml#L50).
- Implement file upload
- Validate field_schema on content entry creation (https://docs.mongodb.com/manual/core/schema-validation/)

## Thanks to...

- [go-fiber/recipes/auth-jwt](https://github.com/gofiber/recipes/tree/master/auth-jwt)
- [bmc blogs: How To Run MongoDB as a Docker Container](https://www.bmc.com/blogs/mongodb-docker-container/)
- [Docker: Base image when deploying a GoLang binary in a container (Fabian Lee)](https://fabianlee.org/2018/05/10/docker-packaging-a-golang-version-of-hello-world-in-a-container/)