FROM golang:1.17-alpine3.15

RUN apk update && apk add --no-cache git curl

WORKDIR /app
COPY go.mod go.sum ./
COPY .air.toml .
RUN go mod download
COPY . .

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
CMD ["air"]
