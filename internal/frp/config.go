// Package frp generates frpc configuration from Kubernetes Service specs.
package frp

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	// DefaultServerPort is the default frps control port.
	DefaultServerPort = 7000
)

// GenerateClientConfig generates a TOML frpc configuration from a Service spec.
// serverAddr is the fly.io Machine's dedicated IPv4 address.
func GenerateClientConfig(svc *corev1.Service, serverAddr string, serverPort int) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("serverAddr = \"%s\"\n", serverAddr))
	b.WriteString(fmt.Sprintf("serverPort = %d\n", serverPort))
	b.WriteString("\n")

	// Build the ClusterIP DNS name for this service.
	localIP := fmt.Sprintf("%s.%s.svc.cluster.local", svc.Name, svc.Namespace)

	for _, port := range svc.Spec.Ports {
		proxyName := fmt.Sprintf("%s-%s", svc.Name, port.Name)
		if port.Name == "" {
			proxyName = fmt.Sprintf("%s-%d", svc.Name, port.Port)
		}

		protocol := strings.ToLower(string(port.Protocol))
		if protocol == "" {
			protocol = "tcp"
		}

		b.WriteString("[[proxies]]\n")
		b.WriteString(fmt.Sprintf("name = \"%s\"\n", proxyName))
		b.WriteString(fmt.Sprintf("type = \"%s\"\n", protocol))
		b.WriteString(fmt.Sprintf("localIP = \"%s\"\n", localIP))
		b.WriteString(fmt.Sprintf("localPort = %d\n", port.Port))
		b.WriteString(fmt.Sprintf("remotePort = %d\n", port.Port))
		b.WriteString("\n")
	}

	return b.String()
}

// GenerateServerConfig generates a minimal TOML frps configuration.
func GenerateServerConfig(bindPort int) string {
	return fmt.Sprintf("bindPort = %d\n", bindPort)
}
