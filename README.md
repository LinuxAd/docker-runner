# Docker Runner

__Currently under construction!__

## Requirements

* Go 1.16+
* Docker (tested with Mac Desktop 20.10.5)

## Usage

Quick test:

```bash
go build
./docker-runner
curl -X POST @test_service.json
docker ps
```

## Current Bugs

The runner locks up after two containers are spun up, constantly trying to create a new container that has the same name as the second container.