# first stage: build the Go application
FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go


# second stage: run the Go application
FROM scratch
COPY --from=build /app/main /main
CMD ["/main"]