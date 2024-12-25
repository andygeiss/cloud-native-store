# Stage 1: Compile the service.
FROM golang:1.23.2 as build
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service ./cmd/service/main.go

# Stage 2: Create the service image.
FROM scratch
COPY --from=build /src/service .
CMD [ "/service" ]
