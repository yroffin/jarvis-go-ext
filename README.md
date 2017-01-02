# jarvis-go-ext
Jarvis golang extention

# Dependencies

Setup wiringPi on your host

# Setup on raspberry pi 2 or zero

    pi@raspberrypi:~ $ mkdir jarvis
    pi@raspberrypi:~ $ cd jarvis
    pi@raspberrypi:~/jarvis $ wget https://github.com/yroffin/jarvis-go-ext/releases/download/1.01b/jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ ls -lrt
    total 11784
    -rw-r--r-- 1 pi pi 12063744 Jan  2 10:07 jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ chmod 755 jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ ./jarvis-rest-module-0.0.1-SNAPSHOT.armel
    pi@raspberrypi:~/jarvis $ ./jarvis-rest-module-0.0.1-SNAPSHOT.armel start --jarvis.option.teleinfo true

