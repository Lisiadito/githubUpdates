# Docker build
Before you start dockerize the application, you need to copy your
**.env** file to the docker directory. Note that the following 
commands expects to be executed in the project root directory.

    cp .env docker

To dockerize "githubUpdates" use the following build command inside of the "docker" directory.

    docker build -t githubupdates .

Please note, that root rights are needed to run docker related commands.
