package config

type joystick struct {
	Channel   int     `yaml:"channel"`
	ZeroValue float32 `yaml:"zero-value"`
}

type remoteControlConfigs struct {
	Configs struct {
		CommandPerSecond int `yaml:"command-per-sec"`
		Joysticks        struct {
			Roll     joystick  `yaml:"roll"`
			Pitch    joystick  `yaml:"pitch"`
			Yaw      joystick  `yaml:"yaw"`
			Throttle joystick  `yaml:"throttle"`
			VRef     float32   `yaml:"v-ref"`
			SPI      SpiConfig `yaml:"spi"`
		} `yaml:"joysticks"`
		Buttons struct {
			FrontLeft   string `yaml:"front-left"`
			FrontRight  string `yaml:"front-right"`
			TopLeft     string `yaml:"top-left"`
			TopRight    string `yaml:"top-right"`
			BottomLeft  string `yaml:"bottom-left"`
			BottomRight string `yaml:"bottom-right"`
		} `yaml:"buttons"`
		Radio RadioConfig `yaml:"radio"`
	} `yaml:"remote-control"`
	RadioConnection RadioConnection `yaml:"radio-connection"`
}

func ReadRemoteControlConfig() remoteControlConfigs {
	return readConfig(remoteControlConfigs{}).(remoteControlConfigs)
}
