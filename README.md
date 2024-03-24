# Koofr HEIC

Convert HEIC to JPG on Koofr.

## Getting started

```sh
go install github.com/revel/cmd/revel

go get github.com/bancek/koofr-heic

export KOOFR_CLIENT_ID="CLIENTID"
export KOOFR_CLIENT_SECRET="CLIENTSECRET"
export KOOFR_REDIRECT_URL="http://localhost:9000/App/Auth"
export APP_SECRET="APPSECRET"

revel run github.com/bancek/koofr-heic
```

Now go to http://localhost:9000/

## Deploy to Docker

```sh
docker build -t koofr-heic .

docker run -d -e KOOFR_CLIENT_ID="CLIENTID" -e KOOFR_CLIENT_SECRET="CLIENTSECRET" -e KOOFR_REDIRECT_URL="http://localhost:8000/App/Auth" -e APP_SECRET="APPSECRET" -p 8000:9000 koofr-heic
```
