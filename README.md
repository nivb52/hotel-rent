# hotel-rent : 
 Rooms & Hotels managment

## ENVIROMENTS
see the env.example to see the options.
and set values in the ```.env``` file
you also can set values in the ```.env.local``` - for your local env
and ```env.test.local``` - for local test env

---
## Dev mode (hot reload using [air](#Air))
```makefile
make dev
```
- tested on WIN, 
if hot reload not working for you try to change the .air file, 
there is an air file for linux as well, 
but I didn't tested it yet.

## Start / Run
```makefile
make start
```
-- see the makefile for more commands & info

---
## API Design

I use the convention of plural for the resources in the API - 
meaning - 
```
GET  /hotels          <---> hotels 
GET  /hotels/1        <---> hotels[1]
GET  /hotels/1/rooms  <---> hotels[1].rooms

PUT  /hotels/1        <---> hotels[1] = data
DELETE  /hotels/1        <---> delete hotels[1]

POST /hotels          <---> hotels.push(data)
POST /hotels/1/rooms  <---> hotels[1].rooms.push(data) 
```
Using what I know, and as well I refresh my knowladge with [stackoverflow answer](https://stackoverflow.com/a/21809963/15039733) 
on this subject

---
## Resources
### Server 
 - Fiber v2 : [Docs](https://docs.gofiber.io/api/ctx)

### MongoDB
  - Driver : [Docs](https://mongodb.com/docs/drivers/go/current/quick-start)
  - Client Installation: 
     ```bash
      go get go.mongodb.org/mongo-driver/mongo
      ```
  - [Work with BSON](https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/)
  - [Lookup aggregation queries](https://www.mongodb.com/docs/manual/reference/operator/aggregation/lookup/#use--lookup-with-an-array)
### Authentication
 - [JWT NewWithClaims](https://pkg.go.dev/github.com/golang-jwt/jwt/v5#NewWithClaims)

### Tests
After a bit of dilema if to connect to a ‘real’ database or mock it, 
and it will be actually integration tests and not unit tests, In some way - 
```
Don’t try to mock SQL. you’ll only end up unhappy
— dylanb (user on the Gopher slack)
```
[more about that](https://markphelps.me/posts/writing-tests-for-your-database-code-in-go/) 
[reddit link](https://www.reddit.com/r/golang/comments/u62emg/mocking_database_or_use_a_test_database/?rdt=52641)

```makefile
make test
```


### Requests Mocking
In the folder ```mock/requests/```
There are several .rest files which come handy with the VS-Code plugin [humao.rest-client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)
In this project files are part files, which cannot stand alone in most cases.
They are combined with a script (written for powershell) and can be executed with make file -
```makefile 
make rest
```

This will combine all the files and then any dependency with data,
can be overcome by using the parameter spesific rest request before the wanted request
for example - 
if the request depends on some user id, liek so - 
```js
@USER_ID = {{GetAllUsers.response.body.$.[0].id}}
```
You should first invoke the dependency - ```GetAllUsers```
Which is also clickable, and will move the page to the GetAllUsers request location.

### Updating Steps

```git
git tag -a v1 -m "Step 1"
git push origin --tags
```
(I use small v not as in the realese notes)
v1 = Lesson 1 = Step 1 

---
## Dev Tools
### Air
   for Live Reload 
 - [Docs](https://github.com/cosmtrek/air)

### More
 - [Date Formating](https://zetcode.com/golang/datetime-format/)
 - [Go for NodeJs Developers](https://github.com/miguelmota/golang-for-nodejs-developers)