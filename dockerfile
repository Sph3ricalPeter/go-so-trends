FROM golang:latest AS build
WORKDIR /app
COPY go.mod main.go ./
RUN go mod tidy
RUN go build -o main main.go

FROM busybox:glibc as run
WORKDIR /app
COPY --from=build /app/main .

EXPOSE 8080
ENTRYPOINT ["/app/main"]
