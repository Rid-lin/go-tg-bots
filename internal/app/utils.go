package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/xelaj/go-dry"
)

func ReadWarningsToStdErr(err chan error) {
	go func() {
		for {
			fmt.Fprintf(os.Stderr, "%v\n", <-err)
		}
	}()
}

func PrepareAppStorage(appStorage string) (appStoragePath string, err error) {

	if !dry.FileExists(appStorage) {
		if !dry.PathIsWritable(appStorage) {
			fmt.Printf("cant create app local storage at %v\n", appStorage)
			os.Exit(1)
		}
		err := os.MkdirAll(appStorage, 0755)
		ExitIfErr(err)
	}
	publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")
	if !dry.FileExists(publicKeys) {
		return "", fmt.Errorf("file tg_public_keys.pem in ./config not found")
	}

	return appStorage, nil
}

func ExitIfErr(args ...any) {
	flagS := false
	for _, v := range args {
		if err, _ := v.(error); err != nil {
			flagS = true
			logrus.Errorf("%+v\n", err)
		}
	}
	if flagS {
		os.Exit(1)
	}
}
