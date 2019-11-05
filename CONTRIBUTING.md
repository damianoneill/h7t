# Contributing to h7t

In h7t commands represent actions/verbs, args are things/nouns/objects and Flags are modifiers/adjective (property/state) for actions.

For example:

| action     | thing/object                        | property/state |
| ---------- | ----------------------------------- | -------------- |
| myapp verb | noun                                | --adjective    |
| git clone  | git@github.com:damianoneill/h7t.git | --bare         |
| go get -u  | github.com/aws/aws-sdk-go/...       |                |

If adding new commands, understand how they will fit into the bigger picture by reviewing below.

## Actions

The following are the list of actions supported / planned for h7t

- provision
- summarize
- replicate

## Things

The following are the list objects supported / planned by h7t

- devices
- device-groups
- helper-files
- playbook-instances
- installation

## States

The following are the list of properties supported / planned by h7t

### Common

- verbose (v)
- config
- help (h)
- authority (a)
- username (u)
- password (p)

### Action specific

- N/A

### Thing specific

- directory (d)
- erase (e)
