# Transforms

Customer data is available in lots of different formats. The tool provides a plugin loading mechanism for transforming that data into the dsl format used by the tool.

This provides an extension mechanism for writing customer plugins for doing the conversion.

The plugin technology is based on hashicorp [go-plugin](https://github.com/hashicorp/go-plugin) a RPC based solution for Go plugins.

The solution ships with a sample plugin for devices that transforms files containing records that are delimited by a definable character set for e.g. comma's, tab's, etc.

## CSV Transform

The sample plugin provides a transform for mapping CSV files into the Device dsl required by Healthbot.

The plugin expects to parse CSV files in the following format:

```csv
device-id,host,username,password
mx1,10.0.0.1,root,changeme
mx2,10.0.0.2,root,"changeme now"
mx3,10.0.0.3,,
```