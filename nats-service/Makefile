cli:
	go install github.com/nats-io/natscli/nats@latest
	
nats:
	podman run --name nats -p 4222:4222 -p 6222:6222 -p 8222:8222 --rm nats:latest

cert:
	mkcert -install
	mkcert -cert-file config/server-cert.pem -key-file config/server-key.pem localhost 127.0.0.1 ::1

.PHONY: cli nats cert
