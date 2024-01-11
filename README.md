# Backend for Bakasov Marat's porfolio website

Uses [generator](https://github.com/Marattttt/Generator) Go package for creating simple drawings

Collects statistics about visitors (No personal data!)

## Usage

### Docker

Docker deployment is to be added

### Manual

All listed commands are run from portfolio/back

1. (optionally) Edit the dbinit.sql file to change the created role's name and give it a password
2. Initialize the database with 
```shell
  psql < dbinit.sql
```
3. Edit the .env_example file to set the environment variables required to run the app
4. Install Go if it isn't already installed from [link](https://go.dev/doc/install)
5. Compile the project
```shell
  go build back/
```
6. Export the variables from .env_example
```shell
  export $(cat .env_example | xargs)
```
7. Run the compiled binary
```shell
  ./portfolio_back
```

## Plan

- [x] Database init
- [ ] Endpoints and openapi
- [ ] CRUD for guests
- [ ] Authentication
