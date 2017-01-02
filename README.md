# jarvis-go-ext
Jarvis golang extention [![Build Status](https://travis-ci.org/yroffin/jarvis-go-ext.svg?branch=master)](https://travis-ci.org/yroffin/jarvis-go-ext)

# Dependencies and setup

## MongoDB

Setup mongodb for store local data

        sudo apt-get install mongodb-server

## WiringPI

Tools is base on wiringPi so setup wiringPi on your host

## Teleinfo option

For teleinfo option change right to 755 on USB file (ex: /dev/ttyUSB0) and apply this command in /etc/rc.local

        stty 1200 cs7 evenp cstopb -igncr -inlcr -brkint -icrnl -opost -isig -icanon -iexten -F /dev/ttyUSB0
    
## Setup on raspberry pi 2 or zero

    pi@raspberrypi:~ $ mkdir jarvis
    pi@raspberrypi:~ $ cd jarvis
    pi@raspberrypi:~ $ export GITHUB=https://github.com/yroffin/jarvis-go-ext/releases/download/1.01b
    pi@raspberrypi:~/jarvis $ wget ${GITHUB}/jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ ls -lrt
    total 11784
    -rw-r--r-- 1 pi pi 12063744 Jan  2 10:07 jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ chmod 755 jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ ./jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ ./jarvis-rest-module-0.0.1-SNAPSHOT.armel start --jarvis.option.teleinfo true
    pi@raspberrypi:~/jarvis $ curl http://192.168.1.47:7000/api/teleinfo

