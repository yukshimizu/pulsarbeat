// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Config struct {
	Client   pulsarClientOptions   `config:"client"`
	Consumer pulsarConsumerOptions `config:"consumer"`
}

type pulsarClientOptions struct {
	URL                        string            `config:"url" validate:"required"`
	ConnectionTimeout          time.Duration     `config:"connection_timeout" validate:"min=0"`
	OperationTimeout           time.Duration     `config:"operation_timeout" validate:"min=0"`
	AuthenticationTLS          authenticationTLS `config:"authentication_tls"`
	AuthenticationAthenz       map[string]string `config:"authentication_athenz"`
	TLSTrustCertsFilePath      string            `config:"tls_trust_certs_file_path"`
	TLSAllowInsecureConnection bool              `config:"tls_allow_insecure_connection"`
	TLSValidateHostname        bool              `config:"tls_validate_hostname"`
	MaxConnectionsPerBroker    int               `config:"max_connections_per_broker"`
}

type authenticationTLS struct {
	CertificatePath string `config:"certificate_path"`
	PrivateKeyPath  string `config:"private_key_path"`
}

type pulsarConsumerOptions struct {
	Topic                       string            `config:"topic"`
	Topics                      []string          `config:"topics"`
	TopicsPattern               string            `config:"topics_pattern"`
	AutoDiscoveryPeriod         time.Duration     `config:"auto_discovery_period" validate:"min=0"`
	SubscriptionName            string            `config:"subscription_name" validate:"required"`
	Properties                  map[string]string `config:"properties"`
	Type                        string            `config:"subscription_type"`
	SubscriptionInitialPosition string            `config:"subscription_initial_position"`
	ReceiverQueueSize           int               `config:"receiver_queue_size"`
	NackRedeliveryDelay         time.Duration     `config:"nack_redelivery_delay" validate:"min=0"`
	Name                        string            `config:"name"`
	ReadCompacted               bool              `config:"read_compacted"`
	ReplicateSubscriptionState  bool              `config:"replicate_subscription_state"`
	NumWorkers                  int               `config:"num_workers" validate:"min=1"`
}

type authProvider int

const (
	authProviderNone authProvider = iota
	authProviderTLS
	authProviderAthenz
)

var DefaultConfig = Config{
	Client: pulsarClientOptions{
		URL:               "pulsar://localhost:6650",
		ConnectionTimeout: 20 * time.Second,
	},
	Consumer: pulsarConsumerOptions{
		Topic:            "my-topic",
		SubscriptionName: "my-sub",
		NumWorkers:       1,
	},
}

func (c *pulsarClientOptions) authValidate() (authProvider, error) {
	if len(c.AuthenticationAthenz) == 0 &&
		c.AuthenticationTLS.CertificatePath == "" && c.AuthenticationTLS.PrivateKeyPath == "" {
		return authProviderNone, nil
	}

	var ok bool
	if len(c.AuthenticationAthenz) != 0 {
		_, ok = c.AuthenticationAthenz["providerDomain"]
		if !ok {
			return authProviderAthenz, errors.New("Athenz providerDomain is not configured")
		}

		_, ok = c.AuthenticationAthenz["tenantDomain"]
		if !ok {
			return authProviderAthenz, errors.New("Athenz tenantDomain is not configured")
		}

		_, ok = c.AuthenticationAthenz["tenantService"]
		if !ok {
			return authProviderAthenz, errors.New("Athenz tenantService is not configured")
		}

		_, ok = c.AuthenticationAthenz["privateKey"]
		if !ok {
			return authProviderAthenz, errors.New("Athenz privateKey is not configured")
		}

		_, ok = c.AuthenticationAthenz["keyId"]
		if !ok {
			return authProviderAthenz, errors.New("Athenz keyId is not configured")
		}

		_, ok = c.AuthenticationAthenz["ztsUrl"]
		if !ok {
			return authProviderAthenz, errors.New("Athenz ztsUrl is not configured")
		}

		if c.AuthenticationTLS.CertificatePath != "" || c.AuthenticationTLS.PrivateKeyPath != "" {
			return authProviderNone, errors.New("Multiple authentication settings are configured")
		}
		return authProviderAthenz, nil
	} else {
		if c.AuthenticationTLS.CertificatePath == "" {
			return authProviderTLS, errors.New("TLS certificatePath is not configured")
		}

		if c.AuthenticationTLS.PrivateKeyPath == "" {
			return authProviderTLS, errors.New("TLS privateKeyPath is not configured")
		}
		return authProviderTLS, nil
	}
}

func NewPulsarClient(clientOptions pulsarClientOptions) (*pulsar.Client, error) {
	var clientConfig pulsar.ClientOptions
	clientConfig.URL = clientOptions.URL
	clientConfig.ConnectionTimeout = clientOptions.ConnectionTimeout
	clientConfig.OperationTimeout = clientOptions.OperationTimeout
	clientConfig.TLSTrustCertsFilePath = clientOptions.TLSTrustCertsFilePath
	clientConfig.TLSAllowInsecureConnection = clientOptions.TLSAllowInsecureConnection
	clientConfig.TLSValidateHostname = clientOptions.TLSValidateHostname
	clientConfig.MaxConnectionsPerBroker = clientOptions.MaxConnectionsPerBroker

	auth, err := clientOptions.authValidate()
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Authentication Settings")
	}
	switch auth {
	case authProviderAthenz:
		clientConfig.Authentication = pulsar.NewAuthenticationAthenz(clientOptions.AuthenticationAthenz)
	case authProviderTLS:
		clientConfig.Authentication = pulsar.NewAuthenticationTLS(clientOptions.AuthenticationTLS.CertificatePath,
			clientOptions.AuthenticationTLS.PrivateKeyPath)
	default:
	}

	client, err := pulsar.NewClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Initializing pulsar client")
	}

	return &client, nil
}

func NewPulsarConsumer(client *pulsar.Client, consumerOptions pulsarConsumerOptions) (*[]pulsar.Consumer, error) {
	var consumerConfig pulsar.ConsumerOptions
	consumerConfig.Topic = consumerOptions.Topic
	consumerConfig.AutoDiscoveryPeriod = consumerOptions.AutoDiscoveryPeriod
	consumerConfig.SubscriptionName = consumerOptions.SubscriptionName
	consumerConfig.Properties = consumerOptions.Properties

	switch consumerOptions.Type {
	case "Shared":
		consumerConfig.Type = pulsar.Shared
	case "Failover":
		consumerConfig.Type = pulsar.Failover
	case "KeyShared":
		consumerConfig.Type = pulsar.KeyShared
	default:
		consumerConfig.Type = pulsar.Exclusive
	}

	if consumerOptions.SubscriptionInitialPosition == "Earliest" {
		consumerConfig.SubscriptionInitialPosition = pulsar.SubscriptionPositionEarliest
	} else {
		consumerConfig.SubscriptionInitialPosition = pulsar.SubscriptionPositionLatest
	}

	consumerConfig.ReceiverQueueSize = consumerOptions.ReceiverQueueSize
	consumerConfig.NackRedeliveryDelay = consumerOptions.NackRedeliveryDelay
	consumerConfig.ReadCompacted = consumerOptions.ReadCompacted
	consumerConfig.ReplicateSubscriptionState = consumerOptions.ReplicateSubscriptionState

	var consumers []pulsar.Consumer
	for i := 1; i <= consumerOptions.NumWorkers; i++ {
		consumerConfig.Name = consumerOptions.Name + "-" + strconv.Itoa(i)
		consumer, err := (*client).Subscribe(consumerConfig)
		if err != nil {
			return nil, errors.Wrap(err, "Initializing pulsar consumer")
		}
		consumers = append(consumers, consumer)
	}
	return &consumers, nil
}
