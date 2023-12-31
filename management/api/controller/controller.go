package controller

import (
	"os"
	"log"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ItemData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CategoryIds []int  `json:"categoryids"`
}

type SitemapData struct {
	Loc        string `json:"loc"`
	Lastmod    string `json:"lastmod"`
	Changefreq string `json:"changefreq"`
	Priority   string `json:"priority"`
}

type ContentData struct {
	Const              ConstData        `json:"constData"`
	ItemList           []ItemData       `json:"itemList"`
	CategoryNameList   []string         `json:"categoryNameList"`
	CategoryItemMap    map[string][]int `json:"categoryItemMap"`
	SitemapDataList    []SitemapData    `json:"sitemap"`
}

type ConstData struct {
	Title       string `json:"title"`
	HeadImage   string `json:"headImage"`
	JsFileName  string `json:"jsFileName"`
	CssFileName string `json:"cssFileName"`
}

type DynamoData struct {
	Id      int    `dynamodbav:"id"`
	Data    string `dynamodbav:"data"`
	Type    int    `dynamodbav:"item_type"`
	Status  int    `dynamodbav:"status"`
	Created string `dynamodbav:"created"`
}

type KVSData struct {
	K string `json:"k"`
	V string `json:"v"`
}

var cfg aws.Config

const (
	DataTypeConst        int    = 0
	DataTypeItem         int    = 1
	DataTypeCategory     int    = 2
	DataTypeItemCategory int    = 3
	DataTypeSitemap      int    = 4
	Layout               string = "2006-01-02 15:04:05"
)

func FixJS(ctx context.Context, filedata string) error {
	return UploadFile(ctx, filedata, "text/javascript")
}

func FixCSS(ctx context.Context, filedata string) error {
	return UploadFile(ctx, filedata, "text/css")
}

func InitConfig(ctx context.Context) {
	var err error
	cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Print(err)
	}
}
