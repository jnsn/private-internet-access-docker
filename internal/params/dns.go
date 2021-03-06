package params

import (
	"fmt"
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetDNSOverTLS obtains if the DNS over TLS should be enabled
// from the environment variable DOT
func (p *paramsReader) GetDNSOverTLS() (DNSOverTLS bool, err error) {
	return p.envParams.GetOnOff("DOT", libparams.Default("on"))
}

// GetDNSOverTLSProviders obtains the DNS over TLS providers to use
// from the environment variable DOT_PROVIDERS
func (p *paramsReader) GetDNSOverTLSProviders() (providers []models.DNSProvider, err error) {
	s, err := p.envParams.GetEnv("DOT_PROVIDERS", libparams.Default("cloudflare"))
	if err != nil {
		return nil, err
	}
	for _, word := range strings.Split(s, ",") {
		provider := models.DNSProvider(word)
		switch provider {
		case constants.Cloudflare, constants.Google, constants.Quad9, constants.Quadrant, constants.CleanBrowsing, constants.SecureDNS, constants.LibreDNS:
			providers = append(providers, provider)
		default:
			return nil, fmt.Errorf("DNS over TLS provider %q is not valid", provider)
		}
	}
	return providers, nil
}

// GetDNSOverTLSVerbosity obtains the verbosity level to use for Unbound
// from the environment variable DOT_VERBOSITY
func (p *paramsReader) GetDNSOverTLSVerbosity() (verbosityLevel uint8, err error) {
	n, err := p.envParams.GetEnvIntRange("DOT_VERBOSITY", 0, 5, libparams.Default("1"))
	return uint8(n), err
}

// GetDNSOverTLSVerbosityDetails obtains the log level to use for Unbound
// from the environment variable DOT_VERBOSITY_DETAILS
func (p *paramsReader) GetDNSOverTLSVerbosityDetails() (verbosityDetailsLevel uint8, err error) {
	n, err := p.envParams.GetEnvIntRange("DOT_VERBOSITY_DETAILS", 0, 4, libparams.Default("0"))
	return uint8(n), err
}

// GetDNSOverTLSValidationLogLevel obtains the log level to use for Unbound DOT validation
// from the environment variable DOT_VALIDATION_LOGLEVEL
func (p *paramsReader) GetDNSOverTLSValidationLogLevel() (validationLogLevel uint8, err error) {
	n, err := p.envParams.GetEnvIntRange("DOT_VALIDATION_LOGLEVEL", 0, 2, libparams.Default("0"))
	return uint8(n), err
}

// GetDNSMaliciousBlocking obtains if malicious hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_MALICIOUS
func (p *paramsReader) GetDNSMaliciousBlocking() (blocking bool, err error) {
	return p.envParams.GetOnOff("BLOCK_MALICIOUS", libparams.Default("on"))
}

// GetDNSSurveillanceBlocking obtains if surveillance hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_SURVEILLANCE
// and BLOCK_NSA for retrocompatibility
func (p *paramsReader) GetDNSSurveillanceBlocking() (blocking bool, err error) {
	// Retro-compatibility
	s, err := p.envParams.GetEnv("BLOCK_NSA")
	if err != nil {
		return false, err
	} else if len(s) != 0 {
		p.logger.Warn("You are using the old environment variable BLOCK_NSA, please consider changing it to BLOCK_SURVEILLANCE")
		return p.envParams.GetOnOff("BLOCK_NSA", libparams.Compulsory())
	}
	return p.envParams.GetOnOff("BLOCK_SURVEILLANCE", libparams.Default("off"))
}

// GetDNSAdsBlocking obtains if ads hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_ADS
func (p *paramsReader) GetDNSAdsBlocking() (blocking bool, err error) {
	return p.envParams.GetOnOff("BLOCK_ADS", libparams.Default("off"))
}

// GetDNSUnblockedHostnames obtains a list of hostnames to unblock from block lists
// from the comma separated list for the environment variable UNBLOCK
func (p *paramsReader) GetDNSUnblockedHostnames() (hostnames []string, err error) {
	s, err := p.envParams.GetEnv("UNBLOCK")
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, nil
	}
	hostnames = strings.Split(s, ",")
	for _, hostname := range hostnames {
		if !p.verifier.MatchHostname(hostname) {
			return nil, fmt.Errorf("hostname %q does not seem valid", hostname)
		}
	}
	return hostnames, nil
}

// GetDNSOverTLSCaching obtains if Unbound caching should be enable or not
// from the environment variable DOT_CACHING
func (p *paramsReader) GetDNSOverTLSCaching() (caching bool, err error) {
	return p.envParams.GetOnOff("DOT_CACHING")
}

// GetDNSOverTLSPrivateAddresses obtains if Unbound caching should be enable or not
// from the environment variable DOT_PRIVATE_ADDRESS
func (p *paramsReader) GetDNSOverTLSPrivateAddresses() (privateAddresses []string) {
	s, _ := p.envParams.GetEnv("DOT_PRIVATE_ADDRESS")
	for _, s := range strings.Split(s, ",") {
		privateAddresses = append(privateAddresses, s)
	}
	return privateAddresses
}

// GetDNSOverTLSIPv6 obtains if Unbound should resolve ipv6 addresses using ipv6 DNS over TLS
//  servers from the environment variable DOT_IPV6
func (p *paramsReader) GetDNSOverTLSIPv6() (ipv6 bool, err error) {
	return p.envParams.GetOnOff("DOT_IPV6")
}
