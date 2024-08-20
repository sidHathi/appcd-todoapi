### Siddharth Hathi

# AppCD TodoList API

This repo contains:
* An API for managing todo items and user todo lists written in Go
* Automated tests for the API
* IaC code written in Terraform to deploy the API in a Kubernetes cluster
* Code to generate OpenAPI docs for the API

## Usage and Deployment

* **To run the API locally,** start a postgres server and enter the following info into a `.env` file at the top level of the codebase:

```
DB_USER="Your postgres user's username"
DB_PASSWORD="Your postgres user's password"
DB_PORT="The port on which the database is running"
DB_NAME="The name of your database"
DB_HOST="The host location of your database (i.e. localhost)"
```
Once the env file is added, you can use the `go get .` and `go run .` to fetch packages and run the API respectively.

* **To run the API in a local, dockerized Kubernetes cluster:**
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
* **To provision a paid AWS EKS cluster and deploy the API on that cluster:**
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

## API Functionality

This todo list API contains support for the following:
* Creating, modifying, deleting, and retrieving users. Users are not expected to authenticate since that's outside the scope of this task. However, resources owned by a user can only be accessed through that use. I.E., you need to provide a user's id to have get any lists associated with that user
* Creating, modifying, deleting, and retrieving todo-lists. Todo-lists are owned by users and constitute a certain number of items that a user has to complete
* Creating, modifying, deleting, and retrieving todo-items. Todo-items are items in a todo list. They represent a task that the user needs to complete and can have sub-items and attachments associated with them. Sub-items are child items that need to be completed in order for the parent to be complete. Attachments are references to file data that contain a url and a type
* Infinite nesting of items. Since an item can have sub-items, and those sub-items can have their own sub-items, items can theoretically nest infinitely
* Sharing todo lists between users. The todo-list/:listId/share endpoint can be used to give one user access to another's todo list

## API Implementation

The data underlying the todo-lists are stored in a postgres database hosted alongside the API. This schema diagram represent's the layout of the database's tables:

![database layout](/appcd-todo-db.png)

The controllers found in the [controllers](/controllers/) subdirectory contain code for handling user input and output through the API routes. Controllers are separated based on the type of resource they pertain to. The services, found in the [services](/services/) subdirectory contain code for updating and retrieving database information pertaining to each resource. The API is then deployed through the `main.go` class in the top level of the repo.
