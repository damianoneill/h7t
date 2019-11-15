# h~~ealthbot~~<sup>7</sup>t

[![GitHub release](https://img.shields.io/github/v/release/damianoneill/h7t.svg)](https://GitHub.com/damianoneill/h7t/releases/)
[![Go Report Card](https://goreportcard.com/badge/damianoneill/h7t)](http://goreportcard.com/report/damianoneill/h7t)
[![license](https://img.shields.io/github/license/damianoneill/h7t.svg)](https://github.com/damianoneill/h7t/blob/master/LICENSE)

A command line tool for interacting with [Juniper Healthbot](https://www.juniper.net/us/en/products-services/sdn/contrail/contrail-healthbot/).

> Note Device Group functionality is limited currently, this will be updated after Healthbot 2.1 is released.

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
│   └── helper-files
├── summarise
│   └── installation
└── transform
    └── devices
```

A full list of the commands and their options is described in the [docs](./docs/h7t.md).

## Transforms

The tool includes a plugin based solution for transforming customer data into a format that can be consumed. Further information is available in the [plugins](./plugins/) directory.
