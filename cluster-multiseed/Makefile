node1:
	podman run --name node1 -p 4222:4222 -p 6222:6222 -v ./config/:/etc/nats --rm nats:latest -c /etc/nats/node1.conf

node2:
	podman run --name node2 -p 4223:4223 -p 6223:6223 -v ./config/:/etc/nats --rm nats:latest -c /etc/nats/node2.conf

node3:
	podman run --name node3 -p 4224:4224 -p 6224:6224 -v ./config/:/etc/nats --rm nats:latest -c /etc/nats/node3.conf

natsbox:
	podman run --rm -it natsio/nats-box

goreman:
	goreman start

.PHONY: node1 node2 node3 natsbox goreman

