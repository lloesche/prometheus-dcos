# [Prometheus](https://prometheus.io/) on [Mesosphere DC/OS](https://dcos.io/)

## Intro
This runs Prometheus on DC/OS. `server.json` contains the service definition for Prometheus itself. `node_exporter.json` contains the service definition for node_exporter. I'm running node_exporter inside a Mesos (cgroups) container so that it sees all of the hosts filesystems without any need for priviliges or translation.

## Environment Variables
| Variable | Function | Example |
|----------|----------|-------|
|`NODE_EXPORTER_SRV` | Mesos-DNS SRV record of the node_exporter | `NODE_EXPORTER_SRV=_node-exporter.prometheus._tcp.marathon.mesos`|
|`SRV_REFRESH_INTERVAL` | How often should we update the targets | `SRV_REFRESH_INTERVAL=60`|
|`ALERT_MANAGER_URI` | AlertManager URL | `ALERT_MANAGER_URI=http://prometheusalertmanager.marathon.l4lb.thisdcos.directory:9093`|

## Building the SRV lookup helper
To run the srv-lookup helper tool inside the minimal prom/prometheus Docker container I statically linked it. To do so yourself install [musl libc](http://www.musl-libc.org/) and compile using:
```
$ CC=/usr/local/musl/bin/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' srv-lookup.go
```
