network:
	podman network prune
	podman network create nats
	
node1:
	podman run --name node1 --network nats -p 4222:4222 -v ./config/:/etc/nats --rm nats:latest -c /etc/nats/node1.conf

node2:
	podman run --name node2 --network nats -p 4223:4223 -v ./config/:/etc/nats --rm nats:latest -c /etc/nats/node2.conf

node3:
	podman run --name node3 --network nats -p 4224:4224 -v ./config/:/etc/nats --rm nats:latest -c /etc/nats/node3.conf

natsbox:
	podman run --network nats --rm -it natsio/nats-box

.PHONY: network node1 node2 node3 natsbox

