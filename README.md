# Diver

This is a tool to interact with the APIs of the Docker Enterprise Edition products enabling an end user to provision, manage and monitor the platform. 

## Building

**QUICK USAGE** From your local machine `go get github.com/thebsdbox/diver`

**Alternatively (or for dev work)**
Clone the repository and build with `make build`, the `make docker` will build a local `scratch `container that only has the binary.

Alternatively you can manually compile `diver` through the use of `go build`.

## Usage

The main commands are `dtr`, `ucp` and `store` that interact directly with those areas of the EE platform or the Docker store.

```
./diver -h
This tool uses the native APIs to "dive" into Docker EE

Usage:
  diver [command]

Available Commands:
  dtr         Docker Trusted Registry
  help        Help about any command
  store       Docker Store
  ucp         Universal Control Plane 

Flags:
  -h, --help   help for diver

Use "diver [command] --help" for more information about a command.
```

## *STATUS*

**UCP**
- Login
- Query/Create/Delete Users
- Query/Create/Delete Organisations
- Query/Create/Delete Teams
- Get/Set Swarm
- Clone and set Roles
- Set Grants from subject, role, object
- Build Images from a local or Remote (github URL) Dockerfile
- Download client bundle
- Inspect Services (Endpoints coming soon)
- Manage Swarm configuration

.. More coming

**DTR**
- Login
- List Replicas and health

.. More coming

**STORE**
- Login
- Retrieve Subscriptions
- Find recent active
- Download licenses

## UCP

### Logging in to UCP

```
./diver ucp login --username docker               \
                  --password password             \
                  --url https://docker01.fnnrn.me \
                  --ignorecert
INFO[0000] Succesfully logged into [https://docker01.fnnrn.me] 
```

### Creating users/organisations

This uses the `auth` command as part of `ucp`.

This will create a new **user** called `bob`, to create an organisation use the `-isorg` flag. The `--action` flag identifies what operation will take place, such as `create`, `delete` and `modify`.

```
./diver ucp auth --active               \
                 --admin                \
                 --fullname "Bob Smith" \
                 --username bob         \
                 --password chess123    \
                 --action create
```


### Working with Roles

Once logged in you can list/get and create roles as per the example below:

```
dan $ ./diver ucp auth roles list | grep jenkins
998612c1-b367-42af-9d82-b2a5de9f8851    false   jenkins

dan $ ./diver ucp auth roles get --rolename jenkins > jenkins.role

dan $ ./diver ucp auth roles create --rolename jenkins2 --ruleset ./jenkins.role
INFO[0000] Role [jenkins2] created succesfully

dan $ ./diver ucp auth roles list | grep jenkins
260976b1-76d7-4ef0-84e2-6ae6b896eed1    false   jenkins2
998612c1-b367-42af-9d82-b2a5de9f8851    false   jenkins
```

### Working with Grants

**Listing**
To list all current grants you can use the following command:
`./diver ucp auth grants list`

To **resolve** the gran UUID to an actual `name` use the `--resolve` flag when listing grants.

**Creating**

To create a grant use the command `./diver ucp auth grants set` with the following flags:

`--collection` - Can either be a collection path or a Kubernetes namespace.

`--subject` - A user or service account.

`--role` - A role that has been created in UCP

`--type` - The type of grant that will be applied to, can be a `collection` grant, a single `namespace` grant or `all` kube namespaces.

**NOTE**: Unless the accounts are pre-configured UCP accounts then the UUIDs will need to be passed to this command.

#### EXAMPLE - Deploying HELM

**Before Installing Helm**

Create a Kubernetes service account:
`kubectl create serviceaccount --namespace kube-system tiller`

Create a grant for the `tiller` service account:

```
  ./diver ucp auth grants set --role fullcontrol       \
  --subject system:serviceaccount:kube-system:tiller  \
  --collection kubernetesnamespaces                    \
  --type all
```

Install (or init) Helm

`helm init`

Correct the service account

`kubectl patch deploy --namespace kube-system tiller-deploy -p ‘{“spec”:{“template”:{“spec”:{“serviceAccount”:”tiller”}}}}’`

Deploy using Helm! 

e.g MySQL deployment.

`helm install --name mysql stable/mysql`

### Downloading the client bundle

Download the client bundle to your local machine.

```
./diver ucp client-bundle
INFO[0000] Downloading the UCP Client Bundle            
```


### Watching Containers

This will present a colour coded output on memory usage of all containers that are running in a swarm cluster.. (using [urchin](http://github.com/thebsdbox/urchin) to hit memory reservations in the demo below)


```
./diver ucp containers top
```

![](img/container-top.jpg)

## DTR

### Logging into DTR

```
./diver dtr login --username docker               \
                  --password password             \
                  --url https://docker02.fnnrn.me \
                  --ignorecert
INFO[0000] Succesfully logged into [https://docker02.fnnrn.me] 
```


### Repositories


### List Replicas

```
./diver dtr info replicas

Replica         Status
a3a8ab213a8b     OK
ecb7a768afc4     OK
```

## Docker Store

### Interacting with Docker Store

Logging into the Docker Store through the following command:

`./diver store --username <user> --password <password>`


To retrieve the Docker Store User ID for this user use the following command:

`./diver store user`

The `ID` field is used for retrieving subscriptions and licenses, other users can be examined by using the `--user <username>` flag.

Retrieve subscriptions for this user with the following command:

`./diver store subscriptions ls --id <ID>`

To retrieve the first active subscription use the `--firstactive` flag.

To retrieve a subscription license use the following command:

`./diver store licenses get --subscription <SUBSCRIPTION>`

This will print the raw output so it is advisable to pipe this to a file with the following addition to the command:

`> ./subscription_ID.lic`

## Debugging Issues

When errors are reported turn up the `--logLevel` to 5, which enables debugging output.
