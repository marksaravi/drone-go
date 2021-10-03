package config

type udpLoggerConfigs struct {
	UdpLoggerConfigs struct {
		Enabled          bool   `yaml:"enabled"`
		IP               string `yaml:"ip"`
		Port             int    `yaml:"port"`
		PacketsPerSecond int    `yaml:"packets-per-second"`
	} `yaml:"logger"`
}

func ReadLoggerConfig() udpLoggerConfigs {
	return readConfig(udpLoggerConfigs{}).(udpLoggerConfigs)
}
