package ssl

import "github.com/cloudogu/cesapp-lib/registry"

type configurationContext interface {
	registry.ConfigurationContext
}

type sslGenerator interface {
	GenerateSelfSignedCert(fqdn string, domain string, certExpireDays int, country string,
		province string, locality string, altDNSNames []string) (string, string, error)
}
