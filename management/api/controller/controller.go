package controller

import (
	"os"
	"log"
	"context"
	"io/ioutil"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
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
	Const              ConstData        `json:"const"`
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
	Id      int    `json:"id"`
	Data    string `json:"data"`
	Type    int    `json:"item_type"`
	Status  int    `json:"status"`
	Created string `json:"created"`
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
	FilePath             string = "sample_data/data.json"
)

func JsonDataLoad()(ContentData, error) {
	content := new(ContentData)
	jsonString, err := ioutil.ReadFile(FilePath)
	if err != nil {
		return *content, err
	}
	json.Unmarshal(jsonString, content)
	return *content, nil
}

func FixJS(ctx context.Context, filedata string) error {
	return UploadFile(ctx, filedata, "text/javascript")
}

func FixCSS(ctx context.Context, filedata string) error {
	return UploadFile(ctx, filedata, "text/css")
}

func init() {
	var err error
	cfg, err = external.LoadDefaultAWSConfig()
	cfg.Region = os.Getenv("REGION")
	if err != nil {
		log.Print(err)
	}
}
