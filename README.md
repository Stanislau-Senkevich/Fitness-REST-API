# Fitness REST API

---------------

## Introduction

This API provides an interface for sphere of fitness.
It can help users and trainers write down their 
workouts and manage partnerships between each other. 

To see all the endpoints follow the link: http://143.198.157.158:8001/swagger/index.html

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

