#! /bin/zsh

mkdir ./postgres_data
initdb ./postgres_data/
pg_ctl -D ./postgres_data start
psql postgres -f init_test_postgres.sql

export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=testpass
export TEST_DB_PORT=5432
export TEST_DB_NAME=postgres
export TEST_DB_HOST=localhost
