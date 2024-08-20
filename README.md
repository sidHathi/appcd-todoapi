### Siddharth Hathi

# AppCD TodoList API

This repo contains:
* An API for managing todo items and user todo lists written in Go
* Automated tests for the API
* IaC code written in terraform to deploy the API in a Kubernetes Cluster
* Code to generate OpenAPI docs for the API

## Usage and Deployment

* To run the API locally, start a postgres server and enter the following info into a `.env` file at the top level of the codebase:

```
DB_USER="Your postgres user's username"
DB_PASSWORD="Your postgres user's password"
DB_PORT="The port on which the database is running"
DB_NAME="The name of your database"
DB_HOST="The host location of your database (i.e. localhost)"
```
* To run the API in a local, dockerized Kubernetes cluster:
1. Install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/), [kubectl](https://kind.sigs.k8s.io/docs/user/quick-start/), [the terraform cli](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli), and [docker](https://docs.docker.com/engine/install/)
2. Create a local kind cluster using the config file found in the [terraform-kind](terraform-kind/) subdirectory of this repo. This can be done using the shell command `kind create cluster --config terraform-kind/kind-config.yaml`
3. Run the shell command `docker build -t sidhathi/appcd-todo .` to build the kubernetes docker image
4. Run `kind load docker-image sidhathi/appcd-todo:latest` to load the docker image into the kind context
5. Run the following commands:
```bash
cd terraform-kind
kubectl config use-context kind-kind
kubectl config config view --raw > kubeconfig
terraform init
terraform apply
```
* To provision a paid AWS EKS cluster and deploy the API on that cluster:
1. Install the aws cli, terraform, kubectl, and docker
2. Authenticate into your aws and hashicorp accounts through the command line
3. Run the following commands:
```bash
cd terraform-eks
terraform init
terraform apply
```

## Automated Testing

The tests for this API were written using Go's built in `testing` framework and the `httptest` request simulation framework. All tests for the repo are contained in the [tests](tests/) subdirectory. To run them use:
```bash
go test ./tests/...
```
For verbose output use:
```bash
go test -v ./tests/...
```
There are tests for each API endpoint - mainly to test that they function as expected and throw the appropriate errors for bad input. There are also integration tests for the API which test combinations of different endpoints to ensure that they work together as expected.

## OpenAPI Docs

OpenAPI docs for this API were generated using the `go-swag` library which parses in-code annotations above each API route controller. Docs are served on the `/docs/index.html` route when the API is run
