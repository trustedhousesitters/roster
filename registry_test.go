package roster

import (
	"log"
	"testing"
)

// Setup
func init() {
	client := NewConfiguredClient("roster_registry_test")
	DeleteTestRegistry(client)
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
	client := NewConfiguredClient("roster_registry_test")

	if err := client.Registry.Create(); err != nil {
		t.Error(err)
	}
}

// Check registry is active
func TestIsActive(t *testing.T) {
	client := NewConfiguredClient("roster_registry_test")

	if isActive,err := client.Registry.IsActive(); err != nil  {
		t.Error(err)
	} else if !isActive {
		t.Error("Table is not in active state")
	}
}

// Delete registry
func TestDelete(t *testing.T) {
	client := NewConfiguredClient("roster_registry_test")

	if err:=client.Registry.Delete(); err != nil {
		t.Error(err)
	}
}

// Check deleted registry is not active
func TestNotIsActive(t *testing.T) {
	client := NewConfiguredClient("roster_registry_test")

	if isActive,err := client.Registry.IsActive(); err == nil || isActive  {
		t.Error("Table not deleted")
	}
}
