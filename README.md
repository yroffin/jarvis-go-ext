# jarvis-go-ext
Jarvis golang extention

[![Build Status](https://travis-ci.org/yroffin/jarvis-go-ext.svg?branch=master)](https://travis-ci.org/yroffin/jarvis-go-ext)

# Dependencies and setup

## MongoDB

Setup mongodb for store local data

        sudo apt-get install mongodb-server

## WiringPI

Connector is based on wiringPi so setup wiringPi on your host [here](http://wiringpi.com)

        sudo apt-get install wiringpi

## Teleinfo option

For teleinfo option change right to 755 on USB file (ex: /dev/ttyUSB0) and apply this command in /etc/rc.local

        stty 1200 cs7 evenp cstopb -igncr -inlcr -brkint -icrnl -opost -isig -icanon -iexten -F /dev/ttyUSB0

## Setup on raspberry pi 2 or zero

        pi@raspberrypi:~ $ sudo userdel -r jarvis
        pi@raspberrypi:~ $ sudo useradd -m jarvis
        pi@raspberrypi:~ $ export GITHUB=https://github.com/yroffin/jarvis-go-ext/releases/download/1.03a
        pi@raspberrypi:~ $ sudo wget ${GITHUB}/jarvis-rest-module-0.0.1-SNAPSHOT.armel -O /home/jarvis/jarvis-rest-module-0.0.1-SNAPSHOT.arm
        or
        pi@raspberrypi:~ $ sudo wget ${GITHUB}/jarvis-rest-module-0.0.1-SNAPSHOT.armhf -O /home/jarvis/jarvis-rest-module-0.0.1-SNAPSHOT.arm
        pi@raspberrypi:~ $ sudo chmod 755 /home/jarvis/jarvis-rest-module-0.0.1-SNAPSHOT.arm
        pi@raspberrypi:~ $ sudo chown jarvis:jarvis /home/jarvis/jarvis-rest-module-0.0.1-SNAPSHOT.arm
        pi@raspberrypi:~ $ sudo wget ${GITHUB}/jarvis-go-service -O /etc/init.d/jarvis-go-service
        pi@raspberrypi:~ $ sudo chmod 755 /etc/init.d/jarvis-go-service
        pi@raspberrypi:~ $ sudo update-rc.d jarvis-go-service defaults
        pi@raspberrypi:~ $ sudo service jarvis-go-service restart
        pi@raspberrypi:~ $ curl http://192.168.1.47:7000/api/teleinfo

## /etc/rc.local sample

        # setup jarvis
        stty 1200 cs7 evenp cstopb -igncr -inlcr -brkint -icrnl -opost -isig -icanon -iexten -F /dev/ttyUSB0
        service mongodb restart
        service jarvis-go-service restart

# Roadmap

# Components

## HUE
------
TODO

## DIO

DIO integration is based on this article
- http://www.homautomation.org/2013/10/09/how-to-control-di-o-devices-with-a-raspberry

## NFC-RC522

Usefull article for RASPBERRY PI 3 rev 1.2
- https://www.raspberrypi.org/forums/viewtopic.php?f=37&t=147291
