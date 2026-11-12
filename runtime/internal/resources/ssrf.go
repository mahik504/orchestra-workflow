package resources

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// RejectPrivateURL refuses localhost, loopback, private, link-local, and the cloud metadata address.
// A DNS failure is not a private address. The caller still has to handle the failed request.
func RejectPrivateURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("malformed URL")
	}
	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("missing host")
	}
	if strings.EqualFold(host, "localhost") {
		return fmt.Errorf("localhost blocked")
	}
	if ip := net.ParseIP(host); ip != nil {
		if privateIP(ip) {
			return fmt.Errorf("blocked IP %s", ip)
		}
		return nil
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil
	}
	for _, ip := range ips {
		if privateIP(ip) {
			return fmt.Errorf("resolved to blocked IP %s", ip)
		}
	}
	return nil
}

func privateIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.String() == "169.254.169.254"
}
