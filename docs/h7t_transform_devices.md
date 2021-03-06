## h7t transform devices

Transform Devices configuration

### Synopsis

Transform Devices configurations from comma separated value (csv) format into the dsl format using a bundled plugin.

```
h7t transform devices [flags]
```

### Options

```
  -h, --help   help for devices
```

### Options inherited from parent commands

```
  -a, --authority string          healthbot HTTPS Authority (default "localhost:8080")
      --config string             config file (default is $HOME/.h7t.yaml)
  -i, --input_directory string    directory where the configuration will be loaded from (default ".")
  -o, --output_directory string   directory where the configuration will be written to (default ".")
  -p, --password string           healthbot Password (default "****")
      --plugin string             name of the plugin to be used (default "csv")
  -u, --username string           healthbot Username (default "admin")
  -v, --verbose                   cause h7t to be more verbose
```

### SEE ALSO

* [h7t transform](h7t_transform.md)	 - Transform things from proprietary formats into Healthbot dsl format

###### Auto generated by spf13/cobra on 29-Nov-2019
