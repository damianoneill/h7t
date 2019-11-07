# h~~ealthbot~~<sup>7</sup>t

[![GitHub release](https://img.shields.io/github/v/release/damianoneill/h7t.svg)](https://GitHub.com/damianoneill/h7t/releases/)
[![Go Report Card](https://goreportcard.com/badge/damianoneill/h7t)](http://goreportcard.com/report/damianoneill/h7t)
[![license](https://img.shields.io/github/license/damianoneill/h7t.svg)](https://github.com/damianoneill/h7t/blob/master/LICENSE)

A command line tool for interacting with [Juniper Healthbot](https://www.juniper.net/us/en/products-services/sdn/contrail/contrail-healthbot/).

## Synopsis

A tool for interacting with Healthbot over the REST API.

The initial focus of this tool is to provide bulk or aggregate functions, that simplify interacting with Healthbot.

## Commands

```console
h7t
├── configure
│   └── devices
├── extract
│   ├── device-groups
│   └── devices
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
