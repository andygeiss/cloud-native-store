# Stage 1: Compile the service.
FROM golang:1.23.2 as build
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service ./cmd/main.go

# Stage 2: Create the service image.
FROM scratch
COPY --from=build /src/service .
ENV ENCRYPTION_KEY="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
ENV SERVER_CERTIFICATE=".tls/server.crt"
ENV SERVER_DOMAIN="localhost"
ENV SERVER_KEY=".tls/server.key"
ENV TRANSACTIONAL_LOG=".cache/transactions.json"
EXPOSE 443
CMD [ "/service" ]
