set dotenv-load

cache-dir := ".cache"
cert-dir := ".tls"

# Run the service.
run:
    @go run cmd/main.go

# Set up the service.
# Create a local CA and sign a server certificate.
# This will only be used if domains = ["localhost"].
setup:
    @brew install mkcert
    @rm -rf {{cache-dir}} ; mkdir {{cache-dir}}
    @rm -rf {{cert-dir}} ; mkdir {{cert-dir}}
    @mkcert -install
    @mkcert -cert-file {{cert-dir}}/server.crt \
        -key-file {{cert-dir}}/server.key \
        localhost 127.0.0.1 ::1

# Test the Go sources (Units).
test:
    @go test -v ./internal/app/core/services/...
