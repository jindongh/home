# home

# Setup
```
go mod init github.com/jindongh/home
go mod tidy
```

# Build
```
go build app.go
```


# Start
* Add service
```
cat > /etc/systemd/system/home.service <<EOF
[Unit]
Description=home web site
Wants=network-online.target
After=network-online.target nss-lookup.target

[Service]
Type=exec
Restart=always
RestartSec=1
User=pi
UMask=0000
ExecStart=/home/pi/home/app

[Install]
WantedBy=multi-user.target
EOF
```

* Start service
```
sudo systemctl daemon-reload
sudo systemctl start home
sudo systemctl enable home
```

