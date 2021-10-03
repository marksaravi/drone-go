package config

type joystick struct {
	Channel   int     `yaml:"channel"`
	ZeroValue float32 `yaml:"zero-value"`
}

type remoteControlConfigs struct {
	RemoteControlConfigs struct {
		Joysticks struct {
			Roll     joystick  `yaml:"roll"`
			Pitch    joystick  `yaml:"pitch"`
			Yaw      joystick  `yaml:"yaw"`
			Throttle joystick  `yaml:"throttle"`
			VRef     float32   `yaml:"v-ref"`
			SPI      spiConfig `yaml:"spi"`
		} `yaml:"joysticks"`
		Buttons struct {
			FrontLeft   string `yaml:"front-left"`
			FrontRight  string `yaml:"front-right"`
			TopLeft     string `yaml:"top-left"`
			TopRight    string `yaml:"top-right"`
			BottomLeft  string `yaml:"bottom-left"`
			BottomRight string `yaml:"bottom-right"`
		} `yaml:"buttons"`
	} `yaml:"remote-control"`
	Radio radioConnection `yaml:"radio-connetion"`
}

func ReadRemoteControlConfig() remoteControlConfigs {
	return readConfig(remoteControlConfigs{}).(remoteControlConfigs)
}
