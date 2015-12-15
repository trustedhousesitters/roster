package roster

import (
	"log"
	"testing"
)

// Create a client to local dynamodb
var c1 = NewClient(LocalConfig{Endpoint: "http://192.168.99.101:8000",BaseConfig: BaseConfig{RegistryName: "roster_registry_test"}})

// Setup
func init() {
	DeleteTestRegistry(c1)
}

func DeleteTestRegistry(client *Client) {
	// Check doesn't already exist (if it does delete)
	if exists,err := client.Registry.Exists(); err != nil {
		log.Println(err)
	} else if exists {
		if err:=client.Registry.Delete(); err != nil {
			log.Println(err)
		}
	}
}

// Create a registry
func TestCreate(t *testing.T) {
	if err := c1.Registry.Create(); err != nil {
		t.Error(err)
	}
}

// Check registry is active
func TestIsActive(t *testing.T) {
	if isActive,err := c1.Registry.IsActive(); err != nil  {
		t.Error(err)
	} else if !isActive {
		t.Error("Table is not in active state")
	}
}

// Delete registry
func TestDelete(t *testing.T) {
	if err:=c1.Registry.Delete(); err != nil {
		t.Error(err)
	}
}

// Check deleted registry is not active
func TestNotIsActive(t *testing.T) {
	if isActive,err := c1.Registry.IsActive(); err == nil || isActive  {
		t.Error("Table not deleted")
	}
}
