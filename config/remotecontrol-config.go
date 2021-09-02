package config

type Joystick struct {
	Channel   int     `yaml:"channel"`
	ZeroValue float32 `yaml:"zero-value"`
}

type Joysticks struct {
	Roll     Joystick `yaml:"roll"`
	Pitch    Joystick `yaml:"pitch"`
	Yaw      Joystick `yaml:"yaw"`
	Throttle Joystick `yaml:"throttle"`
	VRef     float32  `yaml:"v-ref"`
	SPI      SPI      `yaml:"spi"`
}

type Buttons struct {
	FrontLeft   string `yaml:"front-left"`
	FrontRight  string `yaml:"front-right"`
	TopLeft     string `yaml:"top-left"`
	TopRight    string `yaml:"top-right"`
	BottomLeft  string `yaml:"bottom-left"`
	BottomRight string `yaml:"bottom-right"`
}

type RemoteControlConfigs struct {
	Joysticks Joysticks `yaml:"joysticks"`
	Buttons   Buttons   `yaml:"buttons"`
}

type remoteControlConfigs struct {
	RemoteControlConfigs RemoteControlConfigs `yaml:"remote-control"`
}

func ReadRemoteControlConfig() remoteControlConfigs {
	return readConfig(remoteControlConfigs{}).(remoteControlConfigs)
}
