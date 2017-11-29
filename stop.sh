#!/bin/bash

# Import environment config
. envs

# Stop a docker container
docker stop $CONTAINER > /dev/null 2>&1;