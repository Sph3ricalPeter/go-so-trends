### ========== BUILD STAGE ========== ###
# pull latest Go image
FROM golang:latest AS build

# set working directory & copy files
WORKDIR /app
COPY . ./

# install dependencies
RUN make tidy

# build the app & init script
RUN make build
RUN make build-init

### ========== RUN STAGE ========== ###
# pull latest glibc image
FROM busybox:glibc as run

# set working directory & copy binaries & run.sh script
WORKDIR /app
COPY --from=build /app/out/main .

COPY --from=build /app/out/init .
COPY --from=build /app/data/ /app/data/

COPY --from=build /app/run.sh .

RUN chmod +x /app/run.sh

# expose port & run the server app
EXPOSE 8080
ENTRYPOINT ["/app/run.sh"]
