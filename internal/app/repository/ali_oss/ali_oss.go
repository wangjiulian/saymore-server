package ali_oss

import (
	"com.say.more.server/config"
	"com.say.more.server/internal/app/errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

type AliOSSBucket struct {
	bucket            *oss.Bucket
	fileExpired       int64
	exportFileExpired int64
	baseUrl           string
}

func NewAliOSSBucket(cfg *config.AliOSS) (*AliOSSBucket, error) {
	Endpoint := cfg.Endpoint
	AccessKeyId := cfg.AccessKeyId
	AccessKeySecret := cfg.AccessKeySecret

	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(cfg.BucketName)
	if err != nil {
		return nil, err
	}
	return &AliOSSBucket{
		baseUrl:     cfg.BucketUrl,
		bucket:      bucket,
		fileExpired: 180 * 86400,
	}, nil
}

func (a *AliOSSBucket) GetUrl(key string) (string, error) {
	signedGetURL, err := a.bucket.SignURL(key, oss.HTTPGet, a.fileExpired)
	if err != nil {
		return "", errors.Wrap(err)
	}
	return signedGetURL, nil
}

func (a *AliOSSBucket) Get(key string) string {
	return a.baseUrl + "/" + key
}

//func (a *AliOSSBucket) SignedPut(key string) (string, error) {
//	// // Signed direct upload with optional parameters.
//	options := []oss.Option{
//	}
//	signedPutURL, err := a.Bucket.SignURL(key, oss.HTTPPut, a.fileExpired, options...)
//	if err != nil {
//		return "", err
//	}
//	return signedPutURL, nil
//}

func (a *AliOSSBucket) SignedGet(key string) (string, error) {
	signedGetURL, err := a.bucket.SignURL(key, oss.HTTPGet, a.fileExpired)
	if err != nil {
		return "", err
	}
	return signedGetURL, nil
}

func (a *AliOSSBucket) PutObject(key string, src io.Reader, options []oss.Option) error {
	if len(options) > 0 {
		return a.bucket.PutObject(key, src, options...)
	}
	return a.bucket.PutObject(key, src)
}

func (a *AliOSSBucket) PutObjectFromFile(key, filePath string) error {
	return a.bucket.PutObjectFromFile(key, filePath)
}

func (a *AliOSSBucket) GetExpiredDay(operationType int) int64 {
	return a.fileExpired
}
