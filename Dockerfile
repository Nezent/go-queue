FROM golang:1.23-alpine

# Install make, curl, and tzdata
RUN apk add --no-cache make curl tzdata

ENV TZ=Asia/Dhaka

WORKDIR /usr/src/app
COPY . .

RUN go mod download

# Specify the command to run when starting the container
CMD ["make", "run"]

