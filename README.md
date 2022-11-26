# go-mqtt-handler-demo

### example application:
#### demo app using Go language and Paho MQTT library for managing messeges from IoT devices (mostly zigbee but also Tasmota based devices and Arduino like microcontrollers - ESP32, ESP8266 (pure mqtt))

&nbsp;

### an application works with several other components, most running as containers (this repo contains only go app):
- Redis - for storing state and caching temporary data
- Mosquitto - MQTT broker - core component.
- Zigbee2Mqtt - passing messages from zigbee devices onto MQTT broker
- Grafana for dashboards, visualizations and alerts
- Telegram Bot API for push notifications on mobile devices
- Wireguard - VPN tunnel for accessing our stack from external networks
- custom REST API for an easy access to Redis db and for managing containers/images
- Mqtt Dashboard - IoT and Node-RED controller - native client app for Android

&nbsp;

ip and mac addresses are anonymized in this project

&nbsp;

### references:
- https://github.com/eclipse/paho.mqtt.golang
- https://tasmota.github.io/docs/
- https://www.zigbee2mqtt.io/
- https://mosquitto.org/


