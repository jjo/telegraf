# MQTT Consumer Input Plugin

The [MQTT](http://mqtt.org/) consumer plugin reads from
specified MQTT topics and adds messages to InfluxDB.
The plugin expects messages in the
[Telegraf Input Data Formats](https://github.com/influxdata/telegraf/blob/master/DATA_FORMATS_INPUT.md).

### Configuration:

```toml
# Read metrics from MQTT topic(s)
[[inputs.mqtt_consumer]]
  servers = ["localhost:1883"]
  ### MQTT QoS, must be 0, 1, or 2
  qos = 0

  ### Topics to subscribe to
  topics = [
    "telegraf/host01/cpu",
    "telegraf/+/mem",
    "sensors/#",
  ]

  ### Maximum number of metrics to buffer between collection intervals
  metric_buffer = 100000

  ### username and password to connect MQTT server.
  # username = "telegraf"
  # password = "metricsmetricsmetricsmetrics"

  ### Optional SSL Config
  # ssl_ca = "/etc/telegraf/ca.pem"
  # ssl_cert = "/etc/telegraf/cert.pem"
  # ssl_key = "/etc/telegraf/key.pem"
  ### Use SSL but skip chain & host verification
  # insecure_skip_verify = false

  ### Data format to consume. This can be "json", "influx" or "graphite"
  ### Each data format has it's own unique set of configuration options, read
  ### more about them here:
  ### https://github.com/influxdata/telegraf/blob/master/DATA_FORMATS_INPUT.md
  data_format = "influx"
```

### Tags:

- All measurements are tagged with the incoming topic, ie
`topic=telegraf/host01/cpu`
