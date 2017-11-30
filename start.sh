#!/bin/bash

# Import environment config
. envs

# Remove previous container
docker rm $CONTAINER > /dev/null 2>&1

docker run -ti -d --name $CONTAINER --restart always \
	-p $PORT:8088 \
	-v $DIR/config:/config:ro \
	$DEPLOY_ACCOUNT/$DEPLOY_IMAGE:latest -c /config/config.yml