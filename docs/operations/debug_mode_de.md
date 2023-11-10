# Debug-Mode

Die Informationen zum Debug-Mode, wie zum Beispiel Log-Level, werden im Gegensatz zu `cesappd` anders gespeichert.
Für diese Informationen existiert eine Registry in der Configmap `debug-mode-registry`.
Diese wird beim Aktivieren erstellt und beim Deaktivieren wieder gelöscht.
Es wurde von der regulären CES-Registry im ETCD abgesehen, weil im Kubernetes-Kontext noch andere Values in die Registry
einfließen (z.B. Log-Level anderer Komponenten im Cluster).

## Werte, die in der Registry enthalten sind:

### enabled

* YAML-Key: `enabled`
* Typ: `bool`
* Notwendige Konfiguration
* Beschreibung: Gibt an, ob der Debug-Mode aktiviert oder deaktiviert ist.
* Beispiel: `true`

### disable-at-timestamp

* YAML-Key: `disable-at-timestamp`
* Typ: `string`
* Notwendige Konfiguration
* Beschreibung: Gibt an, wann der Debug-Mode automatisch wieder deaktiviert wird.
* Beispiel: `10 Nov 23 10:48 UTC`

> **Hinweis:** Der String ist ein Zeitstempel formatiert nach RFC822. Es wird in regelmäßigen Takt dieser Key abgefragt und entschieden, ob der Debug-Mode deaktiviert werden kann.

### dogus

* YAML-Key: `dogus`
* Typ: `map[string]string`
* Notwendige Konfiguration
* Beschreibung: Gibt die Log-Level an, die die Dogus vor der Aktivierung des Debug-Modes hatten.
* Beispiel:
```yaml
ldap: "ERROR"
cas: ""
postfix: "INFO"
```

> **Hinweis:** Hatte ein Dogu keinen expliziten Wert als Log-Level konfiguriert wird ein leerer String `""` als Wert verwendet.