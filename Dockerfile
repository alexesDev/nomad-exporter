FROM golang:1.15.6-alpine3.12 as build

RUN apk add --update git
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -o server .

FROM scratch
COPY --from=0 /app/server .
CMD ["./server"]
