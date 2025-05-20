package config

import "github.com/spf13/viper"

type Config struct {
	BotToken       string `mapstructure:"LTA_BOT_TOKEN"`
	BotTimeout     int    `mapstructure:"LTA_BOT_TIMEOUT"`
	DataPath       string `mapstructure:"DATA_PATH"`
	DeviceAddr     string `mapstructure:"DEVICE_ADDR"`
	DevicePort     int    `mapstructure:"DEVICE_PORT"`
	PollIntervalMs int    `mapstructure:"POLL_INTERVAL_MS"`
}

func LoadConfig(name, ext string, paths ...string) (Config, error) {
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigName(name)
	viper.SetConfigType(ext)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	config := Config{}

	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
