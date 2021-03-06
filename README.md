# ormesh - onion-routed mesh

ormesh helps you build private backplane to connect infrastructure through
[Tor](https://www.torproject.org/).

## Why?

To abstract away geography and network topologies.

To access services running almost anywhere, from just about anywhere else.

Tor is well-suited to traversing all kinds of networks between services and the
clients that would consume them. Tor provides a resilient infrastructure with
no single point of failure.

Tor's hidden services can be deployed in a private, authenticated mode, which
keeps services from being generally accessible.

ormesh helps manage the configuration and auth token exchange necessary to
deploy a private backplane to connect infrastructure.

## What kind of services?

Text-based protocols (email, messaging), sensors & actuators, home automation,
and file synchronization are just some ideas to get you started.

In general, services that require little bandwidth or tolerate latency. With
ormesh, they can be accessed without the hassle of setting up iptables, NAT
port forwarding, VPNs, TLS, and without relying on central rendezvous servers.

## What ormesh isn't

ormesh is not a VPN in the conventional sense.

ormesh is not intended for operating unauthenticated anonymous hidden services.
Anonymity is an interesting side-effect of building on Tor, but it is not a
primary goal for ormesh, nor it is guaranteed. Users are responsible for
evaluating ormesh (and its Tor configuration) and deciding whether it meets
security requirements and threat models.

Low-latency, high bandwidth applications may not perform well over ormesh's Tor
configuration. Improvements here are possible (by trading anonymity for
improved latency and network throughput) but not yet implemented.

Also keep in mind that Tor only routes TCP traffic.

# Installing

[Install Go](https://golang.org/doc/install).

Add $GOPATH/bin to your $PATH.

Download and build ormesh.

    go get -u github.com/cmars/ormesh

## macOS

[Install Tor Browser](https://www.torproject.org/download/download-easy.html.en). On macOS,
ormesh will launch the `tor.real` executable bundled with the Tor Browser.

## Others

[Install Tor (standalone)](https://www.torproject.org/download/download-unix.html.en) from
official packages for your platform.

## (Experimental) Ubuntu 16.04 LTS amd64

    sudo snap install ormesh --edge

The snap packaging bundles Tor from official archives and operates it in a
sandbox. Why it's experimental quality:

- Needs an automated build set up to track official Tor releases. Until then,
  it's possible that the version of Tor bundled in the package is out-of-date.
- Installing snaps into containers may not work, depending on the container
  image, host kernel and other possible factors.

# Configuring

## Exporting local services

A machine, VM or container running a web application might export access to
HTTP and SSH running locally.

```
$ ormesh export add 22
$ ormesh export add 80
$ ormesh status
service:
  export:
    - 127.0.0.1:22
    - 127.0.0.1:80
```

You can also export services on other hosts.

```
$ ormesh export add 192.168.1.19:8000
```

## Adding clients

Generate a token string that a client can import to access exported services.
This string should be securely sent to the user of `my-MacBook` who is granted
access.

```
$ ormesh client add my-MacBook
Y5Cfw7A5RhP8Rd7xGYfD8N4oyEBpBWNR+6Qkgrbepk0=
```

Status shows the ports exported and the clients authorized to access them.

```
$ ormesh status
service:
  export:
    - 127.0.0.1:22
    - 127.0.0.1:80
  clients:
    - name: my-MacBook
      address: q6jo2z3bw5exkece.onion
```

# Consuming services

## Add a remote service, with client authentication

On the machine `my-MacBook`, add a remote using the authentication string displayed by
`client add` above.

```
$ ormesh remote add website Y5Cfw7A5RhP8Rd7xGYfD8N4oyEBpBWNR+6Qkgrbepk0=
```

```
$ ormesh status
remotes:
  - name: website
    address: q6jo2z3bw5exkece.onion
```

```
$ ormesh remote show website
q6jo2z3bw5exkece.onion
```

## Display an SSH config entry

Display an ssh-config(5) stanza for the remote.

```
$ ormesh remote ssh-config website
Host website
  ProxyCommand nc -X 5 -x localhost:9250 %h %p
  Hostname q6jo2z3bw5exkece.onion
```

## Importing remote services

Set up local port forwarding to remote services with _imports_.

Forward local port 10022 to port 22 on the remote:

```
$ ormesh import add website 22 127.0.0.1:10022
$ ormesh status
remotes:
  - name: website
    address: q6jo2z3bw5exkece.onion
    imports:
    - local-addr: 127.0.0.1:10022
      remote-port: 22
```

Listen on all addresses to create a public ingress to a remote service. Useful
for circumventing inbound port blocks where the service is running. For
example, you want to physically locate your email server in a mobile camper,
your ISP blocks SMTP inbound, and your IP address changes often. Import your
services from a cloud instance with a public IP and DNS.

```
$ ormesh import add mailinabox 25 0.0.0.0:25
$ ormesh import add mailinabox 587 0.0.0.0:587
```

# Operating

An ormesh agent launches tor and operates it automatically, based on ormesh
configuration.

## Launching manually

```
$ ormesh agent run
```

will run the agent as configured above. Configuration changes made while the agent
is running are applied immediately.

## Setting up systemd

Display a systemd unit file that will run ormesh, from its current installed
binary path.

```
$ ormesh agent systemd-unit --user
[Unit]
Description=ormesh - onion-routed mesh

[Service]
ExecStart=/path/to/ormesh agent run
Restart=always

[Install]
WantedBy=default.target
```
