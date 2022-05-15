package app

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
)

func ReadWarningsToStdErr(err chan error) {
	go func() {
		for {
			pp.Fprintln(os.Stderr, <-err)
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
		dry.PanicIfErr(err)
	}
	publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")
	if !dry.FileExists(publicKeys) {
		return "", fmt.Errorf("file tg_public_keys.pem in ./config not found")
	}

	return appStorage, nil
}

type Namespace uint8

const (
	NamespaceUnknown Namespace = iota
	NamespaceGlobal
	NamespaceUser
	NamespaceDirectory
)

func GetAppStorage(appName string, namespace Namespace) (string, error) {
	switch namespace {
	case NamespaceGlobal:
		return filepath.Join("var", "lib", appName), nil
	case NamespaceUser:
		p, err := GetAppStorage(appName, NamespaceGlobal)
		if err != nil {
			return "", err
		}
		u, _ := user.Current()
		userPath, err := GetUserNamespaceDir(u.Username)
		if err != nil {
			return "", err
		}
		return filepath.Join(userPath, p), nil
	default:
		return "", errors.New("Incompatible feature for this namespace")
	}
}

func GetUserNamespaceDir(username string) (string, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return "", errors.Wrapf(err, "looking up '%v'", username)
	}

	return filepath.Join(u.HomeDir, ".local"), nil
}
