Author: Mark Saravi  
Description: Building a drone from scratch with golang  
**This project is still under development and informations in this repo is not enough for a complete drone yet**

# Building your own drone from scratch

## Safety First

Building your own drone from scratch can be very exciting and there is lots of learning but it also can be a hazardous and dangerous hobby if you don't pay attention about the safety. Multiple scars on my fingers are my witnesses and believe me I was careful about safety at least for this project.  
When you test your drone you expose yourself to a system that is under development and the reaction of system can be buggy abd unpredictable. Drones use Brushless motors and although they are small, they are a powerful beast and combine with a propeller which has multiple thousands RPM (imagine something between 3000 to 7000 RPM) it is killing machine. Any bug in your control system can cause sudden out of control move and jump and if it hits someone at best case scenario there will be multiple bleeding wounds and worst case scenario blind eyes. So if you are going to use my code and design to build your own drone please pay attention to the following safety measures:  
**Important**:  never ever install and use the propellers (don't even buy one) until you finish the following steps and
make sure that you have clear understanding about how to turn on and control a Brushless motor. This includes knowing how to calibrate ESC and how to run it. 

## Propeller last
- DON'T BUY AND INSTALL PROPELLERS UNTIL YOU FINISH THESE STEPS.
- Never test or run a Brushless motor without securing and tightening it with proper screws to a fixed position. (DON'T ATTACH THE PROPELLERS)
- Continuously check the screws as screws can get lose due to motor's vibration.
- Make sure you can cut the power of Brushless motors anytime. (DON'T TEST WITH BATTERY AND ALWAYS USE A POWER SUPPLY WHICH OFF SWITCH IS ALWAYS NEXT TO YOU)  
- Always have a power limit in system. My recommendation is 25% for initial tests and learning how to control Brushless motors and initial tests. 40% for flight control tunings. Even when you fly the drone, you hardly need more than 50% of power.
- Test your control software and make sure it always has a safety limit for the applied powers to motors.
- Test your control software against moves and initial conditions. For example your drone might accumulate power by **Integral** control (referring to PID control systems) while it is tilted on floor and increasing the power can end up with violent jump and flipping.  

## Testing with Propellers
- Always use a safety glass or eye protection when you are testing and tuning with propellers.
- If you want to start testing with propellers never test it where someone else can be there. Always make sure no one else is around or can get close to you while you are testing.
- It is better to run initial tests and tuning in a closed space to make sure there is always something like a wall or a ceiling to stop it.
- If you want to test outside make sure no one is around and you are confident that you can stop motors anytime.
- Make sure you are not disobeying any safety regulations based on your local rules.
 

## Install Go in Raspberry Pi

### By apt-get
`sudo apt-get install golang`  

### By precompiled package

get the latest version of precompiled package from [here](https://golang.org/dl/). I used *go1.15.2.linux-armv6l.tar.gz*  
```
wget https://storage.googleapis.com/golang/go1.16.6.linux-armv6l.tar.gz  
sudo tar -C /usr/local -xzf ./go1.16.6.linux-armv6l.tar.gz   
export PATH=$PATH:/usr/local/go/bin
or
add PATH=$PATH:/usr/local/go/bin to your shell profile
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
