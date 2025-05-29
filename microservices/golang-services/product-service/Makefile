ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

migrate-create:
	@ migrate create -ext sql -dir migration -seq $(name)

migrate-up:
	@ migrate -database ${DATABASE_URL_POSTGRES} -path migration up

migrate-down:
	@ migrate -database ${DATABASE_URL_POSTGRES} -path migration down
