# pakhshi

## Introduction
Consider you have an array of brokers but you want to publish and subscribe on all of them at the same time.
Why you may need this setup? consider clients randomly distributed between available clusters and you don't want to check which client is connected to which
broker so you will publish on all cluster and your client is connected to one them.

## How?
This library use [paho](https://github.com/eclipse/paho.mqtt.golang) in the background so you can easily change your applications to use this instead of paho.
It trying to implement all paho interfaces.

## Credits
Based on idea of [Ahmad Anvari](https://github.com/anvari1313).
