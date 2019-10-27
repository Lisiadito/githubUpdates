# Docker build
To dockerize "githubUpdates" use the following build command inside of the "docker" directory.

    `sudo docker build --no-cache -t githubupdates .`

After it is build run 

    `sudo docker run -dit --env-file ../.env githubupdates:latest`

If you want to attach to the container you can run

    `sudo docker exec -it <CONTAINER ID> bash`