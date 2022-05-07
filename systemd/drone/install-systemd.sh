sudo cp ./systemd/drone.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable drone