package telemetry

type Config struct {
	OTLPEndpoint     string
	TraceSampleRate  float64
	MetricsListen    string
}

func (c *Config) MetricsAddr() string {
	if c.MetricsListen == "" {
		return "0.0.0.0:9091"
	}
	return c.MetricsListen
}