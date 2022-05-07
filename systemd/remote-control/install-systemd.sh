sudo cp ./systemd/drone-remote-control.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable drone-remote-control 