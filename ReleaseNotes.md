# Release Notes
### v16: (36)
Test Booking (step 3 - last step)
- FEAT: Test Helper -  Auto Data Compare 
- FEAT: Test GetBookingsById
- FEAT: Test AdminGetBookings
- FIX: Is Admin func


### v15.5: (35)
Test Fixtures (step 2)
- FEAT: DB connection 
- FEAT: Add Booknig Fixtures
- FIX: Auth test (+ renames)


### v15: (35)
Creating Test Fixtures (step 1)
- FEAT: Extract db functions from seed script to fixtures (to use in test)
- FEAT: Extract mmocking functions from seed script to mock package (to use in test)
- FEAT: add postman collection for testing (some function are easier to test this way)
- FIX: room filter build - typo mistake (which affect logic)


### v14: (34)
Cancel booking
- FEAT: Cancel Booking
- FIX: Booking room when booking collection is missing 
    (commit: db/fix IsRoomAvailable room function)


### v13:
Room booking
- API+DB+Types:
    - feat: booking system (store, handler, types)
- Scripts & Tests:
    - new: aggregate all *.rest files into main.rest (single) file 
    - new: booking endpoint tester
        * note endpoint mock requests using VSCODE Plugin [humao.rest-client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)
    - ref: use context with timeout in test db

### <strike> v12: </strike>Skipped

### v11:
Test JWT Auth
- Test: 
    - Authenticate handler
    - ref: use context with timeout in test db
- API:
    - ref: use invalid credentials response 

### v10: (29)
Complete implement JWT Auth
- API: 
    - Protect users route by JWT Auth
    - Create custom token with 4 hour expiration time 
- Scripts: Improve the VSCode api checking using varibles
    
### v9:
- API: 
    - Get user by email
    - Start implement Auth    
- Scripts: add new data + seed users
    

### v8: (27)
- API: 
    - Get Rooms By Ids
    - Get Hotel Rooms
- Types: Room - add size as string & BedType
- Scripts: add new room fields in seed 
- Global: REFCTORING 
    - encapsulating stores, 
    - func names,
    - remove unused

### v7:
- API: Get Hotel/s
- Types: Hotel - add rating

### v6: (26½)
Seeding hotel and rooms - 
Using Insert Many for rooms and update the corresponding hotel,
the room store interface includes the hotel store interface.

- FEAT: Insert Many for rooms

### v5.9:
Seed Rooms + Imp. room types  
Update hotel rooms from the hotel store
 
### v5.5:
Add Hotel & Rooms Types + scripts - Seeding hotels
partly implement rooms
- FIXES:
 - types/fix: ID type should be primitive.ObjectID
 - api/fix: create user should return valid ID 

### v5: (24)
Adding the first test for the User API

### v4: (23)
Users complte CRUD with validation functions

### v3:
Add get users, create user, validation

### v2: 
Integrate with MongoDB, init user handler

### v1
Start