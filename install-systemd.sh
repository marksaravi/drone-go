#!/bin/bash
PS3="Choose an service to install: "
options=(drone drone-remote)
select servicename in "${options[@]}";
do
    if [[ $servicename == "drone" ||  $servicename == "drone-remote" ]]; then
        break;
    fi
done
sudo cp "./systemd/$servicename.service" /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable $servicename 
