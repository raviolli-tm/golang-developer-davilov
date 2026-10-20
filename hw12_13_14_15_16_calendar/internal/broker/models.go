package broker

type KafkaConf struct {
	BrokersUrl string `yaml:"brokers_url"`
	Topic      string `yaml:"topic"`
	Interval   int    `yaml:"interval"`
	GroupId    string `yaml:"group_id"`
}

type RetryPolicyConf struct {
	MaxRetries int `yaml:"max_retries"`
	BaseTime   int `yaml:"base_time_seconds"`
}
