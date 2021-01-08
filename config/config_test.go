// +build !integration

package config

import (
	"github.com/elastic/beats/v7/libbeat/common"
	"testing"
	"time"
)

const (
	url                     = "pulsar://localhost:6650"
	urlTLS                  = "pulsar+ssl://localhost:6651"
	connTimeOut             = 20 * time.Second
	certificatePath         = "/tmp/my-ca/admin.cert.pem"
	privateKeyPath          = "/tmp/my-ca/admin.key-pk8.pem"
	certificatePathInvalid  = "/invalid_path/admin.cert.pem"
	privateKeyPathInvalid   = "/invalid_path/admin.key-pk8.pem"
	TLSTrustCertsFilePath   = "/tmp/my-ca/certs/ca.cert.pem"
	athenzProviderDomain    = "pulsar"
	athenzTenantDomain      = "shopping"
	athenzTenantService     = "some_app"
	athenzPrivateKey        = "file:///tmp/some_app/some_app_private.pem"
	athenzPrivateKeyInvalid = "file:///invalid_path/some_app_private.pem"
	athenzKeyId             = "v0"
	athenzZtsUrl            = "https://athenz.local:8443/zts/v1"
	topicName               = "my-topic"
	subscriptionName        = "my-sub"
)

func TestDefaultConfigLoad(t *testing.T) {
	cfg := common.NewConfig()
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		t.Fatalf("Error unpacking config: %v\n", err)
	}

	t.Logf("Default Config is: %+v\n", c)
}

func TestPulsarClientDefaultConfig(t *testing.T) {
	cfg := common.NewConfig()
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		t.Fatalf("Error unpacking config: %v\n", err)
	}

	t.Logf("Default Config is: %+v\n", c)

	client, err := NewPulsarClient(c.Client)
	if err != nil {
		t.Errorf("Could not instantiate Pulsar client: %v\n", err)
	}

	(*client).Close()
}

func TestPulsarClientAuthTLS(t *testing.T) {
	tests := []struct {
		name    string
		client  pulsarClientOptions
		wantErr bool
	}{
		{
			name: "Successful authentication settings",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationTLS: authenticationTLS{
					CertificatePath: certificatePath,
					PrivateKeyPath:  privateKeyPath,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: false,
		},
		{
			name: "No CertificationPath error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationTLS: authenticationTLS{
					PrivateKeyPath: privateKeyPath,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "No PrivateKeyPath error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationTLS: authenticationTLS{
					CertificatePath: certificatePath,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "Invalid CertificationPath error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationTLS: authenticationTLS{
					CertificatePath: certificatePathInvalid,
					PrivateKeyPath:  privateKeyPath,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "Invalid PrivateKeyPath error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationTLS: authenticationTLS{
					CertificatePath: certificatePath,
					PrivateKeyPath:  privateKeyPathInvalid,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "Multiple authentication configurations error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationTLS: authenticationTLS{
					CertificatePath: certificatePath,
					PrivateKeyPath:  privateKeyPath,
				},
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantDomain":   athenzTenantDomain,
					"tenantService":  athenzTenantService,
					"privateKey":     athenzPrivateKey,
					"keyId":          athenzKeyId,
					"ztsUrl":         athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Client config is: %+v\n", test.client)
			client, err := NewPulsarClient(test.client)
			if test.wantErr {
				if err == nil {
					t.Error("Supposed to have err, but actually no err")
				} else {
					t.Logf("Could not instantiate Pulsar client: %v\n", err)
				}
			} else {
				if err != nil {
					t.Errorf("Could not instantiate Pulsar client: %v\n", err)
				} else {
					(*client).Close()
				}
			}
		})
	}
}

func TestPulsarClientAuthAthenz(t *testing.T) {
	tests := []struct {
		name    string
		client  pulsarClientOptions
		wantErr bool
	}{
		{
			name: "Successful authentication settings",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantDomain":   athenzTenantDomain,
					"tenantService":  athenzTenantService,
					"privateKey":     athenzPrivateKey,
					"keyId":          athenzKeyId,
					"ztsUrl":         athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: false,
		},
		{
			name: "No providerDomain error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"tenantDomain":  athenzTenantDomain,
					"tenantService": athenzTenantService,
					"privateKey":    athenzPrivateKey,
					"keyId":         athenzKeyId,
					"ztsUrl":        athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "No tenantDomain error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantService":  athenzTenantService,
					"privateKey":     athenzPrivateKey,
					"keyId":          athenzKeyId,
					"ztsUrl":         athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "No tenantService error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantDomain":   athenzTenantDomain,
					"privateKey":     athenzPrivateKey,
					"keyId":          athenzKeyId,
					"ztsUrl":         athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "No privateKey error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantDomain":   athenzTenantDomain,
					"tenantService":  athenzTenantService,
					"keyId":          athenzKeyId,
					"ztsUrl":         athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "No keyId error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantDomain":   athenzTenantDomain,
					"tenantService":  athenzTenantService,
					"privateKey":     athenzPrivateKey,
					"ztsUrl":         athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "No ztsUrl error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantDomain":   athenzTenantDomain,
					"tenantService":  athenzTenantService,
					"privateKey":     athenzPrivateKey,
					"keyId":          athenzKeyId,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
		{
			name: "Invalid privateKey path error",
			client: pulsarClientOptions{
				URL:               urlTLS,
				ConnectionTimeout: connTimeOut,
				AuthenticationAthenz: map[string]string{
					"providerDomain": athenzProviderDomain,
					"tenantDomain":   athenzTenantDomain,
					"tenantService":  athenzTenantService,
					"privateKey":     athenzPrivateKeyInvalid,
					"keyId":          athenzKeyId,
					"ztsUrl":         athenzZtsUrl,
				},
				TLSTrustCertsFilePath:      TLSTrustCertsFilePath,
				TLSAllowInsecureConnection: true,
				TLSValidateHostname:        false,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Client config is: %+v\n", test.client)
			client, err := NewPulsarClient(test.client)
			if test.wantErr {
				if err == nil {
					t.Error("Supposed to have err, but actually no err")
				} else {
					t.Logf("Could not instantiate Pulsar client: %v\n", err)
				}
			} else {
				if err != nil {
					t.Errorf("Could not instantiate Pulsar client: %v\n", err)
				} else {
					(*client).Close()
				}
			}
		})
	}
}

/*
The following tests are commented out because they require specific pulsar environment respectively to communicate with.
You can use those tests if required.

func TestPulsarConsumerDefaultConfig(t *testing.T) {
	// Requires pulsar instance listening with port 6650
	cfg := common.NewConfig()
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		t.Fatalf("Error unpacking config: %v\n", err)
	}

	t.Logf("Default Config is: %+v\n", c)

	client, err := NewPulsarClient(c.Client)
	if err != nil {
		t.Errorf("Could not instantiate Pulsar client: %v\n", err)
	}

	defer (*client).Close()

	consumers, err := NewPulsarConsumer(client, c.Consumer)
	if err != nil {
		t.Errorf("Could not instantiate Pulsar conssumer: %v\n", err)
	}

	for _, consumer := range *consumers {
		consumer.Close()
	}
}

var client = pulsarClientOptions{
	URL:               url,
	ConnectionTimeout: connTimeOut,
}

func TestPulsarConsumerNoAuth(t *testing.T) {
	// Requires pulsar instance listening with port 6650
	tests := []struct {
		name     string
		client   pulsarClientOptions
		consumer pulsarConsumerOptions
		wantErr  bool
	}{
		{
			name:   "SubscriptionType Shared",
			client: client,
			consumer: pulsarConsumerOptions{
				Topic:            topicName,
				SubscriptionName: subscriptionName,
				Type:             "Shared",
				NumWorkers:       2,
			},
			wantErr: false,
		},
		{
			name:   "SubscriptionType Failover",
			client: client,
			consumer: pulsarConsumerOptions{
				Topic:            topicName,
				SubscriptionName: subscriptionName,
				Type:             "Failover",
				NumWorkers:       2,
			},
			wantErr: false,
		},
		{
			name:   "SubscriptionType KeyShared",
			client: client,
			consumer: pulsarConsumerOptions{
				Topic:            topicName,
				SubscriptionName: subscriptionName,
				Type:             "KeyShared",
				NumWorkers:       2,
			},
			wantErr: false,
		},
		{
			name:   "SubscriptionType Exclusive cannot instantiate multiple consumers",
			client: client,
			consumer: pulsarConsumerOptions{
				Topic:            topicName,
				SubscriptionName: subscriptionName,
				Type:             "Exclusive",
				NumWorkers:       2,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Client config is: %+v\n", test.client)
			client, err := NewPulsarClient(test.client)
			if err != nil {
				t.Errorf("Could not instantiate Pulsar client: %v\n", err)
			}

			defer (*client).Close()

			t.Logf("Consumer config is: %+v\n", test.consumer)
			consumers, err := NewPulsarConsumer(client, test.consumer)

			if test.wantErr {
				if err == nil {
					t.Error("Supposed to have err, but actually no err")
				} else {
					t.Logf("Could not instantiate Pulsar consumer: %v\n", err)
				}
			} else {
				if err != nil {
					t.Errorf("Could not instantiate Pulsar consumer: %v\n", err)
				} else {
					for _, consumer := range *consumers {
						consumer.Close()
					}
				}
			}
		})
	}
}

func TestPulsarConsumerAuthTLS(t *testing.T) {
	// Requires pulsar instance listening with port 6651
	cfg := common.NewConfig()
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		t.Fatalf("Error unpacking config: %v\n", err)
	}
	c.Client.URL = urlTLS
	c.Client.AuthenticationTLS.CertificatePath = certificatePath
	c.Client.AuthenticationTLS.PrivateKeyPath = privateKeyPath
	c.Client.TLSTrustCertsFilePath = TLSTrustCertsFilePath
	c.Client.TLSAllowInsecureConnection = true
	c.Client.TLSValidateHostname = false

	t.Logf("Default Config is: %+v\n", c)

	client, err := NewPulsarClient(c.Client)
	if err != nil {
		t.Errorf("Could not instantiate Pulsar client: %v\n", err)
	}

	defer (*client).Close()

	consumers, err := NewPulsarConsumer(client, c.Consumer)
	if err != nil {
		t.Errorf("Could not instantiate Pulsar conssumer: %v\n", err)
	}

	for _, consumer := range *consumers {
		consumer.Close()
	}
}

*/
