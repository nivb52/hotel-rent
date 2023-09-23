# hotel-rent : 
 Rooms & Hotels managment

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
---
## Dev Tools
### Air
   for Live Reload 
 - [Docs](https://github.com/cosmtrek/air)

### More
 - [Go for NodeJs Developers](https://github.com/miguelmota/golang-for-nodejs-developers)