package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-dumper/helper"
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	// Imports the Google Cloud Storage client package.
	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type GCS struct {
	Base
	bucket string
	path   string
	jkey   jkey
	client *storage.Client
	ctx    context.Context
}

func (e *GCS) createGCSClient() (err error) {
	e.ctx = context.Background()
	e.bucket = viper.GetString("src.bucket")
	e.jkey = jkey{}
	e.jkey.loadDef()
	e.jkey.loadConf()

	jbts, err := json.Marshal(e.jkey)

	logger.Info("project:", e.jkey.ProjectID)

	if err != nil {
		return errors.Wrap(err, "could not convert api json to bytes")
	}

	conf, err := google.JWTConfigFromJSON(jbts, "https://www.googleapis.com/auth/devstorage.full_control")
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "could not create google jwt config")
	}

	ts := conf.TokenSource(e.ctx)

	client, err := storage.NewClient(e.ctx, option.WithTokenSource(ts))

	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "Could not create storage Client")
	}

	e.client = client
	return nil
}

func (e *GCS) Store(fileKey string) (err error) {
	logger.Info("-----------Storing to Bucket---------------")

	err = e.createGCSClient()

	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "Could not create storage Client")
	}

	f, err := os.Open(fileKey)

	if err != nil {
		logger.Error(err)
		return err
	}
	defer f.Close()

	wc := e.client.Bucket(e.bucket).Object(fileKey).NewWriter(e.ctx)
	if _, err = io.Copy(wc, f); err != nil {
		logger.Error(err)
		return err
	}
	if err := wc.Close(); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("-----------Backup file has been stored successfully---------------")

	return
}

func (e *GCS) Download(fileKey string) (dumpFilePath string, err error) {
	logger.Info("-----------Downloading from Bucket---------------")

	if !helper.IsExistsPath(e.Base.StorePath) {
		helper.MkdirP(e.Base.StorePath)
	}
	err = e.createGCSClient()

	if err != nil {
		logger.Error(err)
		return "", errors.Wrap(err, "Could not create storage Client")
	}
	_, file := filepath.Split(fileKey)

	localFilePathTar := filepath.Join(e.Base.StorePath, file)

	f, err := os.Create(localFilePathTar)

	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer f.Close()
	logger.Info(fileKey)
	logger.Info(e.bucket)
	rc, err := e.client.Bucket(e.bucket).Object(fileKey).NewReader(e.ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("readFile: unable to open file from bucket %v, file %v: %v", e.bucket, fileKey, err))
		return "", err
	}

	if _, err = io.Copy(f, rc); err != nil {
		logger.Error(err)
		return "", err
	}
	if err := rc.Close(); err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Info("-----------File downloaded Successfully---------------")

	dumpFilePath = strings.Replace(localFilePathTar, ".tar.gz", "", -1)

	//Extracting the  tar file
	logger.Info("------------------Extracting the tar file--------------------")
	_, err = helper.Exec("tar", "-xzvf", localFilePathTar, "-C", e.Base.StorePath)

	if err != nil {
		return "", fmt.Errorf("-> Tar extract error: %s", err)
	}

	logger.Info("------------------File extracted successfully--------------------")

	return
}
