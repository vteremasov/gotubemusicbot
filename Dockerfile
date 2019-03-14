# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.11 base image
FROM golang:1.12

# Add Maintainer Info
LABEL maintainer="Me"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/vteremasov/go-music-bot

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go build

RUN chmod +x 'go-music-bot'

RUN apt-get update
RUN apt-get -y install ffmpeg

# This container exposes port 80 to the outside world
EXPOSE 80/tcp

VOLUME ["/go-music-bot/logs"]

# Run the executable
CMD ["go-music-bot"]

#FROM alpine
#ENV LANGUAGE="en"
#COPY /src .
#RUN apk add --no-cache ca-certificates &&\
#    chmod +x code
#
#CMD [ "./bot" ]