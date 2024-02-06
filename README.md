# Fitness REST API

---------------

## Introduction

This API provides an interface for sphere of fitness.
It can help users and trainers write down their 
workouts and manage partnerships between each other. 

To see all the endpoints follow the link: http://droplet.senkevichdev.work:8001/swagger/index.html

### Models 
- Admin
- Trainer
- User(Client)

### Features
#### Admin
- Create, update, delete users
- Get full information (including workouts and partnerships) about users

#### Trainer
- Initialize/end partnership with users
- Accept/deny requests for partnership from users 
- Create, update, delete workouts with his clients
- Get information about his clients and workouts with them

#### User (Client)
- Get information about account, its partnerships and workouts
- Get information about trainers
- Send request for partnership to trainer
- Create, update, delete workouts with trainer with whom partnership was established
- Create, update, delete workout without trainer
------------------
## Technologies
- #### Go 1.18
- #### Gin
- #### PostgreSQL
- #### Docker
- #### Swagger
- #### CI/CD (GitHub Actions)

-----------------
## Realization features
- #### Clean architecture
- #### Authorization with JWT tokens
- #### Unit tests for repository and handlers
- #### Linter
- #### DB Migrations 
- #### JSON logging (logrus)

-----------------

## Tools and libraries

### Database (db)

- `jmoiron/sqlx`: Database interactions.
- `lib/pq`: PostgreSQL driver.

### API (api)

- `gin-gonic/gin`: Web framework for Go.

### Authentication (auth)

- `dgrijalva/jwt-go`: JWT functionality.

### Configuration (config)

- `spf13/viper`: Go configuration solution.

### Testing (tests)

- `golang/mock`: Mocking functionality.
- `stretchr/testify`: Assertion functions.
- `zhashkevych/go-sqlxmock`: Mocking functionality for SQLX.

### Documentation (docs)

- `swaggo/files`: Serving Swagger files.
- `swaggo/gin-swagger`: Gin middleware for API documentation.
- `swaggo/swag`: Go documentation tool.

