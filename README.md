# Leptop

Lepton is a journald log forwarder.

If fetched log entries from `systemd-journal-gatewayd` and sends them to an adapter.
Currently there is only a loggly adapter.

## Usage

```
docker run -it --net=host \
    -v /etc/ssl/:/etc/ssl/:ro \
    -v /usr/share/ca-certificates/:/usr/share/ca-certificates/:ro \
    -e LOGGLY_TOKEN=<your loggly token> \
    pulcy/lepton:latest
```
