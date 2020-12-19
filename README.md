# goddl - go deCONZ data logger

This is a simple data logger, which fetches sensor data from the [deCONZ REST Api](https://dresden-elektronik.github.io/deconz-rest-doc/), and stores it to a csv file.

Obviously, you need a [deCONZ](https://www.dresden-elektronik.com/wireless/software/deconz.html) gateway with a ConBee/RaspBee dongle and some connected sensors.

I started writing this, since I found more generic home automation software like openHAB way to powerful (not to say bloated) for the simple task of data logging, and for the sole purpose of hacking something in Go.

Use it like so:

    ./goddl --ip $YOUR_GATEWAY_IP --storeconfig

When no API key is present (on the first run usually there isn't), it will try to register with the gateway. Press the link button in the deCONZ software before running the binary.

After a key was registered, it will poll all sensors from the gateway and continuously log their values in a csv logfile. With the `--storeconfig` option the new API key will be written to the config file right away, so you can omit the commandline arguments on subsequent runs.

Simple as that.

So far, this small program only logs temperature and humidity. More to come. Eventually.
