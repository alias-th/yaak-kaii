DB_URL=postgresql://root:1234@localhost:5432/yaak_kaii?sslmode=disable
MIGRATION_DIR=internal/db/migration

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down 1

test:
	go test -v -cover -short ./...

new_migration:
	migrate create -ext sql -dir "$(MIGRATION_DIR)" -seq $(name)

sqlc:
	sqlc generate

.PHONY: test migrateup migrateup1 migratedown migratedown1 new_migration sqlc