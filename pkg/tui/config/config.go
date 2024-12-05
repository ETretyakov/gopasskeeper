package config

type AppConfig struct {
	Fullscreen  bool
	EnableMouse bool
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Fullscreen:  true,
		EnableMouse: true,
	}
}
