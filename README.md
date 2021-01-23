# tproxy

tproxy is an experimental project for replaying tcp traffic.

## HOW-TO

There are two namespaces, client and server, and a bridge. They communicate by veth interfaces and the bridge.

![TProxy with bridge](./img/linux-tproxy-with-bridge.png)

In namespace client, its settings are:

```shell
# veth interface is vtci
ip link set dev vtci up
ip addr add 10.0.0.2/24 dev vtci
ip route add default via 10.0.0.1 dev vtci
```

In namespace server, its settings are:

```shell
# veth interface is vtsi
ip link set dev vtsi up
ip addr add 10.0.0.1/24 dev vtsi
ip route add default via 10.0.0.2 dev vtsi
```

Reference `./script/setup.sh`.

All IP packets in namespace client go through `vtci` into namespace server. Program of namespace server receive IP packets from `vtsi`.
