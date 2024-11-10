# TrinityApp
Design and implement a promotional campaign system for the Trinity app that offers a 30% discount on Silver subscription plans to the first 100 users who register through campaign links (note that vouchers generated for the campaign are only valid for a limited period).
### Port: 8080
### Routes: api/server.go 
## 1. Requirements and Setups
- OS: Ubuntu 
- Go 1.23
- [Go Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- Make for Makefile
- Docker
- Steps:
  - Clone this repo
  - Run Make commands with orders (make pullpostgres, make postgres, make migrate-up)
  - Start the application (go run main.go)

## 2. Database schema
![trinity](https://github.com/user-attachments/assets/c3dc3aed-3353-4693-ac74-f71c32ca243d)
[Latest diagrams](https://dbdiagram.io/d/trinity-672ecfede9daa85acad8ad46)

## 3. Flow chart
![Untitled Diagram](https://github.com/user-attachments/assets/c23f8d8c-4a1e-4c27-99ee-07b0d4a05fbe)

## 4. Improvement (Due to time limit, the app still need some improvements)
- Need implement Authenticate and Authorize
- Need implement Pagination and Search
- Need better structure
- Need document
