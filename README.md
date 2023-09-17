# hotel-rent : 
 Rooms & Hotels managment

## Run
try use the makefile

## Resources
### server - fiber
 - [Docs](https://docs.gofiber.io/api/ctx)

### MongoDB
  - Driver : [Docs](https://mongodb.com/docs/drivers/go/current/quick-start)
  - Client: 
     ```bash
      
      go get go.mongodb.org/mongo-driver/mongo
      
      ```
### Tests
After a bit of dilema if to connect to a ‘real’ database or mock it, 
and it will be actually integration tests and not unit tests, In some way - 
```
Don’t try to mock SQL. you’ll only end up unhappy
— dylanb (user on the Gopher slack)
```
[more about that](https://markphelps.me/posts/writing-tests-for-your-database-code-in-go/) 
[reddit link](https://www.reddit.com/r/golang/comments/u62emg/mocking_database_or_use_a_test_database/?rdt=52641)


## Dev Dependencies
### LiveReload - Air 
 - [Docs](https://github.com/cosmtrek/air)

### More
 - [Go for NodeJs Developers](https://github.com/miguelmota/golang-for-nodejs-developers)