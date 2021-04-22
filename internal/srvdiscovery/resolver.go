package srvdiscovery

import (
	discovery "github.com/Confialink/wallet-pkg-discovery/v2"
	"net/url"
)

// Resolver returns default SD resolver
func Resolver() discovery.Resolver {
	topResolver := discovery.NewFallbackResolver(
		discovery.NewEnvResolver(),
		discovery.NewDNSResolver(&discovery.NetSRVResolver{}, "tcp"),
	)
	portMapping := map[string]string{
		PortNameRpc: "http",
	}

	return discovery.NewSchemeDecorator(topResolver, portMapping)
}

// Resolve performs service discovery by port name and service name
func Resolve(portName, serviceName string) (*url.URL, error) {
	return Resolver().Resolve(portName, serviceName)
}

// ResolveRPC performs service discovery by service name where port name is "rpc"
func ResolveRPC(serviceName string) (*url.URL, error) {
	return Resolver().Resolve(PortNameRpc, serviceName)
}
