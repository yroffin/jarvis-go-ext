# jarvis-go-ext
Jarvis golang extention

[![Build Status](https://travis-ci.org/yroffin/jarvis-go-ext.svg?branch=master)](https://travis-ci.org/yroffin/jarvis-go-ext)

# Dependencies and setup

## MongoDB

Setup mongodb for store local data

        sudo apt-get install mongodb-server

## WiringPI

Connector is based on wiringPi so setup wiringPi on your host [here](http://wiringpi.com)

## Teleinfo option

For teleinfo option change right to 755 on USB file (ex: /dev/ttyUSB0) and apply this command in /etc/rc.local

        stty 1200 cs7 evenp cstopb -igncr -inlcr -brkint -icrnl -opost -isig -icanon -iexten -F /dev/ttyUSB0
    
## Setup on raspberry pi 2 or zero

    pi@raspberrypi:~ $ sudo useradd -m -b /home/jarvis jarvis
    pi@raspberrypi:~ $ export GITHUB=https://github.com/yroffin/jarvis-go-ext/releases/download/1.01b6
    pi@raspberrypi:~ $ sudo wget ${GITHUB}/jarvis-rest-module-0.0.1-SNAPSHOT.armel -O /home/jarvis/jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~ $ sudo chmod 755 /home/jarvis/jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~ $ sudo chown jarvis:jarvis /home/jarvis/jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~ $ sudo wget ${GITHUB}/jarvis-go-service -O /etc/init.d/jarvis-go-service
    pi@raspberrypi:~ $ sudo chmod 755 /etc/init.d/jarvis-go-service
    pi@raspberrypi:~ $ sudo update-rc.d jarvis-go-service defaults
    pi@raspberrypi:~ $ sudo service jarvis-go-service restart
    pi@raspberrypi:~ $ curl http://192.168.1.47:7000/api/teleinfo

# Roadmap
