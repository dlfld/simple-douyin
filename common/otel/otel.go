package otel

import (
	"github.com/douyin/common/conf"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func NewOtelProvider(serviceName string) provider.OtelProvider {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint(conf.OtelConfig.Addr),
		provider.WithInsecure(),
	)
	return p
}
