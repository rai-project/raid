# raid [![Build Status](https://travis-ci.org/rai-project/raid.svg?branch=master)](https://travis-ci.org/rai-project/raid)

## Terminology

We will use the following terms throughout this document. For the publically available services, this
document assume reader is familiar with them, their limitations, and semantics:

* **Client/rai**: the client which users use to interact with the RAI system. This includes the RAI client as well as the docker builder website. The client is usually installed on user’s machine and is used primarily to submit tasks to the system. At any point in time, there can be more thanone client running and submitting jobs to the system. The client should work on any OS and does not require any special hardware to be installed on the system.
* **Job Queue**: User’s jobs are submitted onto a queue (using sharding, for example). The queue currently uses Amazon’s SQS queue system. There can be more than one queue, but we currently use only one.
* **Pub/Sub Queue**: Output for jobs are published onto a pub/sub server. RAI currently uses Redis for pub/sub. Multiple redis servers can be used, but we currently use only one.
* **Server/raid**: All work/execution is run on the server. The server listens to the queue and executes jobs if capable. Depending on the load, more than one server can be run at any point in time. The number of jobs that a server can handle is configurable. Linux is required to use the server with GPU capabilities.
* **Docker**: User code execution occurs only within a docker container. Docker is also used to build docker images as well as publishing images to the docker registry. * CUDA Volume: A docker plugin that mounts nvidia driver and cuda libraries within the launched docker container.
* **STS**: Amazon’s STS allows one to place a time constraints on the amazon access tokens (also known as roles) being issued. The current constraint is 15 minutes.
* **Store/File Server**: The location where user’s data files are stored. The rai system currently uses Amazon’s S3 for storage.
* **Auth**: Only users with credentials can submit tasks to rai. Authentication keys can be generated using the rai-keygen tool. In the backend, we use Auth0 as the database. 
* **App secret**: all keys, such as credentials to login to the pub/sub server, are encrypted using AES32 based encryption. For prebuilt binaries, the app secret is embedded in the executable. A user can specify the secret from the command line as well.

## Components

## Execution Flow

1. A client submits a task to RAI
    1. Credentials are validated
    2. The directory is archived and uploaded to the file server
    3. The task is submitted to the queue
    4. The user subscribes to the pub/sub server and prints all messages being received
2. A server accepts a task from the queue
    1. Check if the task is valid
    2. Either build or downloads the docker images required to run the task
    3. Download the user directory
    4. Start a publish channel to the pub/sub server
    5. Start a docker container and run the user commands (if gpu is requested, then the cuda
    volume is mounted)
    6. All stdout/stderr messages from the docker container are forwarded to the publish
    channel
    7. Wait for either tasks to complete or a timeout
    h. The output directory is uploaded to the file server and a link is published
    8. Close the publish channel / docker container

## Prerequisites

### STS Permissions

Set up STS permissions so that users can upload to an S3 bucket and publish to an SQS queue. This
can be done via the IAM AWS console. The STS currently being used only allows one to attach to a rai
role which has very limited permissions:

### Simple Queue Service

Create an SQS queue using the AWS console.

* Navigate to the SQS service page
* "Create New Queue"
* Enter the queue name
* Choose standard for the type.
* Optionally, configure various queue parameters under "configure queue"

Create a **IAM Policy** that allows reading and writing to the new queue.
Use the **Policy Generator**. 
* Select Amazon SQS for AWS Service. 
* Choose the following actions:
    * GetQueueUrl
    * ReceiveMessage
    * DeleteMessage
    * DeleteMessageBatch
    * SendMessage
    * SendMessageBatch
* Create an [ARN](http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html). The account ID may be found on the AWS account page. For example `arn:aws:sqs:*:account-id:rai2` for the `rai2` sqs queue.

The ARN controls which queues the policy applies to.
For example, `arn:aws:sqs:*:account-id:rai*` will apply to all queues that match `rai*`.

### S3 Bucket

Create an S3 bucket using the AWS console

### MongoDB

Create a mongodb to store submissions from the client.

* Create a security group that allows ssh (port 22) and mongodb (port 27017)
* Create an AWS EC2 instance to run the database and add it to that security group

Install docker

    curl -fsSL get.docker.com -o get-docker.sh | sudo sh
    sudo usermod -aG docker $USER

Log out and log back in. 
Start mongo 3.0

    docker run -p 27017:27017 --restart always -d --name rai-mongo -v /data/db:/data/db mongo:3.0 --auth

Takes a while to preallocate some files. YOu can monitor with `docker logs -f rai-mongo`. Then connect to the admin database as admin

    docker exec -it rai-mongo mongo --authenticationDatabase admin admin

Add a rai-root user that can administer any database

    db.createUser({ user: 'rai-root', pwd: 'some-password', roles: [ { role: "root", db: "admin" } ] });

Exit and connect to the admin database as that user

    docker exec -it rai-mongo mongo --authenticationDatabase admin -u rai-root -p some-password admin

Switch to the rai database
This doesn't actually create a database until you put something into the database

    use rai

Create a collection for the submissions.

    db.createCollection("rankings")

Now that the rankings database exists, add a user for the rai-client

    db.createUser({ user: 'rai-client', pwd: 'some-password', roles: [ { role: "readWrite", db: "rai" } ] });

To nuke the database and start from scratch if you goof up:

    docker rm -f `docker ps -a -q` && docker volume prune
 


### Redis Server

Install a redis server. A docker container can also be used.

### CUDA Volume Plugin

Follow the instructions at [rai-project/rai-docker-volume](https://github.com/rai-project/rai-docker-volume#rai-docker-volume).

Prebuilt binaries exist on S3 at /files.rai-project.com/dist/rai-docker-volume/stable/latest.

> _These binaries are not publicly readable, you need an AWS_KEY / SECRET to access them._

## RAI Client Installation

See [rai-project/rai](https://github.com/rai-project/rai#download-binaries)

## RAID Server Installation from Binary

Prebuilt raid binaries exist on s3 in /files.rai-project.com/dist/raid/stable/latest

> _These binaries are not publicly readable, you need an AWS_KEY / SECRET to access them._

## RAID Server Installation from Source

1. Install golang. Either through [Go Version Manager](https://github.com/moovweb/gvm)(recommended) or from the instructions on the [golang site](https://golang.org/). We recommend the Go Version Manager.
2. (Optional) Install [glide](https://github.com/Masterminds/glide#install)
3. Clone the `raid` repository

```sh
        mkdir -p $GOPATH/src/github.com/rai-project
        cd $GOPATH/src/github.com/rai-project
        git clone git@github.com:rai-project/raid.git
```

4. Install the software dependencies using `glide`.
a. If you installed `glide` in step 2

```sh
        cd raid
        glide install
```


b. If you did not


```sh
        cd raid
        go get -u -v ./...
```

5. Create an executable (optionally, embed the secret. You won't have to use the `-s` flag later)

```sh
        go build
```

or


```sh
        go build -ldflags="-s -w -X main.AppSecret=${APP_SECRET}"
```

## Configuration

Much of rai/raid is controlled by configuration files. Services that are shared between the client and server must match. In this section we will explain the minimal configurations needed for both the client and server

> **Note:** One can create secret keys recognizable by rai/raid using [rai-crypto](https://github.com/rai-project/utils/tree/master/rai-crypto) tool.
> If you want to encrypt a string using “PASS” as your app secret, then you’d want to invoke

```sh
    rai-crypto encrypt –s PASS MY_PLAIN_TEXT_STRING
```

Configurations are specified in YAML format and exist either in $HOME/.rai_config.yml or are embedded within the executable. There are many more configurations that can be set, but if omitted then sensible defaults are used.

### Client Configuration

The client configuration configures the client for usage with a cluster of rai servers.

```yaml
    app:
        name: rai # name of the application
        verbose: false # whether to enable verbosity by default
        debug: false # whether to enable debug output by default
    aws:
        access_key_id: AWS_ACCESS_KEY # the aws access key (encrypted)
        secret_access_key: AWS_SECRET_KEY # the aws secret key
        region: us-east-1 # the aws region to use
        sts_account: STS_ACCOUNT # the sts account number
        sts_role: rai # the sts role
        sts_role_duration_seconds: 15m # the maximum time period the sts role can be assumed
    store:
        provider: s3 # the store provider
        base_url: http://s3.amazonaws.com # the base url of the file store
        acl: public-read # the default permissions for the files uploaded to the file store
    client:
        name: rai # name of the client
        upload_bucket: files.rai-project.com # base url or the store buceket
        bucket: userdata # location to store the uploaded user data (user input)
        build_file: rai_build # location to store the result build data (user output)
    auth:
        provider: auth0 # the authentication provider
        domain: raiproject.auth0.com # the domain of the authentication provider
        client_id: AUTH0_CLIENT_ID # the client id from the authentication. The auth0 client id for example
        client_secret: AUTH0_CLIENT_SECRET # the client secret from the authentication. The auth0 client secret for example
    pubsub:
        endpoints: # list of endpoints for the pub sub service
            - pubsub.rai-project.com:6379 # the pubsub server location + port
        password: PUBSUB_PASSWORD # password to the pub/sub service
```

> **Note:** During the travis build process the client configurations are embedded into the client binary. Therefore the $HOME/.rai_config.yml is never read.

### Server Configuration

All servers within a cluster share the same configuration. Here is the configuration currently being used:

```yaml
    app:
        name: rai # name of the application
        verbose: false # whether to enable verbosity by default
        debug: false # whether to enable debug output by default
    logger:
        hooks: # log hooks to enable
            - syslog # enable the syslog hook
    aws:
        access_key_id: AWS_ACCESS_KEY # the aws access key (encrypted)
        secret_access_key: AWS_SECRET_KEY # the aws secret key
        region: us-east-1 # the aws region to use
    store:
        provider: s3 # the store provider
        base_url: http://s3.amazonaws.com # the base url of the file store
        acl: public-read # the default permissions for the files uploaded to the file store
    broker: # broker/queue configuration section
        provider: sqs # the queue provider
        serializer: json # the serialization method to use for messages. Json is the default
        autoack: true # enable auto acknowledgement of messages
    client:
        name: rai # name of the client
        upload_bucket: files.rai-project.com # base url or the store buceket
        bucket: userdata # location to store the uploaded user data (user input)
        build_file: rai_build # location to store the result build data (user output)
    auth:
        provider: auth0 # the authentication provider
        domain: raiproject.auth0.com # the domain of the authentication provider
        client_id: AUTH0_CLIENT_ID # the client id from the authentication. The auth0 client id for example
        client_secret: AUTH0_CLIENT_SECRET # the client secret from the authentication. The auth0 client secret for example
    pubsub:
        endpoints: # list of endpoints for the pub sub service
            - pubsub.rai-project.com:6379 # the pubsub server location + port
        password: PUBSUB_PASSWORD # password to the pub/sub service
```

Other useful configuration options are `docker.time_limit` (default 1 hour), `docker.memory_limit` (default
16gb)

### Start/Stop Server

With the absence of integration of raid with system service management (such as SystemD, UpStart, ...) one needs to start and stop the raid server manually. \
Assuming you’ve already compiled the raid executable, you can run the server using the following command:

```sh
    ./raid –s ${MY_SECRET}
```

There are a few options that are available to control settings:

Table 1 : Command line options for raid server

| Description | Option |
| -- | -- |
| **Debug Mode** | -d |
| **Verbose Mode** | -v |
| **Application Secret** | -s MY_SECRET |

The above command will exit when a user exists the terminal session. Use the nohup command to avoid that

```sh
    nohup ./raid –d -v –s ${MY_SECRET} &
```

#### Setting up NVIDIA Persistence Mode

If using CUDA, make sure to [enable persistence mode](http://docs.nvidia.com/deploy/driver-persistence/index.html).


Copy the the systemd service in `build/systemd/nvidia-persistenced.service` and modify the line

```
ExecStart=/usr/bin/nvidia-persistenced --user ubuntu
```

to launch using the user you want (here it's launching under the `ubuntu` user).
Once modified, copy the file to `/etc/systemd/system/` and start the systemd service

```sh
sudo mv nvidia-persistenced.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start nvidia-persistenced.service
```

Make sure that the `nvidia-persistenced.service` service has been started. 
The command

```sh
sudo systemctl status nvidia-persistenced.service
```

should give an output similar to

```sh
● nvidia-persistenced.service - NVIDIA Persistence Daemon
   Loaded: loaded (/etc/systemd/system/nvidia-persistenced.service; disabled; vendor preset: enabled)
   Active: active (running) since Fri 2017-11-03 18:33:42 CDT; 5min ago
  Process: 71100 ExecStopPost=/bin/rm -rf /var/run/nvidia-persistenced (code=exited, status=0/SUCCESS)
  Process: 71103 ExecStart=/usr/bin/nvidia-persistenced --user abduld (code=exited, status=0/SUCCESS)
 Main PID: 71107 (nvidia-persiste)
    Tasks: 1
   Memory: 492.0K
      CPU: 3ms
   CGroup: /system.slice/nvidia-persistenced.service
           └─71107 /usr/bin/nvidia-persistenced --user abduld

Nov 03 18:33:42 whatever systemd[1]: Starting NVIDIA Persistence Daemon...
Nov 03 18:33:42 whatever nvidia-persistenced[71107]: Started (71107)
Nov 03 18:33:42 whatever systemd[1]: Started NVIDIA Persistence Daemon.
```

Finally, check that the driver is run in persistence mode

```sh
$ nvidia-smi -a | grep Pe
    Persistence Mode                : Enabled
        Pending                     : N/A
        Pending                     : N/A
    Performance State               : P5
        Pending                     : N/A
        Pending                     : N/A
```

#### Creating RAI Accounts

Either build the rai-keygen or download the prebuilt binaries which exist on s3 in /files.rai-project.com/dist/rai-keygen/stable/latest

> **Note:** These binaries are not publically readable and you need an AWS_KEY / SECRET to access them.

One can use the [rai-keygen](https://github.com/rai-project/rai-keygen) to generate RAI and email account information to the people enrolled in the class.
The mailing process uses mailgun.
A prebuilt rai-keygen includes a builtin configuration file, but if compiling from source, then you need to add the email configuration options

    email:
        provider: mailgun # the email provider
        domain: email.webgpu.com # the domain of the
        source: postmaster@webgpu.com # the source email
        mailgun_active_api_key: API_KEY
        mailgun_email_validation_key: VALIDATION_KEY

You will not need the above if you do not need to email the generated keys.

> **Note:** Docker builder does not require account generation, since account information is embedded into the webserver.

### Administration Tips

The following tips are based on past experience administering clusters and managing arbitrary user
execution:
* If using CUDA, make sure to [enable persistence mode](http://docs.nvidia.com/deploy/driver-persistence/index.html). A systemd service exists within the raid repository
* A system can become unstable when executing arbitrary code. Consult the logs (ideally a distributed logging) when trying to identify why certain tasks succeed while other fail.
* Install Cadvistor (github.com/google/cadvisor) to examine the health of the docker container and monitor them.
* Make sure that you have enough disk space. For example, last year the redis server ran out of disk space 2-3 days before the deadline.


#### AWS Admin Notes


##### Reboot all AWS Instances

```
instances=$(aws ec2 describe-instances --filters "Name=tag:name,Values=ece408.project" "Name=instance-state-code,Values=16" | jq -j '[.Reservations[].Instances[].InstanceId] | @sh')
echo ${instances}
for instance in ${instances}
do
  echo ${instance}
  #aws ec2 reboot-instances --dry-run --instance-ids ${instance}
done
```

or

```
instances=$(aws ec2 describe-instances --filters "Name=tag:name,Values=ece408.project" "Name=instance-state-code,Values=16" | jq -j '[.Reservations[].Instances[].InstanceId] | join(" --instance-ids ")')
aws ec2 reboot-instances --no-dry-run --instance-ids ${instances}
```

## Logs

If `journald` is enabled, then you can view the server logs using `journalctl -f -u raid.service`


## License

NCSA/UIUC © [Abdul Dakkak](http://impact.crhc.illinois.edu/Content_Page.aspx?student_pg=Default-dakkak)
