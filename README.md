# Roster: A library for simple service discovery using Dynamodb for Golang

Instead of having to manage a separate distributed key store like etcd or ZooKeeper, leverage AWS's own Dynamodb so you don't have to worry about hardware provisioning, setup and configuration, replication, software patching, or cluster scaling.

As the only dependency is Dynamodb, it has the added benefit of allowing service discovery between applications running on EC2 and dockerized applications running within ECS.

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

### Examples

Within the examples folder there is an example echo client and server that demonstrate basic registry and discovery. Ensuring you have Docker Compose installed, simply run `docker-compose up`. This creates 1 client and 2 instances of the server service - when one of these servers is stopped, the client should remain operational and shouldn't drop any requests.

## Contributing

### Installation

1. Clone the repository
2. Install [Glide](https://github.com/Masterminds/glide) which is used to handle vendoring and dependency management
3. Install dependencies `$ glide install`

### Running tests

`go test .` will run tests. By default the tests will connect to the Amazon hosted Dynamodb in the default region. You can override the region by setting `AWS_REGION` environment variable as per the [AWS SDK](https://github.com/aws/aws-sdk-go). Also if the `DYNAMODB_PORT` environment variable is set, tests will be run against the local dynamodb instance.
