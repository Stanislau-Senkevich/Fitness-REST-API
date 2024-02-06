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
### Database
- **jmoiron/sqlx**: Used for database interactions in a more ergonomic way than raw SQL.

- **lib/pq**: v1.10.9: PostgreSQL driver for Go's database/sql package.

### API
- **gin-gonic/gin**: Web framework for Go. It features a Martini-like API with much better performance, up to 40 times faster.

### Authentication
 - **dgrijalva/jwt-go**: Provides a straightforward way to create and verify JWT.

### Configuration

- **spf13/viper**: Complete configuration solution for Go applications.

### Testing
- **golang/mock**: Testing package generates mock implementations of Go interfaces.

- **stretchr/testify**: Toolkit with various packages that provide many assertion functions.

- **zhashkevych/go-sqlxmock**: Provides mocking functionality for the sqlx package.

### Documentation
- **swaggo/files**: Serves swagger files from the general filesystem.

- **swaggo/gin-swagger**: Gin middleware to automatically generate RESTful API documentation with Swagger 2.0.

- **swaggo/swag**: Go documentation tool for API documentation. It extracts comments from Go source files and turns them into a formatted documentation.

