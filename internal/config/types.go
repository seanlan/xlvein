package config

type Config struct {
	Debug        bool           `yaml:"debug"`
	App          string         `yaml:"app"`
	Host         string         `yaml:"host"`
	Exchange     Exchange       `yaml:"exchange"`
	Applications []Applications `yaml:"applications"`
}

type Exchange struct {
	Type         string `yaml:"type"`
	Rabbitmq     string `yaml:"rabbitmq"`
	ExchangeName string `yaml:"exchange_name"`
	QueueName    string `yaml:"queue_name"`
	Redis        string `yaml:"redis"`
}

type Applications struct {
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
}

func (c *Config) GetApp(appKey string) *Applications {
	for _, v := range c.Applications {
		if v.AppKey == appKey {
			return &v
		}
	}
	return nil
}
