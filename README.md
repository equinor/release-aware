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
cp server.env.template server.env
# populate variables in server.env as described below
docker-compose build
docker-compose up
```

* GITHUB_TOKEN: login to GitHub and generate a Personal access token in Personal settings -> Developer settings -> Personal access tokens

## Release

```
./release.sh <VERSION>
```

First you need to login to Docker Hub.
This is done by the SDP team.
