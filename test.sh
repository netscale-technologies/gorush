#!/bin/bash

# Import environment config
.envs

curl \
	-XGET \
	-H "Accept: application/json" \
 	"localhost:$PORT$VERSION/stats/test" | python -mjson.tool