package roster

import (
	"testing"
    "time"
    "net/url"
)

// Setup
func init() {
	client := NewConfiguredClient("roster_client_test")
	DeleteTestRegistry(client)
}

func NewConfiguredClient(registryName string) *Client {
	return NewClient(ClientConfig{RegistryName: registryName})
}

func TestRegister(t *testing.T) {
	client := NewConfiguredClient("roster_client_test")

    localIp,err := client.GetLocalIP()
    if err != nil {
        t.Error(err)
    }

    endpoint := &url.URL{Scheme: "http", Host: localIp + ":8889"}

    service,err := client.Register("test-service",endpoint.String())
    if err != nil {
        t.Error(err)
    }

    if _,err := client.Discover("test-service"); err != nil {
        t.Error(err)
    }

    service.Unregister()
}

func TestUnregister(t *testing.T) {
	client := NewConfiguredClient("roster_client_test")

    localIp,err := client.GetLocalIP()
    if err != nil {
        t.Error(err)
    }

    endpoint := &url.URL{Scheme: "http", Host: localIp + ":8889"}

    service,err := client.Register("test-service",endpoint.String())
    if err != nil {
        t.Error(err)
    }

    service.Unregister()

    // TTL is set to 5 seconds, so wait 6 seconds to see if expired
    time.Sleep(6 * time.Second)

    if _,err := client.Discover("test-service"); err == nil {
        t.Error("Service not unregistering")
    }

}
