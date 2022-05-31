LOCAL_DB_DSN:="postgres://postgres:postgres@localhost:5432/demo?sslmode=disable"
POSTGRES_DB_DSN:="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
LOCAL_DB_NAME:=demo

db:
	psql -d $(POSTGRES_DB_DSN) -c "drop database if exists \"$(LOCAL_DB_NAME)\""
	psql -d $(POSTGRES_DB_DSN) -c "create database \"$(LOCAL_DB_NAME)\""
	goose -dir migrations postgres $(LOCAL_DB_DSN) up

build:
	go build -o ./expr_demo ./cmd/expr_demo/main.go

