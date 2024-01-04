cli:
	go install github.com/nats-io/natscli/nats@latest
	
nats:
	podman run --name nats -p 4222:4222 -p 6222:6222 -p 8222:8222 -v ./config/:/etc/nats --rm nats:latest -c /etc/nats/nats.conf -js

.PHONY: cli nats
