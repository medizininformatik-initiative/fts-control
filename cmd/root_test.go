package cmd

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {

	rootDir, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}

	err = os.Chdir(rootDir)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestConfigIsLoaded(t *testing.T) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		t.Fatalf("config.yaml konnte nicht geladen werden: %v", err)
	}
}
