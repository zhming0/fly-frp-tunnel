package frp

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerateClientConfig(t *testing.T) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "envoy-gateway",
			Namespace: "envoy-gateway-system",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     80,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "https",
					Port:     443,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}

	config := GenerateClientConfig(svc, "137.66.1.1", 7000)

	expected := `serverAddr = "137.66.1.1"
serverPort = 7000

[[proxies]]
name = "envoy-gateway-http"
type = "tcp"
localIP = "envoy-gateway.envoy-gateway-system.svc.cluster.local"
localPort = 80
remotePort = 80

[[proxies]]
name = "envoy-gateway-https"
type = "tcp"
localIP = "envoy-gateway.envoy-gateway-system.svc.cluster.local"
localPort = 443
remotePort = 443

`

	if config != expected {
		t.Errorf("unexpected config:\ngot:\n%s\nwant:\n%s", config, expected)
	}
}

func TestGenerateClientConfigUnnamedPort(t *testing.T) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minecraft",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     25565,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}

	config := GenerateClientConfig(svc, "10.0.0.1", 7000)

	if !contains(config, `name = "minecraft-25565"`) {
		t.Errorf("expected proxy name 'minecraft-25565' in config:\n%s", config)
	}
}

func TestGenerateServerConfig(t *testing.T) {
	config := GenerateServerConfig(7000)
	expected := "bindPort = 7000\n"
	if config != expected {
		t.Errorf("unexpected server config: got %q, want %q", config, expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
