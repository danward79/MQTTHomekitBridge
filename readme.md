## MQTTHomekitBridge
Bridge devices which communicate by MQTT to standard HomeKit devices, such that they then become available to use on iOS devices.

### Use case
Various devices communicate in my home using simple MQTT topic and payload messages, which are managed using [Mosquitto](https://mosquitto.org) as the broker. Messages are of the form below, but could change at any moment to a different form:

```
home/bedroom/temp = 29.18
home/bedroom/light = 14

home/balcony/light = 49
home/balcony/humi = 1.00
home/balcony/temp = 26.90
home/balcony/battery = 3160

home/lounge/pressure = 1004.38
home/lounge/battery = 3260
home/lounge/light = 90
home/lounge/temp = 26.60
```

It would be nice to have a method of linking arbitrary message data to HomeKit accessory devices, via a configuration file of some form. Flexibility in message topic and payloads would also be useful to minimise the cases where a change of device impacts the ability to bridge.

### Config
A simple config file should allow mapping the devices to be bridged. TOML may be one method of encoding these details.

``` TOML
# Configuration of MQTT devices to be bridged

#pin = "10340567"
broker = "192.168.1.22:1883"

[[devices.temperaturesensor]]
displayname = "Balcony Temperature"
topic = "home/balcony/temp"

[[devices.temperaturesensor]]
displayname = "Lounge Temperature"
topic = "home/lounge/temp"

[[devices.lightsensor]]
displayname = "Balcony Light"
topic = "home/balcony/light"

[[devices.lightsensor]]
displayname = "Lounge Light"
topic = "home/lounge/light"
```
