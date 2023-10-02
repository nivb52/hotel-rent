# Release Notes

### v9: (28)
API: 
    - Get user by email
    - Start implement Auth    
Scripts: add new data + seed users
    

### v8: 
API: 
    - Get Rooms By Ids
    - Get Hotel Rooms
Types: Room - add size as string & BedType
Scripts: add new room fields in seed 
REFCTORING: 
    - encapsulating stores, 
    - func names,
    - remove unused

### v7: 
API: Get Hotel/s
Types: Hotel - add rating

### v6:
Seeding hotel and rooms - 
Using Insert Many for rooms and update the corresponding hotel,
the room store interface includes the hotel store interface.

FEAT: 
- Insert Many for rooms

### v5.9:
Seed Rooms + Imp. room types  
Update hotel rooms from the hotel store
 
### v5.5:
Add Hotel & Rooms Types + scripts - Seeding hotels
partly implement rooms
- FIXES:
 - types/fix: ID type should be primitive.ObjectID
 - api/fix: create user should return valid ID 

### v5:
Adding the first test for the User API

### v4:
Users complte CRUD with validation functions

### v3:
Add get users, create user, validation

### v2: 
Integrate with MongoDB, init user handler

### v1 
Start