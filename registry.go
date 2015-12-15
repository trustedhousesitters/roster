package roster

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var (
	ErrRegistryNotActive = errors.New("roster: Registry table has taken longer than expected to reach ACTIVE state")
)

type Registry struct {
	name *string
    svc *dynamodb.DynamoDB
}

// Creates a new registry
func NewRegistry(svc *dynamodb.DynamoDB, name *string) *Registry {
    return &Registry{svc: svc, name: name}
}

// Does the registry exist
func (r *Registry) Exists() (bool, error) {
    params := &dynamodb.DescribeTableInput{
		TableName: r.name,
	}
	_, err := r.svc.DescribeTable(params)

    if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			// Table not found so create
			if awsErr.Code() == "ResourceNotFoundException" {
				return false,nil
			}
	    }
		return false,err
	}

	return true,nil
}

// Is the registry in an active state
func (r *Registry) IsActive() (bool, error) {
    params := &dynamodb.DescribeTableInput{
		TableName: r.name,
	}
	resp, err := r.svc.DescribeTable(params)

    if err != nil {
		return false,err
	} else {
        return *resp.Table.TableStatus == dynamodb.TableStatusActive, nil
    }
}

// Create table with 2 attributes (Name and Expiry)
func (r *Registry) Create() error {
    _, err := r.svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String("Endpoint"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("Endpoint"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: r.name,
	})

	if err != nil {
		return err
	} else {
		// Table was created, but it's asynchronous...so block whilst it finishes being created
        attempts := 0
        for {
            if isActive,err := r.IsActive(); err == nil && isActive {
                break
            } else if err == nil {
                attempts++
            } else {
                return err
            }

            if attempts > 10 {
                return ErrRegistryNotActive
            }
        }

        return nil
	}
}

// Delete the registry
func (r *Registry) Delete() error {
    _, err := r.svc.DeleteTable(&dynamodb.DeleteTableInput{
    	TableName: r.name,
    })

    return err
}
