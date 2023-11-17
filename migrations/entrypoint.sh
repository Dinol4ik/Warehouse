#!/bin/bash

DBSTRING="host=postgres user=user password=postgres dbname=lamoda sslmode=disable"

goose postgres "$DBSTRING" up