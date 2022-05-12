package config

type SenderConfig struct {
	EventTopic          string   `envconfig:"SENDER_EVENT_TOPIC" default:"collector.services"`
	ParalelUploads      int      `envconfig:"SENDER_PARALLEL_UPLOADS" default:"2"`
	CollectDelay        int      `envconfig:"SENDER_COLLECT_DELAY" default:"5"`
	MatchPatterns       []string `envconfig:"SENDER_MATCH_PATTERNS" required:"true" default:"upload/*.json"`
	MaxCollectBatchSize int      `envconfig:"SENDER_MAX_COLLECT_BATCH_SIZE" default:"5"`
}
