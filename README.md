# btdaemon

A tiny program for enabling / disabling Bluetooth based on the SSID your host 
is connected to.

## Why?

Bluetooth is susceptible to a number of security issues. To avoid those,
this program only enables Bluetooth when your host is connected to a known
network.

It's convenient to have Bluetooth enabled even when devices aren't connected
when using devices that power down (and disconnect) while idle (like mice).

## Caveats

I wrote this program for myself. Which means there's a few bits that are hardcoded:

- You must have [blueutil](https://github.com/toy/blueutil) installed.
- The `airport` command must be present and located at 
  `/System/Library/PrivateFrameworks/Apple80211.framework/Resources/airport`.
- Logs are always printed to `/var/log/net.codeviking.btdaemon/stdout.log`.

## Installation

Build a binary:

```
go build -o btdaemon ./main.go
```

Put it in `/usr/local/bin`:

```
sudo mv btdaemon /usr/local/bin/.
sudo chown root:wheel /usr/local/bin/btdaemon
sudo chmod 755 /usr/local/bin/btdaemon
```

Put the plist file in `/Library/LaunchAgents`:

```
sudo cp net.codeviking.btdaemon.plist /Library/LaunchAgents
sudo chown root:wheel /Library/LaunchAgents/net.codeviking.btdaemon.plist
sudo chmod 644 /Library/LaunchAgents/net.codeviking.btdaemon.plist
```

Make a directory for the daemon's logs. The permission juggling is important,
as the agent runs as your user:

```
sudo mkdir /var/log/net.codeviking.btdaemon
sudo chown $(id -u -n):$(id -g -n) /var/log/net.codeviking.btdaemon
```


Load and start the daemon:

```
launchctl load /Library/LaunchAgents/net.codeviking.btdaemon.plist
```

See what it's doing:

```
tail -f /var/log/net.codeviking.btdaemon/stdout.log
```

To stop the daemon run:

```
launchctl unload /Library/LaunchAgents/net.codeviking.btdaemon.plist
```


