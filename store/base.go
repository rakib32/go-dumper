package store

import (
	"fmt"

	"github.com/spf13/viper"
)

type Base struct {
	StoreType string
	StorePath string
}

func newBase() (base Base) {
	base = Base{
		StoreType: viper.GetString("src.store_type"),
		StorePath: viper.GetString("src.store_path"),
	}

	return
}

// Run storage
func New() (storer Storer, err error) {
	base := newBase()

	switch base.StoreType {
	case "gs":
		storer = &GCS{Base: base}
	default:
		return nil, fmt.Errorf("[%s] storage type has not implemented", base.StoreType)
	}

	return storer, err
}
