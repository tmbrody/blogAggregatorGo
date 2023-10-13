#!/bin/bash

source .env

goose -dir sql/schema postgres $CONN down