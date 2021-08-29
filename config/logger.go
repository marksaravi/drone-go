package config

type UdpLoggerConfigs struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"packets-per-second"`
	MaxDataPerPacket int    `yaml:"max-data-per-packet"`
}

type udpLoggerConfigs struct {
	UdpLoggerConfigs UdpLoggerConfigs `yaml:"logger"`
}

func ReadLoggerConfig() udpLoggerConfigs {
	return readConfig(udpLoggerConfigs{}).(udpLoggerConfigs)
}
