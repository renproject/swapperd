sudo echo "[Unit]
Description=RenEx's Swapper Daemon
After=network.target

[Service]
ExecStart=/home/ubuntu/.swapper/bin/swapper 
Restart=on-failure
StartLimitBurst=0

# Specifies which signal to use when killing a service. Defaults to SIGTERM.
# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)
KillSignal=SIGHUP

[Install]
WantedBy=default.target" >> /etc/systemd/system/swapper.service

sudo systemctl daemon-reload
sudo systemctl enable swapper.service
sudo systemctl start swapper.service
