# Backend for a selfmade CMS using [Fiber](https://github.com/gofiber/fiber)

-------------------------
- [Inspired by](#inspired-by)
- [Available Routes](#available-routes)
-------------------------

## Inspired by
- [go-fiber/recipes auth-jwt](https://github.com/gofiber/recipes/tree/master/auth-jwt)
- [Quick Start: Golang & MongoDB - Modeling Documents with Go Data Structures](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures)

## Available routes

| Route           |  Method   | Description |
| :-------------- | :-------- | :------------------------------------------------------------------------------- |
| `/api/sample`   | `GET`     | Health Check
|
|               | `POST`    | Creates a new Sample Entry in the `samples` collection. Specify to String Fields in the body |
| `/api/auth/`    | `POST`    | Sign in with username or email and password. If it's successful, then generates a token |
| `/api/user/`    | `GET`     | Returns all users present in the collection `users` collection


## Usage

FILL THIS OUT, WHEN CONFIGURATION VIA CONFIG FILE IS AVAILABLE

HERE: MONGODB NOT SELF-HOSTET=> URI FOR ATLAS IS HARDCODED EXCEPT USER, DB NAME AND CREDENTIALS

## To Do:
 - Content Handler allgemein machen:
  - model 'ContentType', um eine Collection zu machen, wo dann wiederum Modelle für verschiedene Contenttypen enthalten sind:
   - Felder: ID, Type, Collection, Status (sowas wie veröffentlicht/unveröffentlicht) []Fields (Beim Ausgeben dann ein eigenes output struct, bei dem zumindest Collection nicht drin ist (wird nicht benötigt, da bei GET Frage mithilfe des Types in der 'ContentType' Collection nachgeschaut werden ka)nn, welche Collection dafür vorgesehen ist: Felder können auch geschachtelt werden) Initialisierung für "Blogposts" und "Events" (Funktion newContentType(typeName, collection, []fields) benutzen?)


- User Handler:
 - Update not only "names" but also username, password... 