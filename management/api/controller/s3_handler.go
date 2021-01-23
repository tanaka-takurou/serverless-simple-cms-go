package controller

import (
	"os"
	"log"
	"time"
	"bytes"
	"errors"
	"strings"
	"context"
	"path/filepath"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type S3ContentData struct {
	Key          string `json:"key"`
	Size         int    `json:"size"`
	LastModified string `json:"lastmodified"`
}

var s3Client *s3.Client

const (
	Layout2        string = "20060102150405"
	Layout3        string = "2006/01/02 15:04:05"
	ImgFilePath    string = "img"
	StaticFilePath string = "static"
)

func UploadImage(ctx context.Context, filename string, filedata string)(string, error) {
	t := time.Now()
	b64data := filedata[strings.IndexByte(filedata, ',')+1:]
	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Print(err)
		return "", err
	}
	extension := filepath.Ext(filename)
	var contentType string

	switch extension {
	case ".jpg":
		contentType = "image/jpeg"
	case ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".png":
		contentType = "image/png"
	default:
		return "", errors.New("this extension is invalid")
	}
	filename_ := string([]rune(filename)[:(len(filename) - len(extension))]) + t.Format(Layout2) + extension
	uploader := manager.NewUploader(s3.NewFromConfig(cfg))
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		ACL: types.ObjectCannedACLPublicRead,
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key: aws.String(ImgFilePath + "/" + filename_),
		Body: bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Print(err)
		return "", err
	}
	imgUrl := "https://" + os.Getenv("BUCKET_NAME") + ".s3-" + os.Getenv("REGION") + ".amazonaws.com/" + ImgFilePath + "/" + filename_
	return imgUrl, nil
}

func UploadFile(ctx context.Context, filedata string, contentType string) error {
	t := time.Now()
	var filename string
	switch contentType {
	case "text/css":
		filename = t.Format(Layout2) + ".css"
		SetCssFileName(ctx,"https://" + os.Getenv("BUCKET_NAME") + ".s3-" + os.Getenv("REGION") + ".amazonaws.com/" +  StaticFilePath + "/" + filename)
	case "text/javascript":
		filename = t.Format(Layout2) + ".js"
		SetJsFileName(ctx, "https://" + os.Getenv("BUCKET_NAME") + ".s3-" + os.Getenv("REGION") + ".amazonaws.com/" + StaticFilePath + "/" + filename)
	default:
		filename = t.Format(Layout2) + ".txt"
	}
	uploader := manager.NewUploader(s3.NewFromConfig(cfg))
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		ACL: types.ObjectCannedACLPublicRead,
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key: aws.String(StaticFilePath + "/" + filename),
		Body: bytes.NewReader([]byte(filedata)),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func GetS3Data(ctx context.Context)(string, interface{}, error) {
	if s3Client == nil {
		s3Client = s3.NewFromConfig(cfg)
	}
	input := &s3.ListObjectsInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
	}
	res, err := s3Client.ListObjects(ctx, input)
	if err != nil {
		return "", nil, err
	}
	var s3Contents []S3ContentData
	for _, v := range res.Contents {
		s3Contents = append(s3Contents, S3ContentData{
			Key: aws.ToString(v.Key),
			Size: int(v.Size),
			LastModified: aws.ToTime(v.LastModified).Format(Layout3),
		})
	}
	return os.Getenv("BUCKET_NAME"), s3Contents, nil
}
