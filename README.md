![Yopass-horizontal](https://user-images.githubusercontent.com/37777956/59544367-0867aa80-8f09-11e9-8d6a-02008e1bccc7.png)

# Yopass - Share Secrets Securely

[![Go Report Card](https://goreportcard.com/badge/github.com/jhaals/yopass)](https://goreportcard.com/report/github.com/jhaals/yopass)
[![codecov](https://codecov.io/gh/jhaals/yopass/branch/master/graph/badge.svg)](https://codecov.io/gh/jhaals/yopass)

![demo](https://ydemo.netlify.com/yopass-demo.gif)

### Based on [jhaals/yopass](https://github.com/jhaals/yopass/)

Yopass is a project for sharing secrets in a quick and secure manner\*.
The sole purpose of Yopass is to minimize the amount of passwords floating around in ticket management systems, Slack messages and emails. The message is encrypted/decrypted locally in the browser and then sent to yopass without the decryption key which is only visible once during encryption, yopass then returns a one-time URL with specified expiry date.

There is no perfect way of sharing secrets online and there is a trade off in every implementation. Yopass is designed to be as simple and "dumb" as possible without compromising on security. There's no mapping between the generated UUID and the user that submitted the encrypted message. It's always best to send all the context except password over another channel.

**[Demo available here](https://yopass.se)**. It's recommended to host yopass yourself if you care about security.

- End-to-End encryption using [OpenPGP](https://openpgpjs.org/)
- Secrets can only be viewed once
- No accounts or user management required
- Secrets self destruct after X hours
- Custom password option
- Limited file upload functionality

## History

Yopass was first released in 2014 and has since then been maintained by me and contributed to by this fantastic group of [contributors](https://github.com/jhaals/yopass/graphs/contributors). Yopass is used by many large corporations none of which are currently listed in this readme.
If you are using yopass and want to support other then by code contributions. Give your thanks in an email, consider donating or by giving consent to list your company name as a user of Yopass in this readme(Trusted by)

## Command-line interface

The main motivation of Yopass is to make it easy for everyone to share secrets easily and quickly via a simple webinterface. Nevertheless, a command-line interface is provided as well to support use cases where the output of a program needs to be shared.

```console
$ yopass --help
Yopass - Secure sharing for secrets, passwords and files

Flags:
      --api string          Yopass API server location (default "https://api.yopass.se")
      --decrypt string      Decrypt secret URL
      --expiration string   Duration after which secret will be deleted [1h, 1d, 1w] (default "1h")
      --file string         Read secret from file instead of stdin
      --key string          Manual encryption/decryption key
      --one-time            One-time download (default true)
      --url string          Yopass public URL (default "https://yopass.se")

Settings are read from flags, environment variables, or a config file located at
~/.config/yopass/defaults.<json,toml,yml,hcl,ini,...> in this order. Environment
variables have to be prefixed with YOPASS_ and dashes become underscores.

Examples:
      # Encrypt and share secret from stdin
      printf 'secret message' | yopass

      # Encrypt and share secret file
      yopass --file /path/to/secret.conf

      # Share secret multiple time a whole day
      cat secret-notes.md | yopass --expiration=1d --one-time=false

      # Decrypt secret to stdout
      yopass --decrypt https://yopass.se/#/...

Website: https://yopass.se
```

The following options are currently available to install the CLI locally.

- Compile from source (needs Go >= v1.15)

  ```console
  export GO111MODULE=on && go get github.com/z0x0z/yopass/cmd/yopass && go install github.com/z0x0z/yopass/cmd/yopass
  ```

- Arch Linux ([AUR package](https://aur.archlinux.org/packages/yopass/))

  ```console
  yay -S yopass
  ```

## Installation / Configuration

Here are the server configuration options.

Command line flags:

```console
$ yopass-server -h
      --address string     listen address (default 0.0.0.0)
      --database string    database backend ('memcached' or 'redis') (default "memcached")
      --max-length int     max length of encrypted secret (default 10000)
      --memcached string   Memcached address (default "localhost:11211")
      --metrics-port int   metrics server listen port (default -1)
      --port int           listen port (default 1337)
      --redis string       Redis URL (default "redis://localhost:6379/0")
      --tls-cert string    path to TLS certificate
      --tls-key string     path to TLS key
```

Encrypted secrets can be stored either in Memcached or Redis by changing the `--database` flag.

## Development

- **To run application locally using go build**
  Clone the repo and follow steps below
  To download go dependencies ( go should be installed)
  ```bash
  go mod tidy
  ```
  To run cli client
  ```bash
  cd cmd/yopass
  go build
  ./yopass --expiration="0" --one-time=false --api http://localhost:1337  --key "random" --url http://localhost:1337 <<< 'testing'
  ```
  To run server
  ```bash
  cd cmd/yopass-server
  go build
  ./yopass-server --database "redis"
  ```
  To run website
  ```bash
  cd website
  yarn install
  REACT_APP_BACKEND_URL='http://localhost:1337' yarn start
  ```
  To run DB, using docker
  ```bash
  docker run -p 6379:6379 redis
  ```
  Browse the application using the URL,
  ```bash
  http://localhost:3000
  ```
- **To run application locally using Docker containers**
  Build Yopass docker container using the Dockerfile present in https://github.com/z0x0z/Yopass
  ```bash
  docker build . -t yopass
  ```
  start the redis container
  ```bash
  docker run --name redis -p 6379:6379 redis
  ```
  start the application container
  ```bash
  docker run -p 1337:1337 yopass:latest --database=redis --redis=redis://<redis container ip>:6379
  ```
  To run website
  ```bash
  cd website
  yarn install
  REACT_APP_BACKEND_URL='http://localhost:1337' yarn start
  ```
  Browse the application using the URL,
  ```bash
  http://localhost:3000
  ```
- **To run application locally using Docker-Compose**
  Build Yopass docker container using the Dockerfile present in https://github.com/z0x0z/Yopass
  ```bash
  docker build . -t yopass
  ```
  Save the below code in a docker-compose.yml file
  ```bash
  version: '3.8'
  services:
    redis:
      image: redis
      container_name: redis
      restart: always
      ports:
        - 6379:6379
    yopass:
      image: yopass
      container_name: yopass
      ports:
        - 1337:1337
      restart: always
      depends_on:
        - redis
      command: "--database redis --redis=redis://redis:6379"
  ```
  run the below command spin up application and redis containers
  ```bash
  docker compose up
  ```
  To run website
  ```bash
  cd website
  yarn install
  REACT_APP_BACKEND_URL='http://localhost:1337' yarn start
  ```
  Browse the application using the URL,
  ```bash
  http://localhost:3000
  ```
- **To run application in AWS EC2 using Docker containers**

  Build the application locally and push to any container repo (ex. DockerHub)

  make sure you have logged in to container repository

  ```bash
  docker build . -t <reponame>/yopass
  docker push <reponame>/yopass
  ```

  Create an ElasticCache or Memorydb cluster in AWS which acts as redis database

  make sure the connections, security groups and appropriate permissions are in place

  Pull the docker container inside EC2 and run the container

  ```bash
  docker pull <reponame>/yopass
  docker run -p 1337:1337 z0x0z/yopass:linux --database redis --redis=rediss://<username>:<password>@clustercfg.yopass.6nwcyd.memorydb.ap-south-1.amazonaws.com:6379
  ```

  Browse the application using the URL,

  ```bash
  http://<instance_ip>:1337
  ```

- **To run application in AWS ECS using Docker containers**

  Build the application locally and push to any container repo (ex. DockerHub)

  make sure you have logged in to container repository

  ```bash
  docker build . -t <reponame>/yopass
  docker push <reponame>/yopass
  ```

  Create an ECS Cluster - Give it a name and use AWS Fargate as infrastructure

  Create a task definition using below json

  Note: Review & change the docker commands and url of Memorydb in below file

  ```json
  {
    "taskDefinitionArn": "arn:aws:ecs:ap-south-1:0216750xxx95:task-definition/yopass-definition:2",
    "containerDefinitions": [
      {
        "name": "yopass",
        "image": "docker.io/z0x0z/yopass:linux",
        "cpu": 0,
        "portMappings": [
          {
            "name": "yopass-1337-tcp",
            "containerPort": 1337,
            "hostPort": 1337,
            "protocol": "tcp",
            "appProtocol": "http"
          }
        ],
        "essential": true,
        "command": [
          "--database",
          "redis",
          "--redis",
          "rediss://clustercfg.yopass.6xxxxd.memorydb.ap-south-1.amazonaws.com:6379"
        ],
        "environment": [],
        "mountPoints": [],
        "volumesFrom": [],
        "logConfiguration": {
          "logDriver": "awslogs",
          "options": {
            "awslogs-create-group": "true",
            "awslogs-group": "/ecs/yopass-definition",
            "awslogs-region": "ap-southeast-1",
            "awslogs-stream-prefix": "ecs"
          }
        },
        "systemControls": []
      }
    ],
    "family": "yopass-definition",
    "executionRoleArn": "arn:aws:iam::021611124395:role/ecsTaskExecutionRole",
    "networkMode": "awsvpc",
    "revision": 2,
    "volumes": [],
    "status": "ACTIVE",
    "requiresAttributes": [
      {
        "name": "com.amazonaws.ecs.capability.logging-driver.awslogs"
      },
      {
        "name": "ecs.capability.execution-role-awslogs"
      },
      {
        "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
      },
      {
        "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
      },
      {
        "name": "ecs.capability.task-eni"
      },
      {
        "name": "com.amazonaws.ecs.capability.docker-remote-api.1.29"
      }
    ],
    "placementConstraints": [],
    "compatibilities": ["EC2", "FARGATE"],
    "requiresCompatibilities": ["FARGATE"],
    "cpu": "1024",
    "memory": "3072",
    "runtimePlatform": {
      "cpuArchitecture": "X86_64",
      "operatingSystemFamily": "LINUX"
    },
    "registeredAt": "2024-03-09T19:39:32.178Z",
    "registeredBy": "arn:aws:iam::021611134395:user/gopikrishna",
    "tags": []
  }
  ```

  Create a service in the ECS Cluster and ensure the correct task definition is selected. Enable public ip in Networking. Set needed inbound rules in Security group

  Once the service is deployed successfully, Navigate to tasks tab inside cluster and click on the task.

  Copy the public ip of the container

  Browse the application using the URL,

  ```bash
  http://<public_ip>:1337
  ```

- **To run application behind VPN with SSL using AWS ECS using Docker containers**

  Build the application locally and push to any container repo (ex. DockerHub)

  make sure you have logged in to container repository

  ```bash
  docker build . -t <reponame>/yopass
  docker push <reponame>/yopass
  ```

  Create an ElasticCache or Memorydb cluster in AWS which acts as redis database

  Create a security group which allow 0.0.0.0/0 traffic and attach to Memoryd

  Create two Security Groups under VPC which is peered to VPN

  1. ALB-SG → to Allow inbound traffic only on port 443 from anywhere
  2. ECS-SG →to Allow inbound all traffic only from ALB-SG

  Create an ECS Cluster - Give it a name and use AWS Fargate as infrastructure

  Create a task definition using below json

  Note: Review & change the docker commands and url of Memorydb in below file

  ```json
  {
    "taskDefinitionArn": "arn:aws:ecs:ap-south-1:021673454395:task-definition/yopass-definition:2",
    "containerDefinitions": [
      {
        "name": "yopass",
        "image": "docker.io/z0x0z/yopass:linux",
        "cpu": 0,
        "portMappings": [
          {
            "name": "yopass-1337-tcp",
            "containerPort": 1337,
            "hostPort": 1337,
            "protocol": "tcp",
            "appProtocol": "http"
          }
        ],
        "essential": true,
        "command": [
          "--database",
          "redis",
          "--redis",
          "rediss://clustercfg.yopass.6xxxxd.memorydb.ap-south-1.amazonaws.com:6379"
        ],
        "environment": [],
        "mountPoints": [],
        "volumesFrom": [],
        "logConfiguration": {
          "logDriver": "awslogs",
          "options": {
            "awslogs-create-group": "true",
            "awslogs-group": "/ecs/yopass-definition",
            "awslogs-region": "ap-southeast-1",
            "awslogs-stream-prefix": "ecs"
          }
        },
        "systemControls": []
      }
    ],
    "family": "yopass-definition",
    "executionRoleArn": "arn:aws:iam::021123434395:role/ecsTaskExecutionRole",
    "networkMode": "awsvpc",
    "revision": 2,
    "volumes": [],
    "status": "ACTIVE",
    "requiresAttributes": [
      {
        "name": "com.amazonaws.ecs.capability.logging-driver.awslogs"
      },
      {
        "name": "ecs.capability.execution-role-awslogs"
      },
      {
        "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
      },
      {
        "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
      },
      {
        "name": "ecs.capability.task-eni"
      },
      {
        "name": "com.amazonaws.ecs.capability.docker-remote-api.1.29"
      }
    ],
    "placementConstraints": [],
    "compatibilities": ["EC2", "FARGATE"],
    "requiresCompatibilities": ["FARGATE"],
    "cpu": "1024",
    "memory": "3072",
    "runtimePlatform": {
      "cpuArchitecture": "X86_64",
      "operatingSystemFamily": "LINUX"
    },
    "registeredAt": "2024-03-09T19:39:32.178Z",
    "registeredBy": "arn:aws:iam::021672345395:user/gopikrishna",
    "tags": []
  }
  ```

  Create a service in the ECS Cluster and ensure the correct task definition is selected. Disable public ip in Networking. attach ECS-SG security group.

  **Create Application Load balancer** → Select scheme “internal” → select VPC which is peered to VPN → choose ALB-SG security group → Create Listener → https protocol 443 port → create target group and choose it → Under Default SSL/TLS server certificate → choose From ACM and select \*.company.com certificate → click on create Load balancer

  **Create Target group** → Select ip addresses → http protocol and 1337 port → select VPC which is peered to VPN → paste the private ip of the container and enter port as 1337 → click on include as pending below and click create

  Now copy the DNS name of the load balancer and add a CNAME entry (secrets) in cloudfare.

  Now connect to pritunl VPN,

  Browse the application using the URL,

  ```bash
  https://secrets.company.com
  ```

### Docker Compose

Use the Docker Compose file `deploy/with-nginx-and-letsencrypt/docker-compose.yml` to set up a yopass instance with TLS transport encryption and certificate auto renewal using [Let's Encrypt](https://letsencrypt.org/). First point your domain to the host you want to run yopass on. Then replace the placeholder values for `VIRTUAL_HOST`, `LETSENCRYPT_HOST` and `LETSENCRYPT_EMAIL` in `deploy/with-nginx-and-letsencrypt/docker-compose.yml` with your values. Afterwards change the directory to `deploy/with-nginx-and-letsencrypt` and start the containers with:

```console
docker-compose up -d
```

Yopass will then be available under the domain you specified through `VIRTUAL_HOST` / `LETSENCRYPT_HOST`.

Advanced users that already have a reverse proxy handling TLS connections can use the `insecure` setup:

```console
cd deploy/docker/compose/insecure
docker-compose up -d
```

Afterwards point your reverse proxy to `127.0.0.1:80`.

### Docker

With TLS encryption

```console
docker run --name memcached_yopass -d memcached
docker run -p 443:1337 -v /local/certs/:/certs \
    --link memcached_yopass:memcached -d z0x0z/yopass --memcached=memcached:11211 --tls-key=/certs/tls.key --tls-cert=/certs/tls.crt
```

Afterwards yopass will be available on port 443 through all IP addresses of the host, including public ones. If you want to limit the availability to a specific IP address use `-p` like so: `-p 127.0.0.1:443:1337`.

Without TLS encryption (needs a reverse proxy for transport encryption):

```console
docker run --name memcached_yopass -d memcached
docker run -p 127.0.0.1:80:1337 --link memcached_yopass:memcached -d z0x0z/yopass --memcached=memcached:11211
```

Afterwards point your reverse proxy that handles the TLS connections to `127.0.0.1:80`.

### AWS Lambda

_Yopass website is a separate component in this step which can be deployed to [netlify](https://netlify.com)_ for free.

You can run Yopass on AWS Lambda backed by dynamodb

```console
cd deploy/aws-lambda && ./deploy.sh
```

### Kubernetes

```console
kubectl apply -f deploy/yopass-k8.yaml
kubectl port-forward service/yopass 1337:1337
```

_This is meant to get you started, please configure TLS when running yopass for real._

## Monitoring

Yopass optionally provides metrics in the [OpenMetrics][] / [Prometheus][] text
format. Use flag `--metrics-port <port>` to let Yopass start a second HTTP
server on that port making the metrics available on path `/metrics`.

Supported metrics:

- Basic [process metrics][] with prefix `process_` (e.g. CPU, memory, and file descriptor usage)
- Go runtime metrics with prefix `go_` (e.g. Go memory usage, garbage collection statistics, etc.)
- HTTP request metrics with prefix `yopass_http_` (HTTP request counter, and HTTP request latency histogram)

[openmetrics]: https://openmetrics.io/
[prometheus]: https://prometheus.io/
[process metrics]: https://prometheus.io/docs/instrumenting/writing_clientlibs/#process-metrics

## Translations

Yopass has third party support for other languages. That means you can write translations for the language you'd like or use a third party language file. Please note that yopass itself is english only and any other translations are community supported.

Here's a list of available translations:

- [German](https://github.com/Anturix/yopass-german)
- [French](https://github.com/NicolasStr/yopass-french)
- [Spanish](https://github.com/nbensa/yopass-spanish)
- [Polish](https://github.com/mdurajewski/yopass-polish)
- [Dutch](https://github.com/KevinRosendaal/yopass-dutch)
- [Russian](https://github.com/karpechenkovkonstantin/yopass-russian)
