# pakhshi

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/snapp-incubator/pakhshi/lint.yaml?label=lint&logo=github&style=flat-square&branch=main)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/snapp-incubator/pakhshi/test.yaml?label=test&logo=github&style=flat-square&branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/snapp-incubator/pakhshi.svg)](https://pkg.go.dev/github.com/snapp-incubator/pakhshi)
[![Codecov](https://img.shields.io/codecov/c/gh/snapp-incubator/pakhshi?logo=codecov&style=flat-square)](https://codecov.io/gh/snapp-incubator/pakhshi)

## Introduction

Consider you have an array of brokers but you want to publish and subscribe on all of them at the same time.
Why you may need this setup? consider clients are randomly distributed between available clusters and you don't want to check which client is connected to which
broker so you will publish on all cluster and your client is connected to one them.

## How?

This library use [paho](https://github.com/eclipse/paho.mqtt.golang) in the background so you can easily change your applications to use this instead of paho.
It trying to implement all paho interfaces.

## Examples

The following example shows how to subscribe on the same topic on two brokers.

```go
opts := mqtt.NewClientOptions()
opts.AddBroker("tcp://127.0.0.1:1883")
opts.AddBroker("tcp://127.0.0.1:1884")

c := client.NewClient(opts)

if token := c.Connect(); token.Wait() && token.Error() != nil {
  assert.NoError(t, token.Error())
}

if token := c.Subscribe("hello", 0, func(c mqtt.Client, m mqtt.Message) {
  ch <- string(m.Payload())
}); token.Wait() && token.Error() != nil {
  assert.NoError(t, token.Error())
}
```

## Credits

Based on idea of [Ahmad Anvari](https://github.com/anvari1313).
