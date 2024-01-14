package latency

import "time"

// Config defines setting for moving average latency calculations
type Config struct {
	Decay          float32        `yaml:"decay" json:"decay"`                                                              // Weight of new latency measurements
	WarmupSamples  int            `yaml:"warmup_samples" json:"warmup_samples"`                                            // The number of latency probes required to init moving average
	UpdateInterval *time.Duration `yaml:"update_interval,omitempty" json:"update_interval" swaggertype:"primitive,string"` // How often gateway should probe models with not the lowest response latency
}

func DefaultConfig() *Config {
	defaultUpdateInterval := 30 * time.Second

	return &Config{
		Decay:          0.06,
		WarmupSamples:  3,
		UpdateInterval: &defaultUpdateInterval,
	}
}
