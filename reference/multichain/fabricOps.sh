#!/bin/bash

docker-compose -f ./docker-compose.yml run --rm assets.generator.fabric.network /bin/bash -c '${PWD}/generate-artifacts.sh'