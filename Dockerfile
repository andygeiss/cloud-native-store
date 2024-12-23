# Stage 1: Compile the service.
FROM golang:latest as build
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service ./cmd/main.go

# Stage 2: Create the service image.
FROM scratch
COPY --from=build /src/service .
EXPOSE 443
CMD [ "/service" ]