FROM golang:latest AS build
WORKDIR /app
COPY . ./
RUN make tidy
RUN make build

FROM busybox:glibc as run
WORKDIR /app
COPY --from=build /app/out/main .

EXPOSE 8080
ENTRYPOINT ["/app/main"]
