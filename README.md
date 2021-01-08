
# Pulsarbeat

Welcome to Pulsarbeat. Pulsaarbeat is an elastic [Beat](https://www.elastic.co/beats/) that reads messages from [Apache Pulsar](https://pulsar.apache.org/en/) topics and forwards them to [Logstash](https://www.elastic.co/logstash) (or any other configured output).

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/yukshimizu/pulsarbeat`

## Getting Started with Pulsarbeat

### Requirements

* [Golang](https://golang.org/dl/) >= 1.7
* [Pulsar Go client](https://github.com/apache/pulsar-client-go) v0.3.0


### Build

To build the binary for Pulsarbeat run the following commands. `make` will generate a binary
in the same directory with the name pulsarbeat.
```
mkdir -p ${GOPATH}/src/github.com/yukshimizu/pulsarbeat
git clone https://github.com/yukshimizu/pulsarbeat ${GOPATH}/src/github.com/yukshimizu/pulsarbeat
make
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Configure

An example configuration can be found in the file `pulsarbeat.yml`. The configuration adheres fundamentally to Pulsar clients and consumers configurations. Please refer to [Pulsar Go client](https://pulsar.apache.org/docs/en/client-libraries-go/) for more information. One additional parameter is num_workers, which specifies the number of workers receiving Pulsar messages.
```
pulsarbeat:
  # Configure pulsar client options.
  client:
    # Configure the service URL for the Pulsar service.
    # This parameter is required
    url: "pulsar://localhost:6650"
    # Timeout for the establishment of a TCP connection (default: 30 seconds).
    connection_timeout : 30s
    # Set the operation timeout (default: 30 seconds).
    # Producer-create, subscribe and unsubscribe operations will be retried until
    # this interval, after which the operation will be marked as failed.
    operation_timeout: 30s
    # Configure either TLS or Athenz authentication provider.
    #authentication_tls:
    #  certificate_path: "/path_to/admin.cert.pem"
    #  private_key_path: "/path_to/admin.key-pk8.pem"
    #authentication_athenz: {
    #  "providerDomain":"pulsar",
    #  "tenantDomain":"shopping",
    #  "tenantService":"some_app",
    #  "privateKey":"file:///path_to/some_app_private.pem",
    #  "keyId":"v0",
    #  "ztsUrl":"https://athenz.local:8443/zts/v1"
    #}
    # Set the path to the trusted TLS certificate file.
    #tls_trust_certs_file_path: "/path_to/ca.cert.pem"
    # Configure whether the Pulsar client accept untrusted TLS certificate from
    # broker (default: false).
    #tls_allow_insecure_connection: false
    # Configure whether the Pulsar client verify the validity of the host name from
    # broker (default: false).
    #tls_validate_hostname: false
    # Max number of connections to a single broker that will kept in the pool
    # (Default: 1 connection).
    max_connections_per_broker: 1

  # Configure pulsar consumer options.
  consumer:
    # Specify the topic this consumer will subscribe on.
    # Either a topic, a list of topics or a topics pattern are required when subscribing.
    topic: "my-topic"
    # Specify a list of topics this consumer will subscribe on.
    # Either a topic, a list of topics or a topics pattern are required when subscribing.
    #topics: ["my-topic"]
    # Specify a regular expression to subscribe to multiple topics under the same namespace.
    #topics_pattern:
    # Specify the interval in which to poll for new partitions or new topics
    # if using a TopicsPattern.
    #auto_discovery_period: 60s
    # Specify the subscription name for this consumer.
    # This argument is required when subscribing.
    subscription_name: "my-sub"
    # Attach a set of application defined properties to the consumer.
    # This properties will be visible in the topic stats.
    #properties: {"key", "value"}
    # Select the subscription type to be used when subscribing to the topic.
    # Default is `Exclusive`.
    subscription_type: "Exclusive"
    # InitialPosition at which the cursor will be set when subscribe.
    # Default is `Latest`.
    subscription_initial_position: "Latest"
    # Sets the size of the consumer receive queue.
    # The consumer receive queue controls how many messages can be accumulated
    # by the `Consumer` before the application calls `Consumer.receive()`.
    # Using a higher value could potentially increase the consumer throughput
    # at the expense of bigger memory utilization.
    # Default value is `1000` messages and should be good for most use cases.
    receiver_queue_size: 1000
    # The delay after which to redeliver the messages that failed to be processed.
    # Default is 1min (See `Consumer.Nack()`).
    nack_redelivery_delay: 60s
    # Set the consumer name.
    name: "my-consumer"
    # If enabled, the consumer will read messages from the compacted topic rather
    # than reading the full message backlog of the topic. This means that,
    # if the topic has been compacted, the consumer will only see the latest value for
    # each key in the topic, up until the point in the topic message backlog
    # that has been compacted. Beyond that point, the messages will be sent as normal.
    #
    # ReadCompacted can only be enabled subscriptions to persistent topics,
    # which have a single active consumer (i.e. failure or exclusive subscriptions).
    # Attempting to enable it on subscriptions to a non-persistent topics or on a
    # shared subscription, will lead to the subscription call throwing a PulsarClientException.
    read_compacted: false
    # Mark the subscription as replicated to keep it in sync across clusters.
    replicate_subscription_state: false
    # Number of go routine workers
    num_workers: 1
```

### Run

To run Pulsarbeat with debugging output enabled, run:

```
./pulsarbeat -c pulsarbeat.yml -e -d "*"
```

### Example

If the payload of pulsar message is "Hello-Pulsar", Pulsarbeat will emit the following event (Elastic output):
```
{
  "_index": "pulsarbeat-0.1.0-2021.01.07-000001",
  "_type": "_doc",
  "_id": "tVxp3HYB8llAnZpavSDf",
  "_version": 1,
  "_score": null,
  "_source": {
    "@timestamp": "2021-01-07T10:34:44.400Z",
    "pulsar": {
      "key": "",
      "timestamp": "2021-01-07T10:34:44.387Z",
      "topic": "persistent://public/default/my-topic",
      "producer": "standalone-0-0"
    },
    "message": "Hello-Pulsar",
    "ecs": {
      "version": "1.6.0"
    },
    ...
  }
}
```



### Test

To test Pulsarbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean Pulsarbeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
