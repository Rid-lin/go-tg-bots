package config

import (
	"os"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"

	log "github.com/sirupsen/logrus"
)

// var Location *time.Location
var ConfigFilePath string

type Config struct {
	Words         []string `default:"" usage:"words for search"`
	ChatsIDSearch []int64  `default:"" usage:"ID chats for Search"`
	ChatID        int64    `default:"" usage:"ID chats for U"`
	APIID         string   `default:"" usage:"APIID for api telegram app"`
	APIHash       string   `default:"" usage:"APIHash for api telegram app"`
	LogLevel      string   `default:"info" usage:"Log level: panic, fatal, error, warn, info, debug, trace"`
	PhoneNumber   string   `default:"" usage:"Number phone"`
	Password      string   `default:"" usage:"Password for, TG"`
}

func New() *Config {
	// fix for https://github.com/cristalhq/aconfig/issues/82
	args := []string{}
	for _, a := range os.Args {
		if !strings.HasPrefix(a, "-test.") {
			args = append(args, a)
		}
	}
	// fix for https://github.com/cristalhq/aconfig/issues/82

	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		// feel free to skip some steps :)
		// SkipEnv:      true,
		// MergeFiles: true,
		SkipFiles:          false,
		AllowUnknownFlags:  true,
		AllowUnknownFields: true,
		SkipDefaults:       false,
		SkipFlags:          false,
		FailOnFileNotFound: false,
		EnvPrefix:          "TGBOTFV",
		FlagPrefix:         "",
		FileFlag:           "config",
		Files: []string{
			"./config.yaml",
			"./config/config.yaml",
			"/etc/tgbotfv/config.yaml",
			"/etc/tgbotfv/config/config.yaml",
			"/usr/local/tgbotfv/config.yaml",
			"/usr/local/tgbotfv/config/config.yaml",
			"/opt/tgbotfv/config.yaml",
			"/opt/tgbotfv/config/config.yaml",
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			// from `aconfigyaml` submodule
			// see submodules in repo for more formats
			".yaml": aconfigyaml.New(),
		},
		Args: args[1:], // [1:] важно, см. доку к FlagSet.Parse
	})
	if err := loader.Load(); err != nil {
		log.Error(err)
	}

	return &cfg
}
