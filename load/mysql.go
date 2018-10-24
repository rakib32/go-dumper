package load

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-dumper/helper"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MySQL struct {
	Base
	host     string
	port     string
	database string
	username string
	password string
}

func (e *MySQL) loadConfig() {
	e.host = viper.GetString("destination.host")
	e.port = viper.GetString("destination.port")
	e.database = viper.GetString("destination.database")
	e.username = viper.GetString("destination.username")
	e.password = viper.GetString("destination.password")

	if e.CreateNewSchema {
		currTime := time.Now()
		e.database = e.database + "_" + currTime.Format("20060102")
	}
}

func (e *MySQL) Load() (err error) {
	e.loadConfig()

	if len(e.database) == 0 {
		err = errors.New("mysql database config is required")
		return err
	}

	if err == nil && e.Base.Storer != nil {
		downloadFilePath, err := e.Base.Storer.Download(e.Base.DumpFilePath)

		if err != nil {
			return err
		}

		e.Base.DumpFilePath = downloadFilePath
	}

	if e.CreateNewSchema {
		err = e.createSchema()

		if err != nil {
			return err
		}
	}

	err = e.load()

	return err
}

func (e *MySQL) getCommonOptions() []string {
	loadArgs := []string{}

	if len(e.host) > 0 {
		loadArgs = append(loadArgs, "--host", e.host)
	}
	if len(e.port) > 0 {
		loadArgs = append(loadArgs, "--port", e.port)
	}
	if len(e.username) > 0 {
		loadArgs = append(loadArgs, "-u", e.username)
	}
	if len(e.password) > 0 {
		loadArgs = append(loadArgs, `-p`+e.password)
	}

	return loadArgs
}

func (e *MySQL) loadOptions() []string {
	loadArgs := []string{}

	loadArgs = e.getCommonOptions()
	loadArgs = append(loadArgs, e.database)
	loadArgs = append(loadArgs, "-e")

	loadArgs = append(loadArgs, "SET FOREIGN_KEY_CHECKS = 0;SET UNIQUE_CHECKS = 0;source "+e.Base.DumpFilePath+";SET FOREIGN_KEY_CHECKS = 1;SET UNIQUE_CHECKS = 1;")

	return loadArgs
}

func (e *MySQL) updateOptions() []string {
	loadArgs := []string{}

	loadArgs = e.getCommonOptions()
	loadArgs = append(loadArgs, e.database)
	loadArgs = append(loadArgs, "-e")

	loadArgs = append(loadArgs, "UPDATE users set email=CONCAT(CONCAT('email', id), '@example.net'), number = id;")

	return loadArgs
}

func (e *MySQL) schemaCreateOptions() []string {
	loadArgs := []string{}

	loadArgs = e.getCommonOptions()
	loadArgs = append(loadArgs, "-e")

	loadArgs = append(loadArgs, "CREATE SCHEMA "+e.database+" DEFAULT CHARACTER SET utf8;")

	return loadArgs
}

func (e *MySQL) load() (err error) {

	logger.Info("--------- Loading MySQL---------------")
	out, err := helper.Exec("mysql", e.loadOptions()...)

	if err != nil {
		return fmt.Errorf("-> Load error: %s", err)
	}
	logger.Info("Load path:", e.Base.DumpFilePath)
	logger.Info("Output:", out)

	if err != nil {
		return fmt.Errorf("-> Tar error: %s", err)
	}

	return err
}

func (e *MySQL) Update() (err error) {
	e.loadConfig()

	if len(e.database) == 0 {
		err = errors.New("mysql database config is required")
		return err
	}

	err = e.update()

	return err
}

func (e *MySQL) update() (err error) {

	logger.Info("-------------Updating user table data---------------")
	out, err := helper.Exec("mysql", e.updateOptions()...)

	if err != nil {
		return fmt.Errorf("-> Update error: %s", err)
	}

	logger.Info("Output:", out)

	if err != nil {
		return fmt.Errorf("-> Tar error: %s", err)
	}

	logger.Info("--------- Updated User email and phone---------------")

	return err
}

func (e *MySQL) createSchema() (err error) {
	logger.Info("---------Creating schema----: ", e.database)
	out, err := helper.Exec("mysql", e.schemaCreateOptions()...)

	if err != nil {
		return fmt.Errorf("-> Create Schema error: %s", err)
	}

	logger.Info("Output:", out)

	if err != nil {
		return fmt.Errorf("-> Create schema error: %s", err)
	}

	logger.Info("--------- Created Schema Successfully---------------")

	return err
}
