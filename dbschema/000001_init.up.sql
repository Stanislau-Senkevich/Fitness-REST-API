CREATE TABLE admins (
    id serial NOT NULL PRIMARY KEY,
    login varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL
);

CREATE TABLE users (
    id serial NOT NULL PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    role varchar(255) NOT NULL DEFAULT 'user',
    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);


CREATE TABLE partnerships (
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id),
    trainer_id int NOT NULL REFERENCES users(id),
    status varchar(255),
    created_at timestamp NOT NULL DEFAULT NOW(),
    ended_at timestamp
);

CREATE TABLE workouts (
    id serial NOT NULL PRIMARY KEY,
    title varchar(255) NOT NULL,
    user_id int NOT NULL REFERENCES users (id),
    trainer_id int REFERENCES users(id),
    description varchar(255),
    date timestamp NOT NULL DEFAULT NOW()
);