# Roster: A library for simple service discovery using Dynamodb for Golang

Instead of having to manage a separate distributed key store like etcd or ZooKeeper, leverage AWS's own Dynamodb so you don't have to worry about hardware provisioning, setup and configuration, replication, software patching, or cluster scaling.

## Requirements

1. Go 1.5 or above
2. AWS credentials setup as per [AWS SDK documentation](https://github.com/aws/aws-sdk-go)

## Usage

### Using a local Dynamodb instance for development

The easiest way to get a local instance of Dynamodb running is to pull down [this](https://hub.docker.com/r/tutum/dynamodb/) docker image. Then the following steps (assumes Docker Machine is being used with default name):

1. `docker pull tutum/dynamodb`
2. `docker run -d -p 5000:5000 tutum/dynamodb`
3. 'export DYNAMODB_PORT=http://$(docker-machine ip default):5000'
4. Run your Go app.

## Contributing

### Installation

1. Clone the repository
2. Install [Glide](https://github.com/Masterminds/glide) which is used to handle vendoring and dependency management
3. Install dependencies `$ glide install`

### Running tests

`go test` will run tests. By default the tests will connect to the Amazon hosted Dynamodb in the default region. You can override the region by setting `AWS_REGION` environment variable as per the [AWS SDK](https://github.com/aws/aws-sdk-go). Also if the `DYNAMODB_PORT` environment variable is set, tests will be run against the local dynamodb instance.
