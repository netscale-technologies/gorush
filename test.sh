#!/bin/bash

PORT=10422
VERSION="/push/v2"

curl \
	-XGET \
	-H "Accept: application/json" \
 	"localhost:$PORT$VERSION/stats/test" | python -mjson.tool