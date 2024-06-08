# Notifier for High Memory Usage

A simple tool to notify you when your system memory usage is too high.

Features:

- uses `notify-send` to send notifications when memory usage reaches a certain threshold
- major Linux distributions are supported
- configurable threshold, check memory usage interval and resend notification delay

Example desktop notification:

![notifymem](https://raw.githubusercontent.com/shayanderson/notify-mem/main/notifymem.png)

## Usage

Example usage that monitors memory usage by checking every 2 seconds, sends a notification when memory usage is above 80% and will not send another notification until after 10 seconds.

```bash
notifymem -threshold 80 -delay 10 -interval 2
```

Run a test which will check memory usage and send a test notification:

```bash
notifymem -test
```

## Installation

First, download the latest release or clone the repository and build locally using `make build`, which will create a binary in the `bin` directory. Then, copy the binary to a directory, like `/opt/notifymem/bin`, and make the file executable, like `chmod +x /opt/notifymem/bin/notifymem`.

Next, create a systemd service file in `/etc/systemd/system/notifymem.service` with the following content:

```ini
[Unit]
Description=notifymem service: notify when memory usage reaches a threshold
After=network.target,systemd-user-sessions.service,systemd-journald.service

[Service]
Type=simple
Restart=always
RestartSec=3
User=your-username
ExecStart=/opt/notifymem/bin/notifymem -threshold 80 -delay 30 -interval 2
Environment="DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus"

[Install]
WantedBy=multi-user.target
```

Make sure to:

- replace `your-username` with your actual username
- set `ExecStart` to the path where you copied the binary and use the desired options
- set `Environment` to the correct `DBUS_SESSION_BUS_ADDRESS` value for your system, which can be found by running `echo $DBUS_SESSION_BUS_ADDRESS`

Finally, enable and start the service:

```bash
systemctl enable notifymem
systemctl start notifymem
```

Tips:

- check the status of the service using `systemctl status notifymem`
- check the logs of the service using `journalctl -u notifymem` (follow the logs using `-f`)
- reload the systemd daemon after making changes to the service file using `systemctl daemon-reload`
