#!/bin/bash

# Import environment config
.envs

# Load a docker image from a .tar.gz file 
gunzip < $IMAGE_FILE.tar.gz | docker load
