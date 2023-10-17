package local

import (
	"github.com/freddyfeng-fy/mucy-core/storage/storage"
	"github.com/freddyfeng-fy/mucy-core/utils/strs"
	"github.com/freddyfeng-fy/mucy-core/utils/urls"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	RootDir string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}

type local struct {
	config *Config
}

var (
	l    *local
	once *sync.Once
)

func Init(config Config) (storage.Storage, error) {
	once = &sync.Once{}
	once.Do(func() {
		l = &local{
			config: &config,
		}
		storage.Register(storage.Local, l)
	})
	return l, nil
}

func (l *local) PutImage(data []byte, contentType string) (*storage.InputRes, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := storage.GenerateImageKey(data, contentType)
	if err := l.PutObject(key, data, contentType); err != nil {
		return nil, err
	}
	return &storage.InputRes{
		Url: urls.UrlJoin(l.config.AppUrl, key),
		Key: key,
	}, nil
}

func (l *local) PutVideo(data []byte, contentType string) (*storage.InputRes, error) {
	if strs.IsBlank(contentType) {
		contentType = "video/mp4"
	}
	key := storage.GenerateVideoKey(data, contentType)
	if err := l.PutObject(key, data, contentType); err != nil {
		return nil, err
	}
	return &storage.InputRes{
		Url: urls.UrlJoin(l.config.AppUrl, key),
		Key: key,
	}, nil
}

func (l *local) PutObject(key string, data []byte, contentType string) error {
	if err := os.MkdirAll("/", os.ModeDir); err != nil {
		return err
	}
	filename := filepath.Join(l.config.RootDir, key)
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (l *local) CopyImage(originUrl string) (*storage.InputRes, error) {
	data, contentType, err := storage.Download(originUrl)
	if err != nil {
		return nil, err
	}
	return l.PutImage(data, contentType)
}

func (l *local) GetUrl(key string) (string, error) {
	return "", nil
}

func (l *local) DeleteObject(key string) error {
	return nil
}
