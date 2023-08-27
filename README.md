# Golang-internship
Colța Maria

## Run Application with Docker

More information about [Docker](https://www.docker.com/).
To run the application type this command in the root folder.

```bash
$ docker compose up --build
```

You might have to run this command twice if it doesn't work the first time :)

## API Docs

To access swagger-ui go to [localhost:8080/api/swagger/index.html](http://localhost:8080/api/swagger/index.html)  

To generate new docs you need [swag](https://github.com/swaggo/swag) installed and added to `PATH`. Then type this command in the root folder.

```bash
$ swag init --parseDependency
```