# h~~ealthbot~~<sup>7</sup>t

[![Docker Release](https://img.shields.io/docker/v/damianoneill/h7t?label=Docker%20Release)](https://hub.docker.com/r/damianoneill/h7t)
[![Go Report Card](https://goreportcard.com/badge/damianoneill/h7t)](http://goreportcard.com/report/damianoneill/h7t)
[![license](https://img.shields.io/github/license/damianoneill/h7t.svg)](https://github.com/damianoneill/h7t/blob/master/LICENSE)

A command line tool for interacting with [Juniper Healthbot](https://www.juniper.net/us/en/products-services/sdn/contrail/contrail-healthbot/).

> Note Device Group functionality is limited currently, this will be updated after Healthbot 2.1 is released.

[![asciicast](https://asciinema.org/a/FsdPZpORfIuciQ70rTgGvtf4M.svg)](https://asciinema.org/a/FsdPZpORfIuciQ70rTgGvtf4M)

## Synopsis

A tool for interacting with Healthbot over the REST API.

The initial focus of this tool is to provide bulk or aggregate functions, that simplify interacting with Healthbot. The initial use case is Extract Transform and Load (ETL) based.

## Commands

```console
h7t
├── configure
│   └── devices
├── extract
│   ├── device-groups
│   ├── devices
│   └── installation
├── load
│   ├── device-groups
│   ├── devices
│   └── files
│       ├── helper-files
│       ├── playbooks
│       └── rules
├── summarise
│   └── installation
└── transform
    └── devices
```

A full list of the commands and their options is described in the [docs](./docs/h7t.md).

## Transforms

The tool includes a plugin based solution for transforming customer data into a format that can be consumed. Further information is available in the [plugins](./plugins/) directory.

## Docker

An image is maintained on Docker Hub.

It can be pulled as follows, consider using a versioned release rather than latest:

```console
docker pull damianoneill/h7t:latest
```

And as an example, run the configure devices, with a sample.rpc in your current directory, as follows: 

```console
docker run -v "$(pwd):/config" damianoneill/h7t:v1.3.0 configure devices -i /config -f /config/sample.rpc
```
