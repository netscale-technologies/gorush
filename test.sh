#!/bin/bash

# Import environment config
. envs
URL="localhost:${PORT}${VERSION_PATH}/stats/test"

curl \
	-XGET \
	-H "Accept: application/json" \
 	$URL | python -mjson.tool