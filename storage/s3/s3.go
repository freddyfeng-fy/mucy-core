package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/freddyfeng-fy/mucy-core/core"
	"github.com/freddyfeng-fy/mucy-core/storage/storage"
	"github.com/freddyfeng-fy/mucy-core/utils/strs"
	"github.com/freddyfeng-fy/mucy-core/utils/urls"
	"sync"
	"time"
)

type Config struct {
	UploadHost   string `mapstructure:"upload_host" json:"upload_host" yaml:"upload_host"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	AccessId     string `mapstructure:"access_id" json:"access_id" yaml:"access_id"`
	AccessSecret string `mapstructure:"access_secret" json:"access_secret" yaml:"access_secret"`
	Bucket       string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
}

type s3 struct {
	config    *Config
	s3session *session.Session
}

var (
	s    *s3
	err  error
	once *sync.Once
)

func Init(config Config) (storage.Storage, error) {
	once = &sync.Once{}
	once.Do(func() {
		s = &s3{}
		s.config = &config
		s.s3session, err = session.NewSession(&aws.Config{
			Endpoint:         &s.config.UploadHost,
			Credentials:      credentials.NewStaticCredentials(s.config.AccessId, s.config.AccessSecret, ""),
			Region:           aws.String("apac"),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(false),
		})
		if err != nil {
			core.App.Log.Error(err.Error())
		}
		storage.Register(storage.S3, s)
	})
	return s, nil
}

func (s3 *s3) PutImage(data []byte, contentType string) (*storage.InputRes, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := storage.GenerateImageKey(data, contentType)
	if err := s3.PutObject(key, data, contentType); err != nil {
		return nil, err
	}
	return &storage.InputRes{
		Url: urls.UrlJoin(s3.config.Host, key),
		Key: key,
	}, nil
}

func (s3 *s3) PutVideo(data []byte, contentType string) (*storage.InputRes, error) {
	if strs.IsBlank(contentType) {
		contentType = "video/mp4"
	}
	key := storage.GenerateVideoKey(data, contentType)
	if err := s3.PutObject(key, data, contentType); err != nil {
		return nil, err
	}
	return &storage.InputRes{
		Url: urls.UrlJoin(s3.config.Host, key),
		Key: key,
	}, nil
}

func (s3 *s3) PutObject(key string, data []byte, contentType string) error {
	uploader := s3manager.NewUploader(s3.s3session)
	_, errUpS3 := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s3.config.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if errUpS3 != nil {
		return errUpS3
	}
	return nil
}

func (s3 *s3) CopyImage(originUrl string) (*storage.InputRes, error) {
	data, contentType, err := storage.Download(originUrl)
	if err != nil {
		return nil, err
	}
	return s3.PutImage(data, contentType)
}

func (s3 *s3) GetUrl(key string) (string, error) {
	svc := awss3.New(s3.s3session)
	req, _ := svc.GetObjectRequest(&awss3.GetObjectInput{
		Bucket: aws.String(s3.config.Bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(3 * time.Hour)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}

func (s3 *s3) DeleteObject(key string) error {
	svc := awss3.New(s3.s3session)
	_, err := svc.DeleteObject(&awss3.DeleteObjectInput{
		Bucket: aws.String(s3.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	err = svc.WaitUntilObjectNotExists(&awss3.HeadObjectInput{
		Bucket: aws.String(s3.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}
