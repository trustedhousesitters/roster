package roster

import (
	"time"
	"errors"
	"math/rand"
	"net"
	"os"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	HeartbeatInterval = 1000 * time.Millisecond
	ServiceTTL = 5
)

var (
	ErrServiceNotFound = errors.New("roster: No matching service found")
	ErrLocalIpNotFound = errors.New("roster: No non loopback local IP address could be found")
)

type Service struct {
	Name string
	Endpoint string
	Expiry int64
	stopHeartbeat chan bool
}

// Unregister the service
func (s *Service) Unregister() {
	s.stopHeartbeat <- true

	// Pause (block) for connections to drain
	time.Sleep(ServiceTTL * time.Second)
}

type ClientConfig struct {
	RegistryName string
	Region string
	Endpoint string
}

func (cc ClientConfig) GetRegistryName() *string {
	if cc.RegistryName != "" {
		return aws.String(cc.RegistryName)
	} else {
		// Default
		return aws.String("roster")
	}
}

func (cc ClientConfig) GetConfig() *aws.Config {

	//Explicitley set
	region := cc.Region

	// Environment var
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	// Default
	if region == "" {
		region = "us-west-2"
	}

	// Create AWS config
	config := aws.NewConfig().WithRegion(region)

	// Check to see if DynamoDB is running locally
	endpoint := cc.Endpoint
	if endpoint == "" {
		// Env variable must be this format for Docker Compose to work, and also must replace tcp with http
		endpoint = os.Getenv("DYNAMODB_PORT")
		endpoint = strings.Replace(endpoint, "tcp", "http", 1)
	}

	if endpoint == "" {
		return config
	} else {
		return config.WithEndpoint(endpoint)
	}

}

type Client struct {
	config ClientConfig
	svc *dynamodb.DynamoDB
	Registry *Registry

}

func NewClient(config ClientConfig) *Client {
	svc := dynamodb.New(session.New(), config.GetConfig())
	Registry := NewRegistry(svc,config.GetRegistryName())

	return &Client{svc: svc, config: config, Registry: Registry}
}

// Register the service in the registry
func (c *Client) Register(name string, endpoint string) (*Service,error) {

	// Check whether the registry has been previously created. If not create before registration.
	if exists,err := c.Registry.Exists(); err != nil {
		return nil,err
	} else if !exists {
		if err := c.Registry.Create(); err != nil {
			return nil,err
		}
	}

	// Create Service
	service := &Service{Name: name, Endpoint: endpoint, stopHeartbeat: make(chan bool)}

	// Heartbeat function - updates expiry
	heartbeat := func() {
		// Update service Expiry based on TTL and current time
		service.Expiry = time.Now().Unix() + ServiceTTL

		// Update service entry in registry
		if av, err := dynamodbattribute.ConvertToMap(*service); err != nil {
			return
		} else {
			_, err := c.svc.PutItem(&dynamodb.PutItemInput{
				Item: av,
				TableName: c.config.GetRegistryName(),
			})

			if err != nil {
				return
			}
		}
	}

	// Ensure call heartbeat at least once
	heartbeat()

	// Start goroutine to send heartbeat
	go func() {
	    for {
			select {
				case <- service.stopHeartbeat:
					return
		        default:
					// Pause for interval
					time.Sleep(HeartbeatInterval)

					// Call heartbeat function
					heartbeat()
		    }
	    }
	}()

	return service,nil
}

// Query the registry for named service
func (c *Client) Discover(name string) (*Service, error) {

	// Make sure registry is active
	if c.Registry.IsActive() {
		expressionAttributeValues := map[string]interface{}{
	    	":NameVal": name,
	    	":ExpiryVal":   time.Now().Unix(),
		}

		ean := map[string]*string{
	    	"#N": aws.String("Name"),
		}

		eav, err := dynamodbattribute.ConvertToMap(expressionAttributeValues)
		if err != nil {
			return nil,err
		}

		resp, err := c.svc.Query(&dynamodb.QueryInput{
			TableName: c.config.GetRegistryName(),
			KeyConditionExpression: aws.String("#N = :NameVal"),
			FilterExpression: aws.String("Expiry > :ExpiryVal"),
			ExpressionAttributeValues: eav,
			ExpressionAttributeNames: ean,
		})

		if err != nil {
			return nil, err
		}

		if len(resp.Items) > 0 {
			// Randomly select one of the available endpoints (in effect load balancing between available endpoints)
			service := Service{}
			err = dynamodbattribute.ConvertFromMap(resp.Items[rand.Intn(len(resp.Items))], &service)
			if err != nil {
				return nil, err
			}
			return &service, nil
		} else {
			// No service found
			return nil,ErrServiceNotFound
		}
	} else {
		return nil,ErrRegistryNotActive
	}
}

// Returns the non loopback local IP of the host the client is running on
func (c *Client) GetLocalIP() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "",err
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String(),nil
            }
        }
    }
    return "",ErrLocalIpNotFound
}
