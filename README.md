![go-clean-arch-tt](docs/img/purpleGopher.svg)

# go clean architecture (tahassurtalih)



wallet management service sadece crud operasyonları için kullanılsın. (kendine özel postgresqli olacak, sadece bu işlemler için. created, updated, delete işlemleri için ayrıca diğer dblerle diyalog kurulabilir ama en sona koyabilirisn.)
assetmanageemnt service (command api olsun)- event journala yazacak olan api bu olacak, withdraw, deposit vs hepsi burada olacak. 
asset - querry service (query api olsun, kendine özel postgresqli read model olacak)


go: finding module for package github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs
go: finding module for package github.com/ozlemugur/go-cqrs-event-sourcing-tt/internal/entity
go: github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/controller/http/v1 imports
        github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs: no matching versions for query "latest"
go: github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/usecase/repo imports
        github.com/ozlemugur/go-cqrs-event-sourcing-tt/internal/entity: no matching versions for query "latest"



[Error - 6:45:21 PM] Request textDocument/hover failed.
  Message: no package data for import "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs"
  Code: 0 
[Error - 6:45:21 PM] Request textDocument/hover failed.
  Message: no package data for import "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs"
  Code: 0 
[Error - 6:45:21 PM] Request textDocument/hover failed.
  Message: no package data for import "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs"
  Code: 0 
[Error - 6:45:22 PM] Request textDocument/hover failed.
  Message: no package data for import "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs"
  Code: 0 
[Error - 6:45:22 PM] Request textDocument/hover failed.
  Message: no package data for import "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs"
  Code: 0 
[Error - 6:45:23 PM] Request textDocument/hover failed.
  Message: no package data for import "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs"
  Code: 0 
[Error - 6:45:27 PM] 2025/01/29 18:45:27 command error: err: exit status 1: stderr: go: finding module for package github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs
go: finding module for package github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/docs
no matching versions for query "latest"

 swag init -g asset-management-service/internal/controller/http/v1/router.go -o asset-management-service/docs

 go mod tidy çalışmıyordu.


** congifursyonların envden alınmasını ayarla.
** KRaft mode çalışılabilir.



TODOS:

- laest swagger go install github.com/swaggo/swag/cmd/swag@latest



- Don’t reinvent the wheel.

- Polymorphism gives you the ability to create one module calling another and yet have the compile time dependency point against the flow of control instead of with the flow of control.
you have absolute control over your depedency structure you can avoid wiriting fragile rigid and non reusable modules.  (Robert C. Martin)

- the dependencies oppose the flow of control. this inversion of depedencies prevents the system from rotting because it stops the fan out of the copy module from growing. the copy module doesn't need to be modified because all of its outgoing depedencies terminate at the file abstracation ,new devices can be added ad nauseam without affecting the copy program one little bit.  (Robert C. Martin)

- high level modules should not depend upon low level modules. Both should depend upon abstractions.

- This project is developed based on the https://github.com/evrone/go-clean-template repository and was selected from among other repositories in the starred list for further development.


# What was paid attention to ?

Evrone repository rules: [https://github.com/evrone/go-clean-template?tab=readme-ov-file#the-main-principle](https://github.com/evrone/go-clean-template?tab=readme-ov-file#the-main-principle)

**The inner layer** with business logic should be clean. It should:

- Not have package imports from the outer layer.
- Use only the capabilities of the standard library.
- Make calls to the outer layer through the interface (!).

The business logic doesn't know anything about Postgres or a specific web API. Business logic has an interface for working with an _abstract_ database or _abstract_ web API.

**The outer layer** has other limitations:

- All components of this layer are unaware of each other's existence. How to call another from one tool? Not directly, only through the inner layer of business logic.
- All calls to the inner layer are made through the interface (!).
- Data is transferred in a format that is convenient for business logic (`internal/entity`).

For example, you need to access the database from HTTP (controller). Both HTTP and database are in the outer layer, which means they know nothing about each other. The communication between them is carried out through `usecase` (business logic):



## Content
- [Quick start](#quick-start)
- [Project structure](#Project-structure)
- [Local Debugging Cheat Sheet](#local-debugging-cheat-sheet)
- [Postgresql monitoring](#Postgresql-monitoring)
- [Start over from the beginning](#Start-over-from-the-beginning)
- [docker helper commands](#docker-helper-commands)
- [missing parts](#missing-parts)



## Project structure


### Main Entry Point
#### cmd/app/main.go

    •	The entry point of the application.
    •	This file defines the starting flow of the program.

### Configuration
#### config/

    •	Contains configuration files for the application.
    •	Examples: database connections, environment variables, API keys.

### Documentation
### docs
    •	Stores Swagger documentation files and other technical documentation.
    •	Use the swag init command to update Swagger documentation here.

### Application Logic
### internal/app

    •	Contains the Run function of the program.
    •	Manages the initialization of components, such as the HTTP server and database connections.

#### Controller Layer (Outer Layer)
### internal/controller/http 

    •	Handles the REST API layer.
    •	Built using the Gin Framework, defining HTTP endpoints.
    •	Responsible for REST versioning.
    •	Entry point for generating Swagger documentation (swag init should point to this folder).


#### Business Logic Layer (Inner Layer)
### internal/usecase/
    •	Contains the core business logic of the application.
    •	Called by the controller layer to execute specific use cases.


#### Repository Layer (Inner Layer)
### internal/repo/

    •	Manages all database-related operations (CRUD).
    •	Isolates database logic from the rest of the application.

#### Web API Layer (Inner Layer)
### internal/webapi/ (inner layer)

    •	Handles interactions with external APIs.
    •	Example: Sending messages using third-party API calls.

#### Entity Layer (Inner Layer)
### internal/entity/ (inner layer)

    •	Defines the entities used in the business logic.
    •	These entities are shared and accessible across all layers.

#### Reusable Packages
### pkg/ 
    •	Contains reusable packages and modules.
    •	Any package here can be imported and used by anyone who imports the module.


## Quick Start

```sh

# To bring everything up together, use the compose-up command.  Postgres, App
$ make compose-up

# To insert mock messages.
$ make mock-messages:

# check application health
$ check-health:

```

## Local Debugging Cheat Sheet

You can debug the application locally using Visual Studio Code. Here’s an example launch.json configuration:
```sh
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug API",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/app",// The main directory of the application
            "cwd": "${workspaceFolder}",
            "env": {
                "GIN_MODE": "debug" // Environment variables
            },
            "args": [] // Add any arguments required for your API to run here
        }
    ]
}
```

Before debugging, ensure you have a .env file:

```sh
# we should copy env file, it works in unix based operating systems
$ cp .env.example .env 

# prepare the environment
$ make prepare

# create swagger files
$ make swag-v1

#  Bring up other components without starting the app
$ make compose-up-without-app

# check the containers
$ docker ps

# migrate the tables
$ make migrate-up 
# if you see these messages, you should download migrate tool.
# migrate -path migrations -database 'postgres://user:pass@localhost:5432/postgres?sslmode=disable' up
# make: migrate: No such file or directory
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# retry migrate the tables
$ make migrate-up 

# Access the PostgreSQL database running in a Docker container as the specified user and database.
$ docker exec -it postgres psql -U user -d postgres

# View table structure:
$ \d messages

# if you see this message "Did not find any relations." you should onload migrate tool.
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

#you can check the status of the new tool
$ which migrate

# # Run the application locally in debug mode.

# To insert mock messages.
$ make mock-messages:

# check application health
$ check-health:

```



## Postgresql monitoring

To see what is happening in the database:

```sh
# Access the PostgreSQL database running in a Docker container as the specified user and database.
$ docker exec -it postgres psql -U user -d postgres

# List tables:
$ \dt

# View table structure:
$ \d messages

# select the data from tha table

$ select * from messages;

#  Quit psql:
$ \q 
```


## Start over from the beginning

For the fresh start you can track these steps:

```sh
# # Stop and remove all containers
$ make compose-down

# Remove Docker volumes (if you prefer to keep database changes, skip this step)
$ make docker-rm-volume

# Check and remove unused Docker images
$ docker images
$ docker rmi "imageid"

# Resolve port conflicts (e.g., port 8080 in use)
$ lsof -i :8080
# COMMAND     PID      USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
# app       42073 ozlemugur    9u  IPv6 0x9be0dbb3dbb48eb7      0t0  TCP *:http-alt (LISTEN)
# Arc\x20He 82052 ozlemugur   34u  IPv6  0xf726809424a9b2f      0t0  TCP localhost:49923->localhost:http-alt (CLOSED)
$ kill -9 42073
```


## docker helper commands


```sh
# Volume management:
$ docker volume ls  
$ docker volume ls -f dangling=true
$ docker volume rm "volumename"
$ docker volume inspect "volumename"
```

```sh
# Container management:
$ docker ps
$ docker ps -a 
$ docker rm "containerid"
$ docker start/stop "containerid"
$ docker images
$ docker rmi "imageid"
$ docker exec -it postgres psql -U user -d postgres
$ docker logs -tail=all <containerid>
```



#### mock messages :

```sh
{"content": " Diego, libre dans sa tête Song by Johnny Hallyday https://open.spotify.com/track/0qJW9XIdyvr4yQrlUFP8xq?si=34cbcd7be5b34ea7", "recipient_phone": "05057136986"}
{"content": " Legends never die Song by League of legends, Against The Current https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c", "recipient_phone": "05057136986"}
{ "content": " Canta Per me by Yoki kajiura  https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c","recipient_phone": "05057136988" }
{ "content": " king Song by Florence and the machine https://open.spotify.com/track/1VSngtLdJhrlfHkLxTyOXK?si=d9292df2504e4da0", "recipient_phone": "05057136986" }
{ "content": " mother nature Song by The Hu and LP https://open.spotify.com/track/35SoEGEXsaNnfi8PsT8xEC?si=4347af98187349fa", "recipient_phone": "05057136986"}
{"content": " winding river by yu-peng chen https://open.spotify.com/track/04WnFdVesT0VLu1Fc57VoI?si=c12404e1bbb34564","recipient_phone": "05057136986"}
{ "content": " guizhong's lullaby by jordy chandra, beside bed https://open.spotify.com/track/0n2sLg3mtyxqGZMtGf0Uow?si=0617f0ebd3b148b2", "recipient_phone": "05057136986"}
{ "content": " nana para mi song by clara peya silvia perez cruz https://open.spotify.com/track/5IY5cuo1nQbcJvzE8h2YvF?si=8049c9c2cc0b4d4c","recipient_phone": "05057136986"}
```


## missing parts

-In the FetchAndSendMessages use case, after a successful web call, if we cannot update the database, we would need to create a retry mechanism to ensure data consistency. outbox pattern or dead letter queue solution would be ideal.



#### we should add these features

- TODO: env setup should be reorganized adherence to 12 factor.
- TODO: we should consider to add JWT Authentication Middleware
- TODO: we should consider pagination for the GetSentMessages endpoint
- TODO: outbox pattern or dead letter queue can be suitable to retain data consistency. At least we ensure our data.
- TODO: is migrate tool able to bulk insertion.
- TODO: unit test should be added to the project
- TODO: integration test should be extended.


#### Useful links
- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/)

## the fin

  
