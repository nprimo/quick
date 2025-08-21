DB_PATH ?= ./dev.db

dev:
	air

# - make change me come from input arg? (nice to have)
db-create:
	goose sqlite3 -dir=./db/migrations/ $(DB_PATH) create changeme sql

db-up:
	goose sqlite3 -dir=./db/migrations/ $(DB_PATH) up

db-down:
	goose sqlite3 -dir=./db/migrations/ $(DB_PATH) down

db-seed:
	cat ./db/seed/*sql | sqlite3 $(DB_PATH)
