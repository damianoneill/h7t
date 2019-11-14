# Contributing to h7t

## Setup

Assuming you have a go setup on your host, run make, which will trigger a download of the tools required to write code in the repo.

## Command Structure

In h7t commands represent actions/verbs, args are things/nouns/objects and Flags are modifiers/adjective (property/state) for actions.

For example:

| action     | thing/object                        | property/state |
| ---------- | ----------------------------------- | -------------- |
| myapp verb | noun                                | --adjective    |
| git clone  | git@github.com:damianoneill/h7t.git | --bare         |
| go get -u  | github.com/aws/aws-sdk-go/...       |                |

If adding new commands, understand how they will fit into the bigger picture by reviewing below.

### Actions

The following are the list of actions supported / planned for h7t

- extract
- transform
- load
- summarise
- configure

## Things

The following are the list objects supported / planned by h7t, this includes Healthbot DSL objects + others for e.g. Junos Devices

- device
- devices
- device-groups
- helper-files
- installation

### States

The following are the list of properties supported / planned by h7t

#### Common

- verbose (v)
- config
- help (h)
- authority (a)
- username (u)
- password (p)

#### Action specific

- N/A

#### Thing specific

- output_directory (o)
- input_directory (i)
- erase (e)
- plugin
- netconf-rpc (f)
