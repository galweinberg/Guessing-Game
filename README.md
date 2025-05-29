#docker deployment instructions  -
-----------------------------------------------------------------------------------
Multiplayer Game - Deployment Guide-

Pre-requisites: What software needs to be installed?
Before users can run the game, they need the following:

Docker: Docker is required to build and run the game in a container.


Git: Git is needed to clone the repository from GitHub.


Go (Optional for Development):

If the user wants to modify the game or run it locally without Docker, Go should be installed.

----------------------------------------------------------------------------------
ii. Clone the Repository:
Clone the repository to your local machine:

in bash:

git clone https://github.com/your-username/Guessing-game.git
cd Guessing-Game
------------------------------------------------------------------------------------------
iii. Building the Docker Image:
Navigate to the folder containing the Dockerfile:

bash-
cd game
Build the Docker image:

bash-
docker build -t multiplayer-game-task .

------------------------------------------------------------------------------------------
iv. Running the Docker Container:
Run the container with:

bash-
docker run -p 8080:8080 multiplayer-game-task

The game will now be running, and players can connect via the provided terminal or client interface.
------------------------------------------------------------------------------------------
v. How to Connect and Play the Game:
Once the server is running, users can connect to the game through a terminal using the client application (if provided).

Players will take turns guessing the secret number, with each turn lasting 20 seconds.

Type "exit" during your turn to leave the game.
