# User authkey

This plugin is used to set the SSH login user and public key using environment variables `${CVM_USER}` and `${CVM_AUTH_KEY}`.

The default value of `${CVM_USER}` is "cvm", and users can customize it as shown below.
```
export CVM_USER=<user>
```

The `${CVM_AUTH_KEY}` has no default value, users need to set it themselves. If `${CVM_AUTH_KEY}` is not specified like below, this plugin will be skipped.

```
export CVM_AUTH_KEY=<ssh public key>
```