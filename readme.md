## fileverse

API for uploading and retrieving large file in IPFS storage.

* Used Domain Driven Design ( Clean Architecture)
* Handle Large Files (Efficiently With direct streaming from ipfs to http client)

Notes:
* I find "streaming" mechanism is more efficient than concurrent upload/download by chunks,
which I tried in this [PR](https://github.com/ashtishad/fileverse/pull/5)
* Also, IPFS itself is designed to handle upload/download efficiently by making Chunks, Content Addressing,
  Deduplication, Retrieval.

### Tools

* Golang
* Gin
* Postgresql(driver: pgx)
* [IPFS Node](https://github.com/ipfs/kubo)
* [IPFS client](https://pkg.go.dev/github.com/ipfs/go-ipfs-api)
* [Golangci-lint](https://golangci-lint.run/)


<!-- GETTING STARTED -->

### Getting Started

###### Clone using ssh protocol `git clone git@github.com:ashtishad/fileverse.git`

#### Environment-variables

Change environment variables in Makefile, if empty then default values listed here will be loaded, check
pkg/utils.go -> SanityCheck()

- API_HOST      `[IP Address of the machine]` : `127.0.0.1`
- API_PORT      `[Port of the  api]` : `8000`
- DB_USER       `[Database username]` : `postgres`
- DB_PASSWD     `[Database password]`: `postgres`
- DB_ADDR       `[IP address of the database]` : `127.0.0.1`
- DB_PORT       `[Port of the database]` : `5432`
- DB_NAME       `[Name of the database]` : `fileverse`
- GIN_MODE      `[Name of the gin mode]` : `debug`
- IPFS_ADDR     `[Address of ipfs docker node]`: `127.0.0.1:5001`

#### Postgres-Database-Setup

* Run docker compose: Bring the container up with `docker compose up`. Configurations are in `compose.yaml` file.


#### Run-the-application

* Run the application with `make run` command from project root. or, if you want to run it from IDE, please set
  environment variables by executing commands mentioned in Makefile on your terminal.

<p align="right"><a href="#fileverse">↑ Top</a></p>

<!-- Project Structure -->

### Project Structure
```

├── .github                    
│   └── workflows              
│       └── goci.yaml          <-- Github CI configuration file(for build, lint, test).
├── cmd                        
│   └── app                    
│       ├── app.go             <-- Server setup.
│       └── handlers.go        <-- Gin HTTP handlers for the file service.
├── config                     
│   └── initdb                 
│       └── 01.create-database.sql  <-- SQL script to create the initial database and tables.
├── docs                       
│   └── instructions.md        <-- Instructions given.
├── internal                   
│   ├── domain                 
│   │   ├── file.go            <-- Domain model for file.
│   │   ├── file_dto.go        <-- Data transfer object for file.
│   │   ├── file_repository.go <-- Interface for file repository.
│   │   └── file_repository_db.go  <-- Database implementation of the file repository.
│   ├── infra                  
│   │   ├── database           
│   │   │   ├── postgres.go    <-- Postgres database setup and configuration.
│   │   │   └── postgres_test.go   <-- Test setup for Postgres.
│   │   └── storage           
│   │       └── ipfs.go        <-- IPFS storage configuration and utilities.
│   └── service               
│       └── file_service.go    <-- Service layer that handles business logic for files.
├── pkg                       
│   └── utils                 
│       ├── api_errors.go      <-- API error definitions and utilities.
│       ├── sanity_check.go    <-- Environment variables setup.
│       └── slog_config.go     <-- Structured logging configuration.
├── .gitignore                 <-- Specifies intentionally untracked files to ignore.
├── .golangci.yml              <-- Configuration for golangci-lint tool.
├── compose.yaml               <-- Configuration for Docker compose services.
├── go.mod                     <-- Go module file, tracking dependencies.
├── main.go                    <-- Main application entry point.
├── Makefile                   <-- Defines set of tasks to be executed.
└── readme.md                  <-- Detailed project README file.



```

<p align="right"><a href="#fileverse">↑ Top</a></p>

<!-- Data Flow (Hexagonal architecture) -->

### Data Flow (Domain Driven Design)

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client

<p align="right"><a href="#fileverse">↑ Top</a></p>

### Routes

[POSTMAN WORKSPACE](https://www.postman.com/altimetry-cosmonaut-1609324/workspace/fileverse)

## Upload

POST 127.0.0.1:8000/upload
Body -> form-data(file) -> key: file

Example Response:

201 Created
```
{
"file": {
"fileId": "61d2d11c-b120-4215-86c8-a36cfef9a883",
"fileName": "large.zip",
"size": 139592028,
"timestamp": "2023-11-05T04:52:05.565843+06:00"
}
}
```

409 conflict

{
"error": "a file with the IPFS hash 'QmdDs4dzfRTo225qsQ2HU4HoCdpZkqPFpqQu9tkMjbz1RK' already exists"
}

500 error


## Get file

127.0.0.1:8000/file/:fileId
GET 127.0.0.1:8000/file/61d2d11c-b120-4215-86c8-a36cfef9a883

Example Response:

200 Ok

-> file content <-


500 Error
```
{
"error": "error retrieving file record"
}
```

<p align="right"><a href="#fileverse">↑ Top</a></p>
