# Transforms

Customer data is available in lots of different formats. The tool provides a plugin loading mechanism for transforming that data into the dsl format used by the tool.

This provides an extension mechanism for writing customer plugins for doing the conversion.

The plugin technology is based on hashicorp [go-plugin](https://github.com/hashicorp/go-plugin) a RPC based solution for Go plugins.

The solution ships with a sample plugin for devices that transforms files containing records that are delimited by a definable character set for e.g. comma's, tab's, etc.

## Delimited Transform

The sample plugin provides a generic solution for mapping delimited files into the Device dsl required by Healthbot.

The plugin takes as its first argument a regex for the delimiter e.g. 

```console
h7t transform devices --plugin delimited '[ ]{1,}'
```

This regex will split records into fields that are delimited on one or more spaces.  For additional regex [examples](https://yourbasic.org/golang/regexp-cheat-sheet/). 

Additional arguments are required to instruct the plugin on how to map columns in the delimited file into fields in the Devices dsl. 
JSON filter notation from [jq](https://stedolan.github.io/jq/) is used to define this mapping. 
As an example, this [jq playground snippet](https://jqplay.org/s/NcDU2mcijx) shows an example of defining a filter to label the password used for NETCONF authentication. 

> Devices dsl contains json keys that include hypens, jq requires any filters on special characters to be wrapped in "", for e.g. on the playground to get the device-id you would need to use ."device-id", for our purposes we will ignore this requirement. 

```console
h7t transform devices --plugin delimited '[ ]{1,}' '1:.host' '2:.device-id' '5:.authentication.password.username' '6:.authentication.password.password'
```

As you can guess, the number indicates the column and the jq filter indicates the content

An example of a compatible file would look like below, note the Site and Rack columns are ignored

```console
10.0.0.1   pe1   "Site a"   "Rack 1"   root   changeme
10.0.0.2   pe2   "Site a"   "Rack 2"   root   changeme
10.0.0.3   pe3   "Site a"   "Rack 2"   root   changeme
10.0.0.4   pe4   "Site a"   "Rack 3"   root   changeme
```