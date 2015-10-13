# Depcon (*Container Deployment*)

[![Build Status](https://travis-ci.org/gondor/depcon.svg)](https://travis-ci.org/gondor/depcon)  [![GoDoc](https://godoc.org/github.com/gondor/depcon?status.svg)](https://godoc.org/github.com/gondor/depcon) [![Join the chat at https://gitter.im/gondor/depcon](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gondor/depcon?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) 

An advanced deployment tool for docker containers against Mesos/Marathon, Kubernetes and Amazon ECS

## Overview

Depcon makes managing clusters that run docker containers a breeze.  It offers the ability to define environments such as test, pre-prod, production against Marathon (initial release), Kubernetes and Amazon ECS.  

**Some key features of Depcon are:**
- Variable substistution in deployment descriptors
- Output results in Column, YAML and JSON formats for easy integration with automation
- Accepts deployment descriptors in either JSON or YAML format
- Full Mesos/Marathon support
  - Application, Group & Task management
  - Partial application updates: CPU, Memory & Scaling
  - Quick application rollback
  - Leader election, Server info and elapsed response
- Future releases will offer a proposed Open Deployment Descriptor format which will allow Depcon to deploy a common descriptor against Marathon, Kubernetes, ECS and Docker Swarm.
- Ability to wait until a new or updated container deployment is healthy

### Binary Installation 

You can download the binaries (ver 0.1)

 * Architecture i386 [ [linux](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_linux_386.tar.gz?direct) / [windows](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_windows_386.zip?direct) / [darwin](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_darwin_386.zip?direct) / [freebsd](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_freebsd_386.zip?direct) / [openbsd](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_openbsd_386.zip?direct) ]
 * Architecture amd64 [ [linux](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_linux_amd64.tar.gz?direct) / [windows](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_windows_amd64.zip?direct) / [darwin](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_darwin_amd64.zip?direct) / [freebsd](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_freebsd_amd64.zip?direct) / [openbsd](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_openbsd_amd64.zip?direct) ]

Or by installing via deb packages (ver 0.1)

 * [ [amd64](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_amd64.deb?direct) / [armhf](https://dl.bintray.com//content/pacesys/utils/depcon_0.1_armhf.deb?direct) ]


### Build and Install the Binaries from Source

Add Depcon and its package dependencies to your go `src` directory

    go get -v github.com/gondor/depcon

Once the `get` has completed, you should find your new `depcon` (or `depcon.exe`) executable sitting inside the `$GOPATH/bin/`

To update Depcon's dependencies, use `go get` with the `-u` option.

    go get -u -v github.com/gondor/depcon

## License

This software is licensed under the Apache 2 license, quoted below.

Copyright 2015 Jeremy Unruh

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.
