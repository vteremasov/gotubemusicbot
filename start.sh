
set -e

docker build -t go-music-bot .

docker-compose up --build
