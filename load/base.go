package load

import (
	"fmt"
	"strings"

	"net/url"

	"github.com/go-dumper/store"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Base struct {
	DbType          string
	DumpFilePath    string
	FileName        string
	CreateNewSchema bool
	Storer          store.Storer
}

//Init base
func newBase(dumpFilePath string) (base Base) {
	base = Base{
		DbType:          viper.GetString("src.database_type"),
		DumpFilePath:    dumpFilePath,
		CreateNewSchema: viper.GetBool("create-database"),
	}

	parsedUrl, _ := url.Parse(dumpFilePath)

	if parsedUrl.Scheme != "" {
		viper.Set("src.store_type", parsedUrl.Scheme)
		viper.Set("src.bucket", parsedUrl.Host)
		base.DumpFilePath = strings.TrimLeft(parsedUrl.Path, "/")
		result, err := store.New()

		if err == nil {
			base.Storer = result
		}
	}

	return
}

func Run(dumpFilePath string) (err error) {
	base := newBase(dumpFilePath)
	var loader Loader

	switch base.DbType {
	case "mysql":
		loader = &MySQL{Base: base}
	default:
		return fmt.Errorf("config `type: %s` is not implemented", base.DbType)
	}

	logger.Info("=> database type |", base.DbType)

	err = loader.Load()
	if err != nil {
		logger.Info(err)
		return err
	}

	//Update user's email
	logger.Info("-----------Updating User or other table data---------------")

	/*err = loader.Update()
	if err != nil {
		logger.Info(err)
		return err
	}*/

	return nil
}
