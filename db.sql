DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS appeal;
DROP TABLE IF EXISTS service;

CREATE TABLE IF NOT EXISTS users (
 id             bigserial   primary key,
 email          text        not null unique,
 name           text        not null,
 surname        text        not null,
 password       text        not null
);

CREATE TABLE IF NOT EXISTS  appeal (
 id             bigserial   primary key,
 email          text        not null,
 phone_number   text        not null,
 nickname       text        not null,
 subject        text        ,
 message        text        not null,
 document       text
);

CREATE TABLE IF NOT EXISTS  service (
 id             bigserial   primary key,
 price          text        not null,
 description    text        not null,
 name_service   text        not null
);


