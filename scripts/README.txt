How to use:

1: Edit ./config/config.yml to add the API keys and VERSION ("/push/vX")
2: Place the certificates PEM files into config folder (this folder will be mounted by the docker image at launch)
3: Execute "build.sh" to load the image contained in "jaraxasoftware_gorush.tar.gz" file
4: Modify "start.sh" and "test.sh" PORT variable with the desired value, and VERSION variable with the "/push/vX" value
5: Execute "start.sh" to start a container named "js-gorush"
6: Test if the service is running by executing "test.sh"
7: You can stop the container by executing "stop.sh"
8: You can now change some config values and/or certificates and restart the container with "start.sh"
9: You can also stop and remove both the container and the image by executing "remove.sh"
