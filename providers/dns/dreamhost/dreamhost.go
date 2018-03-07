// Package dreamhost implements a DNS provider for solving the DNS-01
// challenge using DreamHost DNS.
package dreamhost

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/xenolf/lego/acme"
)

var dreamhostAPIURL = "https://api.dreamhost.com"

type DNSProvider struct {
	apiKey string
}

func NewDNSProvider() (*DNSProvider, error) {
	apiKey := os.Getenv("DREAMHOST_API_KEY")

        return NewDNSProviderCredentials(apiKey)
}

func NewDNSProviderCredentials(apiKey string) (*DNSProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Dreamhost credentials missing")
	}

	return &DNSProvider{
		apiKey: apiKey,
	}, nil
}

func (c *DNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value, _ := acme.DNS01Record(domain, keyAuth)
	return c.addDnsRecord(fqdn, "TXT", value)
}

func (c *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	fqdn, value, _ := acme.DNS01Record(domain, keyAuth)
	return c.removeDnsRecord(fqdn, "TXT", value)
}

func (c *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return 60 * time.Minute, 30 * time.Second
}

func (c *DNSProvider) addDnsRecord(fqdn, t, value string) error {
	fqdn = strings.TrimRight(fqdn, ".")
	return c.apiRequest("dns-add_record", map[string]string{
		"record": fqdn,
		"type":   t,
		"value":  value,
	})
}

func (c *DNSProvider) removeDnsRecord(fqdn, t, value string) error {
	fqdn = strings.TrimRight(fqdn, ".")
	return c.apiRequest("dns-remove_record", map[string]string{
		"record": fqdn,
		"type":   t,
		"value":  value,
	})
}

func (c *DNSProvider) apiRequest(cmd string, args map[string]string) error {
	u, _ := url.Parse(dreamhostAPIURL)
	query := u.Query()
	query.Add("format", "json")
	query.Add("key", c.apiKey)
	query.Add("cmd", cmd)

	for k, v := range args {
		query.Add(k, v)
	}
	u.RawQuery = query.Encode()

	requestUrl := u.String()

	resp, _ := http.Get(requestUrl)

	decoder := json.NewDecoder(resp.Body)

	type DreamhostApiResponse struct {
		Result string
		Data   string
	}

	var dreamhostResponse DreamhostApiResponse
	decoder.Decode(&dreamhostResponse)

	if dreamhostResponse.Result == "success" {
		return nil
	}

	return errors.New(dreamhostResponse.Data)
}
