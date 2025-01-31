# connect-api

Social networks are bloated, too much distractions and heavy web clients when everything we want is just real content. 

Connect-api aims to solve this, post and comment using only API clients (a.k.a Postman, curl, etc) as god planned.

## How to run

After cloning this project, you just have to run:
```bash
cp .env-example .env

docker compose up
```

The swagger docs should be available in http://localhost/swagger/index.html. But for now is recommended to export the json to PostmanAPI and use it from there.

### Techs
- Golang
- Gin
- Gorm
- Docker
- Nginx

### Roadmap

- [X] User authentication
- [X] User posts, crud
- [X] Make friends

