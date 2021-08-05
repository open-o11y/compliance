package targets

import (
	"fmt"
	"os"
	"path"
)


func RunOtelCollector(opts TargetOptions) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	//hardcoded binary name
	binary := path.Join(cwd, "bin", "otelcol_linux_amd64")

	cfg := fmt.Sprintf(`
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: 'test'
          scrape_interval: 1s
          static_configs:
            - targets: [ '%s' ]
processors:
  batch:
exporters:
  prometheusremotewrite:
    endpoint: '%s'
service:
  pipelines:
    metrics:
      receivers: [prometheus]
      processors: [batch]
      exporters: [prometheusremotewrite]
`, opts.ScrapeTarget, opts.ReceiveEndpoint)
	configFileName, err := writeTempFile(cfg, "config-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(configFileName)

	return runCommand(binary, opts.Timeout, `--metrics-addr=:0`, fmt.Sprintf("--config=%s", configFileName))
}
