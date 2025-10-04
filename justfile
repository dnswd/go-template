set dotenv-load

help:
    @just --list

[working-directory: './backend/']
backend:
    air

live:
    nix run .#live

services:
    nix run .#runServices

[working-directory: './backend/']
sqlc:
    sqlc generate

migrate:
    bash ./infra/postgresql/migration.sh

clean-services:
    rm -rf ./.local/

