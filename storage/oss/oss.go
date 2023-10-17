package oss

import (
	"bytes"
	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/freddyfeng-fy/mucy-core/storage/storage"
	"github.com/freddyfeng-fy/mucy-core/utils/strs"
	"github.com/freddyfeng-fy/mucy-core/utils/urls"
	"sync"
)

type Config struct {
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	AccessKeyId     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	IsSsl           bool   `mapstructure:"is_ssl" json:"is_ssl" yaml:"is_ssl"`
	IsPrivate       bool   `mapstructure:"is_private" json:"is_private" yaml:"is_private"`
}

type oss struct {
	config *Config
	client *alioss.Client
	bucket *alioss.Bucket
}

var (
	o       *oss
	once    *sync.Once
	initErr error
)

func Init(config Config) (storage.Storage, error) {
	once = &sync.Once{}
	once.Do(func() {
		o = &oss{}
		o.config = &config

		o.client, initErr = alioss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
		if initErr != nil {
			return
		}

		o.bucket, initErr = o.client.Bucket(config.Bucket)
		if initErr != nil {
			return
		}

		storage.Register(storage.Oss, o)
	})
	if initErr != nil {
		return nil, initErr
	}
	return o, nil
}

func (aliyun *oss) PutImage(data []byte, contentType string) (*storage.InputRes, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := storage.GenerateImageKey(data, contentType)
	if err := aliyun.PutObject(key, data, contentType); err != nil {
		return nil, err
	}
	return &storage.InputRes{
		Url: urls.UrlJoin(aliyun.config.Host, key),
		Key: key,
	}, nil
}

func (aliyun *oss) PutVideo(data []byte, contentType string) (*storage.InputRes, error) {
	if strs.IsBlank(contentType) {
		contentType = "video/mp4"
	}
	key := storage.GenerateVideoKey(data, contentType)
	if err := aliyun.PutObject(key, data, contentType); err != nil {
		return nil, err
	}
	return &storage.InputRes{
		Url: urls.UrlJoin(aliyun.config.Host, key),
		Key: key,
	}, nil
}

func (aliyun *oss) PutObject(key string, data []byte, contentType string) error {
	var options []alioss.Option
	if strs.IsNotBlank(contentType) {
		options = append(options, alioss.ContentType(contentType))
	}
	if err := aliyun.bucket.PutObject(key, bytes.NewReader(data), options...); err != nil {
		return err
	}
	return nil
}

func (aliyun *oss) CopyImage(originUrl string) (*storage.InputRes, error) {
	data, contentType, err := storage.Download(originUrl)
	if err != nil {
		return nil, err
	}
	return aliyun.PutImage(data, contentType)
}

func (aliyun *oss) GetUrl(key string) (string, error) {
	return "", nil
}

func (aliyun *oss) DeleteObject(key string) error {
	return nil
}
