# raid [![Build Status](https://travis-ci.org/rai-project/raid.svg?branch=master)](https://travis-ci.org/rai-project/raid)

## Terminology

We will use the following terms throughout this document. For the publically available services, this
document assume reader is familiar with them, their limitations, and semantics:

* Client/rai: the client which users use to interact with the RAI system. This includes the RAI client as well as the docker builder website. The client is usually installed on user’s machine and is used primarily to submit tasks to the system. At any point in time, there can be more thanone client running and submitting jobs to the system. The client should work on any OS and does not require any special hardware to be installed on the system.
* Job Queue: User’s jobs are submitted onto a queue (using sharding, for example). The queue currently uses Amazon’s SQS queue system. There can be more than one queue, but we currently use only one.
* Pub/Sub Queue: Output for jobs are published onto a pub/sub server. RAI currently uses Redis for pub/sub. Multiple redis servers can be used, but we currently use only one.
* Server/raid: All work/execution is run on the server. The server listens to the queue and executes jobs if capable. Depending on the load, more than one server can be run at any point in time. The number of jobs that a server can handle is configurable. Linux is required to use the server with GPU capabilities.
* Docker: User code execution occurs only within a docker container. Docker is also used to build docker images as well as publishing images to the docker registry. * CUDA Volume: A docker plugin that mounts nvidia driver and cuda libraries within the launched docker container.
* STS: Amazon’s STS allows one to place a time constraints on the amazon access tokens (also known as roles) being issued. The current constraint is 15 minutes.
* Store/File Server: The location where user’s data files are stored. The rai system currently uses Amazon’s S3 for storage.
* Auth: Only users with credentials can submit tasks to rai. Authentication keys can be generated using the rai-keygen tool. In the backend, we use Auth0 as the database. 
* App secret: all keys, such as credentials to login to the pub/sub server, are encrypted using AES32 based encryption. For prebuilt binaries, the app secret is embedded in the executable. A user can specify the secret from the command line as well.

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

### SQS Queue

Create an SQS queue using the AWS console

### S3 Bucket

Create an S3 bucket using the AWS console

### Redis Server

Install a redis server. A docker container can also be used.

### CUDA Volume Plugin

Follow the instructions at [rai-project/rai-docker-volume](https://github.com/rai-project/rai-docker-volume#rai-docker-volume).

Prebuilt binaries exist on S3 at /files.rai-project.com/dist/rai-docker-volume/stable/latest.

_These binaries are not publicly readable, you need an AWS_KEY / SECRET to access them._

## RAI Client Installation

See [rai-project/rai](https://github.com/rai-project/rai#download-binaries)

## RAID Server Installation from Binary

Prebuilt raid binaries exist on s3 in /files.rai-project.com/dist/raid/stable/latest

_These binaries are not publicly readable, you need an AWS_KEY / SECRET to access them._

## RAID Server Installation from Source

1. Install golang. Either through [Go Version Manager](https://github.com/moovweb/gvm)(recommended) or from the instructions on the [golang site](https://golang.org/). We recommend the Go Version Manager.
2. (Optional) Install [glide](https://github.com/Masterminds/glide#install)
3. Clone the `raid` repository

        mkdir -p $GOPATH/src/github.com/rai-project
        cd $GOPATH/src/github.com/rai-project
        git clone git@github.com:rai-project/raid.git

4. Install the software dependencies using `glide`.
    1. If you installed `glide` in step 2

        cd raid
        glide install

    2. If you did not

        cd raid
        go get -u -v ./...

5. Create an executable (optionally, embed the secret. You won't have to use the `-s` flag later)

        go build

    or

        go build -ldflags="-s -w -X main.AppSecret=${APP_SECRET}"


## Usage


## SystemD



## Logs

If `journald` is enabled, then you can view the server logs using `journalctl -f -u raid.service`


## License

NCSA/UIUC © [Abdul Dakkak](http://impact.crhc.illinois.edu/Content_Page.aspx?student_pg=Default-dakkak)
