#!/bin/bash
PS3="Choose an service to install: "
options=(drone remote)
select servicename in "${options[@]}";
do
    if [[ $servicename == "drone" ||  $servicename == "remote" ]]; then
        break;
    fi
done

sudo systemctl stop $servicename 
sudo systemctl disable $servicename
sudo systemctl daemon-reload
sudo rm "/etc/systemd/system/$servicename.service"
