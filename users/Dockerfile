FROM golang:1.18 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./build/app cmd/main.go


FROM golang:alpine
# https://stackoverflow.com/questions/66963068/docker-alpine-executable-binary-not-found-even-if-in-path/66974607#66974607
RUN apk update && apk add --no-cache libc6-compat gcompat
WORKDIR /usr/src/
COPY --from=builder /app/build/app /usr/src/app
ENTRYPOINT ./app

