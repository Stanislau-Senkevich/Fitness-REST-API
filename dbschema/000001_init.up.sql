CREATE TABLE admins (
    id serial NOT NULL PRIMARY KEY,
    login varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL
);

CREATE TABLE users (
    id serial NOT NULL PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE trainers (
    id serial NOT NULL PRIMARY KEY,
    login varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE users_partnerships (
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id),
    trainer_id int NOT NULL REFERENCES trainers(id),
    status varchar(255),
    created_at timestamp NOT NULL DEFAULT NOW(),
    ended_at timestamp
);

CREATE TABLE users_workouts (
    id serial NOT NULL PRIMARY KEY,
    title varchar(255) NOT NULL,
    user_id int NOT NULL REFERENCES users (id),
    trainer_id int REFERENCES trainers(id),
    description varchar(255),
    date timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE trainers_partnerships (
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users (id),
    trainer_id int NOT NULL REFERENCES trainers (id),
    status varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    ended_at timestamp
);

CREATE TABLE trainers_workouts (
   id serial NOT NULL PRIMARY KEY,
   title varchar(255) NOT NULL,
   user_id int NOT NULL REFERENCES users(id),
   trainer_id int NOT NULL REFERENCES trainers (id),
   description varchar(255),
   date timestamp NOT NULL DEFAULT NOW()
)