# Release Aware

![How it looks like](example.png?raw=true "How it looks like")

## Prerequisites

In order to run the commands described below, you need:
- [Docker](https://www.docker.com/) 
- [Docker Compose](https://docs.docker.com/compose/)
- make (`sudo apt-get install make` on Ubuntu)

## Usage

```
git clone ...
docker-compose build
docker-compose up
```

## Release

```
./release.sh <VERSION>
```

First you need to login to Docker Hub.
This is done by the SDP team.
