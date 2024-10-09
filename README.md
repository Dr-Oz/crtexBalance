# Microservice for working with user balance

## Mechanisms for implementing fault tolerance of the service
- Using transactions when interacting with the database to ensure data integrity in case of errors.
- Logging of errors and operations
- Logging
- Distribution of the system. Separation of business logic, handler logic and database logic.
- Containerization
- RabbitMQ has been implemented (So far only for the balance replenishment function)

## What else can be done to prevent the loss of orders
- Implement RabbitMQ for the rest of the functions
- A monitoring system to track the status of the service and database
- Use Kuber
- It is recommended to write tests

## Instalation
- Clone the repository using the command: _``git clone https://github.com/Dr-Oz/crtexBalance.git ```_
- Go to the _crtexBalance folder_
- Launch the microservice using Docker Compose: ``docker-compose -f deploy/docker-compose.yml up -d```
***

## Parameters
- To specify the path to the configuration file, use the -config flag with the file path (used by default./configs/config.yaml).
- To perform database migration, use the -migrationup flag
- To roll back the database migration, use the -migrationdown flag

Example:
```
./crtexBalance -config ./configs/newconfig.yaml -migrationdown -migrationup
```
***

## To interact with the microservice, use the following functions:
### 1. Replenishment of the user's balance
Send POST request on ```localhost:8081/topup```  with JSON:
```json
{
    "userid":15,
    "amount":500,
    "date":"2022-08-01"
}
```
### 2. Getting user's balance
Send POST request on ```localhost:8081/```  with JSON:
```json
{
    "userid":15
}
```
If success -> get JSON like this:
```json
{
    "userid": 15,
    "balance": 500
}
```
### 3. Transfer of funds from the user to the user
Send POST request on  ```localhost:8081/transfer``` with JSON:
```json
{
    "fromuserid":1,
    "touserid":2,
    "amount":100,
    "date":"2022-10-10"
}
```
If success -> get JSON like thisN:
```json
{
    "message": "transfer funds performed"
}
```
***

## Swagger
Address ``http://localhost:8081/swagger/index.html``
