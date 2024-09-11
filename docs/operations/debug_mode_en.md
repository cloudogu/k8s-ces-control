# Debug mode

In contrast to `cesappd`, the information on debug mode, such as log level, is stored differently.
A registry exists for this information in the `debug-mode-registry` configmap.
This is created when activated and deleted again when deactivated.

## Values that are contained in the registry:

### enabled

* YAML key: `enabled`
* Type: `bool`
* Necessary configuration
* Description: Specifies whether debug mode is enabled or disabled.
* Example: `true`

### disable-at-timestamp

* YAML key: `disable-at-timestamp`.
* Type: `string`
* Necessary configuration
* Description: Specifies when the debug mode is automatically deactivated.
* Example: `10 Nov 23 10:48 UTC`

**Note:** The string is a timestamp formatted according to RFC822. This key is queried at regular intervals and a decision is made whether debug mode can be deactivated.

### dogus

* YAML key: `dogus`
* Type: `map[string]string`.
* Necessary configuration
* Description: Specifies the log level that the dogus had before the debug mode was activated.
* Example:
```yaml
ldap: "ERROR"
cas: ""
postfix: "INFO"
```

> **Note:** If a dogu had no explicit value configured as log level, an empty string `""` is used as the value.