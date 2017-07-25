# Depcon (*Container Deployment*)

[![Build Status](https://travis-ci.org/ContainX/depcon.svg)](https://travis-ci.org/ContainX/depcon)  [![release](http://github-release-version.herokuapp.com/github/ContainX/depcon/release.svg?style=flat)](https://github.com/ContainX/depcon/releases/latest)  [![GoDoc](https://godoc.org/github.com/ContainX/depcon?status.svg)](https://godoc.org/github.com/ContainX/depcon) [![Join the chat at https://gitter.im/ContainX/community](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/ContainX/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

An advanced deployment tool for docker containers against Mesos/Marathon, Kubernetes and Amazon ECS

For full up to date documentation and usage please see:

### [http://depcon.containx.io](http://depcon.containx.io)

## Overview

Depcon makes managing clusters that run docker containers a breeze.  It offers the ability to define environments such as test, pre-prod, production against Marathon (initial release), Kubernetes and Amazon ECS.

**Some key features of Depcon are:**
- Variable interpolation in deployment descriptors
- Output results in Column, YAML and JSON formats for easy integration with automation
- Accepts deployment descriptors in either JSON or YAML format
- **Full Mesos/Marathon support**
  - Application, Group & Task management
  - Partial application updates: CPU, Memory & Scaling
  - Quick application rollback
  - Leader election, Server info and elapsed response
- **Docker compose support**
  - Supports all major operating systems
  - Works with Docker Machine aka Boot2Docker
  - Parameter support -- `${PARAMS}` can be placed in compose files
- Future releases will offer a proposed Open Deployment Descriptor format which will allow Depcon to deploy a common descriptor against Marathon, Kubernetes, ECS and Docker Swarm.
- Ability to wait until a new or updated container deployment is healthy

### Binary Installation

Binaries are available through GitHub releases.  You can download the appropriate binary, package and version from the [Releases](https://github.com/ContainX/depcon/releases) page

### Build and Install the Binaries from Source

**Pre-Requisites**
* GOLANG 1.6+

Add Depcon and its package dependencies to your go `src` directory

    go get -v github.com/ContainX/depcon

Once the `get` has completed, you should find your new `depcon` (or `depcon.exe`) executable sitting inside the `$GOPATH/bin/`

To update Depcon's dependencies, use `go get` with the `-u` option.

    go get -u -v github.com/ContainX/depcon


### Running Depcon in Docker

With each release we publish a very small docker image containing depcon.

**Quick Example**
```
# Run depcon in the background capture cid, add alias
$ cid=$(docker run -itd -v $PWD:/data pacesys/depcon)
$ alias depcon="docker exec -it ${cid} depcon"

# Use depcon like it was native
$ depcon app list

```

For additional instructions in testing depcon within docker (useful for CI systems) see the docker hub repository at: https://hub.docker.com/r/pacesys/depcon/

## Global options in Depcon

#### Output Options

Depcon makes it easy to integrate with third party systems.  Any command or query in depcon has the options to list results in tabular, json or yaml formats.

For example:  `depcon app list -o json` would return a list of running applications in JSON form.  You can also use `-o yaml` for yaml or no option which by default results in table/tabular form.

## Using Depcon with Mesos/Marathon

### Applications

Below are examples with application managements

#### Listing deployed applications

List all applications

```
$ depcon app list
```

#### Getting details about a running application by it's ID

Gets an application details by Id

```
$ depcon app get myapp
```

#### Destroy/Delete a running application

Remove an application [applicationId] and all of it's instances

```
$ decon app destroy myapp
```

#### Scale an Application

Scales [appliationId] to total [instances]

```
$ depcon app scale myapp 2
```

#### Restart a running application

Restarts an application by Id

```
$ depcon app restart myapp
```

#### Update a running application

```
// Update CPU resources
$ depcon app update cpu myapp 0.5

// Update Memory to 400mb
$ depcon app update mem myapp 400
```

## Using Depcon as a Docker Compose client

Depcon supports Docker Compose natively on all major operating systems.  This feature is currently in beta, please report any found issues.

**Available Docker Compose Actions**

```
$ depcon compose

Usage:
  depcon compose [command]

Available Commands:
  build       Build or rebuild services
  kill        Kill containers
  logs        View output from containers
  port        Stops services
  ps          List containers
  up          Create and start containers
  pull        Pulls service imagess
  restart     Restart running containers
  rm          Remove stopped containers
  start       Start services
  stop        Stops services
  up          Create and start containers

Flags:
      --compose-file="docker-compose.yml": Docker compose file
  -h, --help[=false]: help for compose
      --name="depcon_proj": Project name for this composition


Global Flags:
  -e, --env="": Specifies the Environment name to use (eg. test | prod | etc). This can be omitted if only a single environment has been defined
  -o, --output="column": Specifies the output format [column | json | yaml]
      --verbose[=false]: Enables debug/verbose logging


Use "depcon compose [command] --help" for more information about a command.
```

The examples below assume `docker-compose.yml` is found in the execution directory.  If the compose file is located in another location then
the global `--compose-file` flag can be invoked.

#### Creating and Starting containers

```
$ depcon compose up
```

#### Creating and Starting a specific service

```
$ depcon compose up redis
```

#### Stopping compose services

```
$ depcon compose kill
```

### Using parameters within Compose templates

Depcon offers extenability on top of tradditional Docker compose.  It allows params to be placed within compose files in the format of `${PARAM}`.  Depcon allows these params to be resolved via the flag `--param PARAM=value` during use or via exported env variables.

Take a look at the `samples/docker-compose-params.yml` in Depcon source repo.  Here's an example of params using the referenced sample compose file.

```
// Inline params
$ depcon compose up redis --compose-file samples/docker-compose-params.yml -p REDIS_PORT=6379

// As env variables
$ export REDIS_PORT=6379
$ depcon compose up redis --compose-file samples/docker-compose-params.yml
```

## License

This software is licensed under the Apache 2 license, quoted below.

Copyright 2016 Jeremy Unruh

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.
