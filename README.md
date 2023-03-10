# movies_api
A simple api backend program using golang and docker

## Prerequisites
- WSL2
- Docker
- Go 1.19

## Run
1. Define .env, can be seen in .env.example
2. run command "docker-compose up -d"
3. Manually test API through postman
4. Run go test -v in local machine to run automated test

## Important Notes
- After making changes to .env, run docker-compose down first, then docker-compose up -d to reload .env 
- To run automated test, make sure TEST_ENVIRONMENT variable in env is set to staging to avoid database conflict

## API Endpoints and body

### /Movies GET
#### Response code
- 200
- 404

#### Response body
[
    {
        "id": 1,
        "title": "Pengabdi Setan 2 Comunion",
        "description": "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
        "rating": 7,
        "image": "",
        "created_at": "2022-08-01 10:56:31",
        "updated_at": "2022-08-13 09:30:23"
    }
]

### /Movies/:ID GET
#### Response code
- 200
- 404

#### Response body
{
    "id": 1,
    "title": "Pengabdi Setan 2 Comunion",
    "description": "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
    "rating": 7,
    "image": "",
    "created_at": "2022-08-01 10:56:31",
    "updated_at": "2022-08-13 09:30:23"
}

### /Movies POST
#### Request body
{
    "id": 1,
    "title":"Pengabdi Setan 2 Comunion",
    "description" : "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
    "rating":7,
    "image":"",
    "created_at":"2022-08-01 10:56:31",
    "updated_at":"2022-08-13 09:30:23"
}

#### Response code
- 200
- 400

### /Movies/:ID PATCH
#### Request body
{
    "id": 1,
    "title":"Pengabdi Setan 2 Comunion",
    "description" : "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
    "rating":7,
    "image":"",
    "created_at":"2022-08-01 10:56:31",
    "updated_at":"2022-08-13 09:30:23"
}

#### Response code
- 200
- 400

### /Movies/:ID DELETE
#### Response code
- 200

