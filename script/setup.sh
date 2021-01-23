#!/bin/bash -ex

VETH_CLIENT_INNER="vtci"
VETH_CLIENT_OUTER="vtco"
VETH_SERVER_INNER="vtsi"
VETH_SERVER_OUTER="vtso"
NS_CLIENT="nsclient"
NS_SERVER="nsserver"
IP_CLIENT="10.0.0.2"
IP_SERVER="10.0.0.1"
BR_TPROXY="brtproxy"
LISTEN_PORT="443"

# create veth pair for client
ip link add ${VETH_CLIENT_INNER} type veth peer name ${VETH_CLIENT_OUTER}

# prepare namespace client
ip netns add ${NS_CLIENT}
ip link set ${VETH_CLIENT_INNER} netns ${NS_CLIENT}
ip netns exec ${NS_CLIENT} bash -ex -c "
ip link set dev lo up
ip link set dev ${VETH_CLIENT_INNER} up
ip addr add ${IP_CLIENT}/24 dev ${VETH_CLIENT_INNER}
ip route add default via ${IP_SERVER} dev ${VETH_CLIENT_INNER}
"

# create veth pair for server
ip link add ${VETH_SERVER_INNER} type veth peer name ${VETH_SERVER_OUTER}

# prepare namespace server
ip netns add ${NS_SERVER}
ip link set ${VETH_SERVER_INNER} netns ${NS_SERVER}
ip netns exec ${NS_SERVER} bash -ex -c "
ip link set dev lo up
ip link set dev ${VETH_SERVER_INNER} up
ip addr add ${IP_SERVER}/24 dev ${VETH_SERVER_INNER}
ip route add default via ${IP_CLIENT} dev ${VETH_SERVER_INNER}
"

# prepare namespace client's tproxy
ip netns exec ${NS_CLIENT} bash -ex -c "
iptables -t mangle -N DIVERT
iptables -t mangle -A PREROUTING -i ${VETH_CLIENT_INNER} -p tcp -m socket -j DIVERT
iptables -t mangle -A DIVERT -j MARK --set-mark 1
iptables -t mangle -A DIVERT -j ACCEPT
ip rule add fwmark 1 lookup 100
ip route add local default dev ${VETH_CLIENT_INNER} scope host table 100
"

# prepare namespace server's tproxy
ip netns exec ${NS_SERVER} bash -ex -c "
iptables -t mangle -N DIVERT
iptables -t mangle -A PREROUTING -p tcp -m socket -j DIVERT
iptables -t mangle -A DIVERT -j MARK --set-mark 1
iptables -t mangle -A DIVERT -j ACCEPT
iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port ${LISTEN_PORT}
ip rule add fwmark 1 lookup 100
ip route add local default dev ${VETH_SERVER_INNER} scope host table 100
"

# prepare bridge  // bridge-utils is required
brctl addbr ${BR_TPROXY}
brctl addif ${BR_TPROXY} ${VETH_CLIENT_OUTER}
brctl addif ${BR_TPROXY} ${VETH_SERVER_OUTER}
ip link set dev ${VETH_CLIENT_OUTER} up
ip link set dev ${VETH_SERVER_OUTER} up
ip link set dev ${BR_TPROXY} up

# bridge runs with docker0
iptables -t filter -I FORWARD -i ${BR_TPROXY} -o ${BR_TPROXY} -j ACCEPT
