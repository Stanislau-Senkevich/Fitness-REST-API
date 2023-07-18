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

CREATE TABLE partnerships (
    id serial NOT NULL PRIMARY KEY,
    user_id int REFERENCES users (id),
    trainer_id int REFERENCES trainers (id),
    status varchar(255) NOT NULL,
    description varchar(255)
);

CREATE TABLE workouts (
   id serial NOT NULL PRIMARY KEY,
   user_id int REFERENCES users (id),
   trainer_id int REFERENCES trainers (id),
   description varchar(255),
   date timestamp NOT NULL DEFAULT NOW()
)