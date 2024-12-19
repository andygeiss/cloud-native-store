set dotenv-load

# Create a local CA and sign a server certificate.
# This will only be used if domains = ["localhost"].
cert-dir := ".tls"
make-certs:
    @brew install mkcert
    @rm -rf {{cert-dir}} ; mkdir {{cert-dir}}
    @mkcert -install
    @mkcert -cert-file {{cert-dir}}/server.crt \
        -key-file {{cert-dir}}/server.key \
        localhost 127.0.0.1 ::1

# Run the service.
run:
    @go run cmd/main.go

# Test the Go sources (Units).
test:
