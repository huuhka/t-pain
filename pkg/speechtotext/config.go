package speechtotext

type Config struct {
	Key    string
	Region string
}

func NewConfig(key string, region string) *Config {
	return &Config{
		Key:    key,
		Region: region,
	}
}