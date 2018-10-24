package dump

import (
	"fmt"

	"github.com/go-dumper/store"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Base struct {
	DbType       string
	StorePath    string
	StoreType    string
	DumpFilePath string
	FileName     string
	Storer       store.Storer
}

//Init base
func newBase() (base Base) {
	base = Base{
		DbType:    viper.GetString("src.database_type"),
		StorePath: viper.GetString("src.store_path"),
		StoreType: viper.GetString("src.store_type"),
	}

	isSkipBucketStore := viper.GetBool("skip-bucket-store")

	if !isSkipBucketStore {
		result, err := store.New()

		if err == nil {
			base.Storer = result
		}
	}

	return
}

func Run() (result *DumpResult, err error) {
	base := newBase()
	var dumper Dumper

	switch base.DbType {
	case "mysql":
		dumper = &MySQL{Base: base}
	default:
		return nil, fmt.Errorf("config `type: %s` is not implemented", base.DbType)
	}

	logger.Info("-------- database type ----------", base.DbType)

	// Start dummping
	result, err = dumper.Dump()
	if err != nil {
		logger.Info(err)
		return nil, err
	}

	return result, err
}
