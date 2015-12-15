# Roster: A library for simple service discovery using Dynamodb for Golang

Instead of having to manage a separate distributed key store like etcd or ZooKeeper, leverage AWS's own Dynamodb so you don't have to worry about hardware provisioning, setup and configuration, replication, software patching, or cluster scaling.


## Contributing

### Installation

1. Clone the repository
2. Install [Glide](https://github.com/Masterminds/glide) which is used to handle vendoring and dependency management
3. Install dependencies `$ glide install`

### Running tests

`go test ...` will run tests on vendor libraries as well. Instead run tests using `$ go test $(glide novendor)`
