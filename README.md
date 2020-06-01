# Assignment

##User Story 
As a user, I want to be able to select from a list a product, and based on this product to receive a
subscription plan.
### AC:
- I can see a list of products
- I can select a product
- On the page I can have a CTA to "Purchase" the product
- On a dedicated page I want to see in a readable format the following details about my 
   subscription: Starting date, duration, price
- On the above page I have the option to pause/unpause my subscription
- On the above page I can cancel my subscription

##Technical details:
Each product from the list should have a different subscription duration and a price to reflect
the length of the subscription duration

# Running website

## Install
Ensure you have Docker and Docker Compose installed and running.

You can build and start the website with the following command:
`$ docker-compose up --build`

This will cause docker to build and run 3 containers.
- db: This container hosts the database. It will be avalible on `localhost:3306`
- backend: This container hosts the backend application. It will be avalible on `localhost:8080`
- frontend-react: This container host the frontend application. It will be avalible on `localhost:3000`


In your web browser go to url `http://localhost:3000/`. This will take you to the landing page which
lists all the products.

## Frontend 
The home landing page has a list of all the products along with a navigation bar up the top with
a link to a login page.

If you want to Buy a product/look at your orders you will need to log in. I have not built a
registration page so you will need to use an existing user.

Once logged in the nagivation bar will change and you will have access to a list of your
subscriptions. On this page you will be able to start/pause/cancel subscriptions. 
The product listing page will also change so you are allowed to buy subscriptions.

## Database 
### users
store user details
```
+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| id       | int(11)      | NO   | PRI | <null>  | auto_increment |
| name     | varchar(255) | NO   |     | <null>  |                |
| email    | varchar(255) | NO   |     | <null>  |                |
| password | varchar(255) | NO   |     | <null>  |                |
+----------+--------------+------+-----+---------+----------------+
```
### products
store product information 
```
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| id          | int(11)      | NO   | PRI | <null>  | auto_increment |
| name        | varchar(255) | NO   |     | <null>  |                |
| description | text         | NO   |     | <null>  |                |
| price       | int(11)      | NO   |     | <null>  |                |
| duration    | int(11)      | NO   |     | <null>  |                |
+-------------+--------------+------+-----+---------+----------------+
```
### subscriptions
store list of subscriptions to products by users. This table links to users and products table
```
+------------+--------------+------+-----+---------+----------------+
| Field      | Type         | Null | Key | Default | Extra          |
+------------+--------------+------+-----+---------+----------------+
| id         | int(11)      | NO   | PRI | <null>  | auto_increment |
| user_id    | int(11)      | NO   | MUL | <null>  |                |
| product_id | int(11)      | NO   | MUL | <null>  |                |
| start_date | datetime     | NO   |     | <null>  |                |
| status     | varchar(255) | NO   |     | <null>  |                |
| nextStatus | varchar(255) | NO   |     | <null>  |                |
+------------+--------------+------+-----+---------+----------------+
```
## Backend
### Database handler
Conects to the database, applies migration tasks, returns a database object.

### Scheduler
schedules tasks to run daily. I have one task that renews expired subscriptions on the
database

### Server
provides rest api for th frontend. Includes buisness logic and database operations for now.
(see improvements section) 

GET /forbidden
Returns forbidden page.

POST /login
Authenicates user, supplies them with valid session so they can access personal informaiton.

POST /logout
Destroys users session and logs them out.

GET /userInfo
Returns user information. 
Requires user to be authenticated.

GET /products
Returns a list of all avalible products.

POST /product
Returns a spesific product matching the id passed in post request.

POST /newSubscription
Signs user up for a new subscription. will not work if user already has subscribed to this product.
Requires user to be authenticated.

POST /updateSubscription
Updates subscription status based on post data. 
Requires user to be authenticated.

GET /subscriptions
List of all users subscriptions. 
Requires user to be authenticated.

POST /cronUpdate
Triggers database update task. This is included for testing only and should be removed. This task
is normally triggered by the cron scheduler.

## Manual testing

GET /forbidden
`$ curl -X GET http://localhost:8080/forbidden`

POST /login
`$ curl -d '{"email":"ann@cmail.com", "password":"passwd"}' -c cookies.txt -H "Content-Type:application/json" -X POST http://localhost:8080/login`

POST /logout
`$ curl -b cookies.txt -c cookies.txt -X POST http://localhost:8080/logout`

GET /userInfo
`$ curl -b cookies.txt -X GET http://localhost:8080/userInfo`

GET /products
`$ curl -X GET http://localhost:8080/products`

POST /product
`$ curl -d '{"ProductID":1}' -H "Content-Type:application/json" -X POST http://localhost:8080/product`

POST /newSubscription
`$ curl -d '{"productId":1}' -b cookies.txt -X POST http://localhost:8080/newSubscription`

POST /updateSubscription
`$ curl -d '{"productId":1, "status":"ACTIVE"}' -b cookies.txt -H "Content-Type: application/json" -X POST http://localhost:8080/updateSubscription`

GET /subscriptions
`$ curl -b cookies.txt -X GET http://localhost:8080/subscriptions`

POST /cronUpdate
To set up data you might want to start by making some current subscriptions and changing their
dates.

Log into data with your favourite database client and run the following sql statement:

```
USE products;
UPDATE subscriptions SET start_date=DATE_ADD(start_date, INTERVAL -15 DAY);

```
This will set all dates in the database back 15 days. 14 days is the longest subscription being
offered so this should make them candidates for expiary.

Now call the endpoint.
`$ curl -d '{"PrivateToken":"top-secret"}' -b cookies.txt -H "Content-Type: application/json" -X POST http://localhost:8080/cronUpdate`

Check the results with: `SELECT * FROM subscriptions`
All entries should have status set to "EXPIRED" or "CANCELLED" except the latest. This will be the newly created subscriptions with the start date set to todays date.

# Implementation

# Assumptions

There was some ambiguity about how the subscriptions work in the acceptance criteria.
These are some of my assumptions:

- The subscription will last until you cancel it. 

- The duration is the period until the next period.

- Buying a subscription will start your subscription immediately. 

- Cancelling a subscription will cancel it immediately.

- Starting/Pausing doesn't affect the current subscription. It puts in a request for the next
  cycle of the subscription. So you either will/will not be charge. It is about pausing the
  payments.  

## Tech stack

- Docker
- Docker-Compose

### Database:
- mySql

### Frontend:
- React
#### Libraries:
- axios - for api calls
- reach/router - for routing incomming requests
- react-bootstrap - for building smarter forms

### Backend:
- Golang
#### Libraries:
- gorilla - for building endpoints, handling sessions, managing cookies 
- go-sql-driver - for connecting to the database 
- golang-migrate - for handling database migrations
- robfig/cron - for setting up scheduling

## Improvements
### Backend
- The app file has lots of code entanglement. This requires refactoring to seperate buisness
  logic, code to map it into the correct datastructures, and the database queries.

- Idealy I would rewrite the database queries so the input is serialised into an object and
  sanitised before inserting it into the database. As it stands now it is highly likely to
  break or work as an entry point for sql injection attacks.

- The return values from my endpoints are not well defined. Sometimes returning json, 
  sometimes plain text and sometimes just a status code with no body. This should be 
  standised and documented, preferably with some auto-generated documentation based on the
  code if it is avalible.

- Endpoints have no security. I would like to use JWT tokens to secure communication between 
  react and my application, along with nounces to prevent replay attacks. Additionally I 
  should not be storing my passwords as plain text. They should be stored as hashes. 

- The Task scheduler is an excellent candiate for being ripped out into a microservice. This 
  would require making the database a library so I can reuse this code for both applications.
  But I think the seperation here would make a lot of sense. 

### Frontend
I am new to react and was learning how to use it on this project so there is potentially a lot of styliasation that I still need to learn. 

- The most notable thing missing from my frontend is CSS. I focused my efforts into building 
  the logic on these web pages and ran out of time to do this.

- Decide if I want to use the react-bootstrap framework or not. Currently this is only used in
  one form. I should make a decision to use this frame work across the application or remove 
  it completely. It doesn't make sense to add a dependancy here for one form. 

- Improve handling of status codes returned by the rest api. This will be easier once I make 
  the api more defined on the backend. 
