package util

import (
	"github.com/kmou424/ero"
	"net"
	"strings"
)

func GetSLD(domain string) (string, error) {
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return "", ero.Newf("invalid domain: %s", domain)
	}
	domainParts = domainParts[len(domainParts)-2:]
	return strings.Join(domainParts, "."), nil
}

func IsIPv4(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

func IsIPv6(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}

func JoinDomains(domains ...string) string {
	var sb strings.Builder
	for i := range domains {
		if strings.HasPrefix(domains[i], ".") {
			domains[i] = strings.TrimPrefix(domains[i], ".")
		}
		if strings.HasSuffix(domains[i], ".") {
			domains[i] = strings.TrimSuffix(domains[i], ".")
		}
		sb.WriteString(domains[i])
		if i < len(domains)-1 {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

func SafeDomain(domain string) string {
	var sb strings.Builder
	for _, c := range domain {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' {
			sb.WriteRune(c)
		}
		if c >= 'A' && c <= 'Z' {
			sb.WriteRune(c + 32)
		}
		if c == ' ' {
			sb.WriteRune('-')
		}
	}
	return sb.String()
}
