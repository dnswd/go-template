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

migrate:
    bash ./infra/database/migration.sh

clean-services:
    rm -rf ./.local/

