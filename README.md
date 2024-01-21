# Inariam

Inariam is a cutting-edge software solution designed to streamline the management of users, groups, roles, and permissions across multi-cloud environments. With the integration of leading cloud providers – Azure, Google Cloud Platform (GCP), and Amazon Web Services (AWS) – Inariam empowers organizations to efficiently manage access controls, enforce security measures, and ensure compliance with industry standards.

## Development setup

```shell
git clone https://gitea.cloudpcp.com/Inariam/Inariam
cd Inariam
mkdir $HOME/.inariam
cp config.example.yaml /$HOME/.inariam/config.yaml
chmod +x ./start.sh && ./start.sh main.go
```

If you want to deploy the app in a production environment

```shell
./start.sh main.go prod
```

The difference between the two docker-compose files is that one expose all the ports for local development, meanwhile the other will only expose the needed ports to expose the app and the admin interfaces.

### Testing

To run tests, just run

```shell
./start.sh test
```

When adding new directories make sure to add it to the script so you include tests in it

```shell
# --snip--
if [ "$2"=="test" ] ; then
# --snip--
go test ./core
go test ./cmd
# .. add here
# --snip--
```

## Make

### Getting help

```shell
make help
```

### Run API + PGADMIN + REDIS

```shell
make dev
```

## API

### Update swagger documentations  

```shell
./start.sh main.go swagger
```
