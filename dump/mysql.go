package dump

import (
	"errors"
	"fmt"
	"os"
	"path"
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

func (e *MySQL) Dump() (result *DumpResult, err error) {
	e.host = viper.GetString("src.host")
	e.port = viper.GetString("src.port")
	e.database = viper.GetString("src.database")
	e.username = viper.GetString("src.username")
	e.password = viper.GetString("src.password")

	if len(e.database) == 0 {
		err = errors.New("mysql database config is required")
		return result, err
	}

	e.Base.FileName = fmt.Sprintf(`backup_%v_%v.sql`, e.database, time.Now().Unix())
	e.Base.DumpFilePath = path.Join(e.Base.StorePath, e.Base.FileName)

	result, err = e.dump()

	return result, err
}

func (e *MySQL) dumpOptions() []string {
	dumpArgs := []string{}

	dumpArgs = append(dumpArgs, "--single-transaction=TRUE")
	dumpArgs = append(dumpArgs, "--column-statistics=0")
	dumpArgs = append(dumpArgs, "--set-gtid-purged=OFF")
	dumpArgs = append(dumpArgs, "--no-autocommit")
	//dumpArgs = append(dumpArgs, "--ignore-table="+e.database+".details")
	//dumpArgs = append(dumpArgs, "--ignore-table="+e.database+".large_table2")

	if len(e.host) > 0 {
		dumpArgs = append(dumpArgs, "--host", e.host)
	}
	if len(e.port) > 0 {
		dumpArgs = append(dumpArgs, "--port", e.port)
	}
	if len(e.username) > 0 {
		dumpArgs = append(dumpArgs, "-u", e.username)
	}
	if len(e.password) > 0 {
		dumpArgs = append(dumpArgs, `-p`+e.password)
	}

	dumpArgs = append(dumpArgs, e.database)

	if !helper.IsExistsPath(e.Base.StorePath) {
		helper.MkdirP(e.Base.StorePath)
	}

	dumpArgs = append(dumpArgs, "--result-file="+e.Base.DumpFilePath)
	return dumpArgs
}

func (e *MySQL) dump() (result *DumpResult, err error) {
	result = &DumpResult{}

	logger.Info("-------------Dumping MySQL---------------")
	out, err := helper.Exec("mysqldump", e.dumpOptions()...)

	if err != nil {
		return result, fmt.Errorf("-> Dump error: %s", err)
	}
	logger.Info("dump path:", e.Base.DumpFilePath)
	logger.Info("Output:", out)
	result.Path = e.Base.DumpFilePath

	isSkipBucketStore := viper.GetBool("skip-bucket-store")

	if !isSkipBucketStore {
		logger.Info("-------------Archiving sql file---------------")
		result.TarFilePath = e.Base.DumpFilePath + ".tar.gz"
		out, err = helper.Exec("tar", "-czvf", result.TarFilePath, "-C", e.Base.StorePath, e.Base.FileName)

		if err != nil {
			return result, fmt.Errorf("-> Tar error: %s", err)
		}

		if err == nil && e.Base.Storer != nil {
			err = e.Base.Storer.Store(result.TarFilePath)

			if err != nil {
				return result, fmt.Errorf("-> Store error: %s", err)
			}
		}

		deleteDumpFiles := viper.GetBool("delete-dump-files")

		if deleteDumpFiles {
			if helper.IsExistsPath(result.TarFilePath) {
				logger.Info("----------------Removing Local backup tar file---------------")
				os.Remove(result.TarFilePath)
			}
		}
	}

	return result, err
}
