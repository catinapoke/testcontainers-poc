#!/bin/bash

DBHOST="localhost:5432"
DBUSER="user"
DBPASSWORD="password"
DBNAME="users"
DBSSL="disable"
DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"

goose postgres "$DBSTRING" up