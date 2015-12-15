package roster

import (
	"testing"
    "time"
    "net/url"
)

// Create a client to local dynamodb
var c2 = NewClient(LocalConfig{Endpoint: "http://192.168.99.101:8000",BaseConfig: BaseConfig{RegistryName: "roster_client_test"}})

// Setup
func init() {
	DeleteTestRegistry(c2)
}

func TestRegister(t *testing.T) {

    localIp,err := c2.GetLocalIP()
    if err != nil {
        t.Error(err)
    }

    endpoint := &url.URL{Scheme: "http", Host: localIp + ":8889"}

    service,err := c2.Register("test-service",endpoint.String())
    if err != nil {
        t.Error(err)
    }

    if _,err := c2.Discover("test-service"); err != nil {
        t.Error(err)
    }

    service.Unregister()
}

func TestUnregister(t *testing.T) {

    localIp,err := c2.GetLocalIP()
    if err != nil {
        t.Error(err)
    }

    endpoint := &url.URL{Scheme: "http", Host: localIp + ":8889"}

    service,err := c2.Register("test-service",endpoint.String(),5)
    if err != nil {
        t.Error(err)
    }

    service.Unregister()

    // TTL set to 5 seconds, so wait 6 seconds to see if expired
    time.Sleep(6 * time.Second)

    if _,err := c2.Discover("test-service"); err == nil {
        t.Error("Service not unregistering")
    }

}
