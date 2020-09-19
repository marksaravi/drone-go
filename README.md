# drone-go

Author: Mark Saravi  
Description: A drone project with Go  

## Install Go in Raspberry Pi

### By apt-get
`sudo apt-get install golang`  

### By precompiled package

get the latest version of precompiled package from [here](https://golang.org/dl/). I used *go1.15.2.linux-armv6l.tar.gz*  
```
wget https://storage.googleapis.com/golang/go1.15.2.linux-armv6l.tar.gz  
sudo tar -C /usr/local -xzf go1.7.3.linux-armv6l.tar.gz   
export PATH=$PATH:/usr/local/go/bin  
```
Add */usr/local/go/bin* to your profile (e.g ~/.bashrc).  
Check the version by `go version`
Run the following commands to create default folders (no need to setup environment variables like GOPATH)
```
mkdir ~/go/
mkdir ~/go/src/
mkdir ~/go/pkg/
mkdir ~/go/bin/
```