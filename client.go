package roster

import (
	"time"
	"errors"
	"math/rand"
	"net"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	HeartbeatInterval = 300 * time.Millisecond
)

var (
	ErrServiceNotFound = errors.New("roster: No matching service found")
	ErrLocalIpNotFound = errors.New("roster: No non loopback local IP address could be found")
)

type Service struct {
	Name string
	Endpoint string
	Expiry int64
	stopHeatbeat chan bool
}

// Unregister the service
func (s *Service) Unregister() {
	s.stopHeatbeat <- true
}

type ClientConfiger interface {
	GetConfig() *aws.Config
	GetRegistryName() *string
	GetTTL() int64
}

type BaseConfig struct {
	RegistryName string
	TTL int64
}

func (bc BaseConfig) GetRegistryName() *string {
	if bc.RegistryName != "" {
		return aws.String(bc.RegistryName)
	} else {
		// Default
		return aws.String("dsd")
	}
}

// Returns the TTL in seconds
func (bc BaseConfig) GetTTL() int64 {
	if bc.TTL != 0 {
		return bc.TTL
	} else {
		// Default
		return 30
	}
}

type WebServiceConfig struct {
	BaseConfig
	Region string
}

func (wsc WebServiceConfig) GetConfig() *aws.Config {

	//Explicitley set
	region := wsc.Region

	// Environment var
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	// Default
	if region == "" {
		region = "us-west-2"
	}
	
	return aws.NewConfig().WithRegion(region)
}

type LocalConfig struct {
	BaseConfig
	Endpoint string
}

func (lc LocalConfig) GetConfig() *aws.Config {
	return aws.NewConfig().WithRegion("us-west-2").WithEndpoint(lc.Endpoint)
}

type Client struct {
	config ClientConfiger
	svc *dynamodb.DynamoDB
	Registry *Registry

}

func NewClient(config ClientConfiger) *Client {
	svc := dynamodb.New(session.New(), config.GetConfig())
	Registry := NewRegistry(svc,config.GetRegistryName())

	return &Client{svc: svc, config: config, Registry: Registry}
}

// Register the service in the registry
func (c *Client) Register(name string, endpoint string, serviceTTL ...int64) (*Service,error) {

	// Check whether the registry has been previously created. If not create before registration.
	if exists,err := c.Registry.Exists(); err != nil {
		return nil,err
	} else if !exists {
		if err := c.Registry.Create(); err != nil {
			return nil,err
		}
	}

	// If TTL not passed in then use default
	TTL := c.config.GetTTL()
	if len(serviceTTL) > 0 {
    	TTL = serviceTTL[0]
  	}

	// Create Service
	service := &Service{Name: name, Endpoint: endpoint, stopHeatbeat: make(chan bool)}

	// Heartbeat function - updates expiry
	heatbeat := func() {
		// Update service Expiry based on TTL and current time
		service.Expiry = time.Now().Unix() + TTL

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
	heatbeat()

	// Start goroutine to send heartbeat
	go func() {
	    for {
			select {
				case <- service.stopHeatbeat:
					return
		        default:
					// Pause for interval
					time.Sleep(HeartbeatInterval)

					// Call heartbeat function
					heatbeat()
		    }
	    }
	}()

	return service,nil
}

// Query the registry for named service
func (c *Client) Discover(name string) (*Service, error) {

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
