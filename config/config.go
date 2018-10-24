package config

import (
	"errors"
	"fmt"

	logrus "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// App holds app configuration
type App struct {
	Self       Self
	consulURL  string
	consulPath string
}

// Init sets consule url and path
func (a *App) Init() error {
	viper.SetEnvPrefix("app")
	viper.BindEnv("env")
	viper.BindEnv("consul_url")
	viper.BindEnv("consul_path")

	a.consulURL = viper.GetString("consul_url")
	a.consulPath = viper.GetString("consul_path")
	if a.consulURL == "" {
		return errors.New("CONSUL_URL missing")
	}
	if a.consulPath == "" {
		return errors.New("CONSUL_PATH missing")
	}

	return nil
}

// ReadConsule reads config from consule
func (a *App) ReadConsule() error {
	viper.AddRemoteProvider("consul", a.consulURL, a.consulPath)
	viper.SetConfigType("yml")
	err := viper.ReadRemoteConfig()
	if err != nil {
		return fmt.Errorf("%s named \"%s\"", err.Error(), a.consulPath)
	}

	return nil
}

// Load loads app cofig
func (a *App) Load() {
	a.Self.Load()
}

// Prt prints all configuration
func (a App) Prt() {
	a.Self.Prt()
}

// Mock loads mock values
func (a *App) Mock() {
	a.Self.Mock()
}

// Self holds self configuration
type Self struct {
	Port string
}

// Load loads self configuration
func (s *Self) Load() {
	// s.Port = viper.GetString("port") // import from consul
	s.Port = "4200"
}

// Mock loads mock values
func (s *Self) Mock() {
	s.Port = "4201"
}

// Prt prints all configuration
func (s Self) Prt() {
	logrus.Info("self port", s.Port)
}

// New returns App
func New() *App {
	return &App{}
}
