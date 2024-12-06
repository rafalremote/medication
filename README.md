# Medication API

This is a RESTful API for managing medication records.

### Reasoning behind implementaion:

For the sake of simplicity, I built just one table because I wanted to focus on my Go programming rather than designing a fully detailed database. 

If I had designed a more complete database, I would have included tables for storing information about companies producing medications, max. allowed dosages, the number of available items, serial numbers for items, and a flag to mark a medication's serial number as withdrawn from sale.

Also I would use **UUID** as index instead of integers to allow better scalability.

Instead of building a mock of DB I prefer having a real DB (**medication_test**) just for testing that gets clean up after tests. This gives options for more realistic test scenarios.

## Features

- Full CRUD functionality for medication records
- JWT-based authentication (required on **PROD** env)
- Swagger API documentation (only on **DEV** env)
- Pagination
- Seeding controlled by ENV variable (`SEED=true`)
- **DEV**, **PROD** environments are set by ENV variable (`TARGET_RELEASE=DEV` or `TARGET_RELEASE=PROD`)
- Makefile for quick setup and running

## TODO

I've encountered problems running tests inside test container, so I stopped there and I miss tests here.
I would also change import paths directly to github as it's not available there.

## Prerequisites

- Docker
- Go 1.23
- `swag` CLI tool
    - (`go install github.com/swaggo/swag/cmd/swag@latest`)
- databases created beforhand:
    - **medication**
    - **medication_test**

## Running the Application

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/medication.git
   cd medication
   ```

2. Start the application (uses Docker Compose):

    ```bash
    make start
    ```

    but for testing purpose, you can seed database by passing ENV variable **seed**:

    ```bash
    SEED=true make start
    ```

    This will seed 1000 medications to **medication** db, so you can test pagination.

3. Access the API:

- Base URL: http://localhost:8080/
- Swagger UI: http://localhost:8080/swagger/index.html

    Swagger endopoint is available **only on DEV** environment.

4. Stop the application

    ```bash
    make stop
    ```

5.  You can also stop and clean up (remove docker orphans and volumes):


    ```bash
    make clean
    ```
    

## Generating Swagger Documentation

If you make changes to the routes or Swagger annotations, regenerate the documentation:

```bash
make swagger
```

## API Endpoints

### Medications

| Method | Endpoint | Description | 
| --- | --- | --- | 
| GET | /medications | Fetch all medications |
| GET | /medications/?limit={int}&offset={int} | Fetch all medications |
| POST | /medications | Create a new medication |
| GET | /medications/{id} | Fetch a medication by ID |
| PUT | /medications/{id} | Update a medication by ID |
| DELETE | /medications/{id} | Delete a medication by ID |


### Authentication

Authentication should be enabled for production environments.

To enable it, set ENV variable: `TARGET_RELEASE=PROD`.

Once enabled, all routes are protected using JWT-based authentication. 
Provide a valid token in the Authorization header:

`Authorization: Bearer <your-token>`

By default: `JWT_SECRET=your_jwt_secret` is set as ENV.

## Testing

Run unit tests:

`make test`

It has a dedicated docker compose file.
