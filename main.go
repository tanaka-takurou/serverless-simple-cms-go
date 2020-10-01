package main

import (
	"io"
	"os"
	"log"
	"math"
	"sort"
	"bytes"
	"regexp"
	"context"
	"strconv"
	"net/http"
	"html/template"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type PageData struct {
	Title          string        `json:"title"`
	Page           int           `json:"page"`
	PageList       []int         `json:"pageList"`
	Content        ContentData   `json:"content"`
	SitemapList    []SitemapData `json:"sitemapList"`
	SitemapHeadTag string        `json:"sitemapHeadTag"`
	URL            string        `json:"url"`
}

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

type Response events.APIGatewayProxyResponse

var cfg aws.Config
var dynamodbClient *dynamodb.Client

const (
	DataTypeConst        int    = 0
	DataTypeItem         int    = 1
	DataTypeCategory     int    = 2
	DataTypeItemCategory int    = 3
	DataTypeSitemap      int    = 4
	Layout               string = "2006-01-02 15:04"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var err error
	baseTitle := "Example Site "
	tmp := template.New("tmp")
	contentType := "text/html"
	var dat PageData
	dat.URL = os.Getenv("URL")
	p := request.PathParameters
	q := request.QueryStringParameters
	r := request.Resource
	buf := new(bytes.Buffer)
	fw := io.Writer(buf)
	if r == "/sitemap" || p["proxy"] == "sitemap" {
		dat.SitemapHeadTag = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
		contentType = "application/xml"
		dat.SitemapList = scanSitemap(ctx)
		if err != nil {
			log.Println(err)
		}
		tmp = getSitemapTemplates()
	} else {
		dat.Content = scanContentData(ctx)
		if err != nil {
			log.Println(err)
		}
		category, page := extractCategoryAndPage(p["proxy"])
		if len(page) == 0 {
			page = q["page"]
		}
		if len(category) == 0 {
			category = q["category"]
		}
		maxContentPerPage := 10
		maxPage := int(math.Ceil(float64(len(dat.Content.ItemList))/float64(maxContentPerPage)))
		tmp = getDefaultTemplates()
		if false {
		// if contains(dat.CategoryNameList, category) {
		/*
			categoryContentMap_ := scanCategoryContents(ctx)
			dat.Title = baseTitle + category
			dat.CategoryPath = "/category/" + category
			dat.CategoryName = getCategoryDisplayName(category, categoryList_)
			dat.Page = 1
			pageNumber := 1
			dat.PageList = []int{}
			categoryItemList := getCategoryContent(categoryContentMap_[category], contentList_)
			dat.ItemList = getContentRange(pageNumber, maxContentPerPage, len(categoryContentMap_[category]), categoryItemList)
		*/
		} else {
			pageNumber := 1
			if len(page) > 0 {
				pageNumber, _ = strconv.Atoi(page)
			}
			if pageNumber > 1 && pageNumber <= maxPage {
				dat.Title = baseTitle + page
				dat.Page = pageNumber
				dat.PageList = getPageList(pageNumber, maxPage)
				dat.Content.ItemList = getContentRange(pageNumber, maxContentPerPage, len(dat.Content.ItemList), dat.Content.ItemList)
			} else {
				dat.Title = baseTitle
				dat.Page = 1
				dat.PageList = getPageList(1, maxPage)
				dat.Content.ItemList = getContentRange(1, maxContentPerPage, len(dat.Content.ItemList), dat.Content.ItemList)
			}
		}
	}
	if e := tmp.ExecuteTemplate(fw, "base", dat); e != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
		}, e
	}
	return Response{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            string(buf.Bytes()),
		Headers: map[string]string{
			"Content-Type": contentType,
		},
	}, nil
}

func contains(s []string, t string) bool {
	for _, v := range s {
		if t == v {
			return true
		}
	}
	return false
}

func getPageList(current int, max int) []int {
	var res []int
	pagerWidth := 2
	i := 1
	for {
		if i > max {
			break
		}
		if i != 1 && i != max {
			if i < current - pagerWidth {
				res = append(res, 0)
				i = current - pagerWidth
				continue
			}
			if i > current + pagerWidth {
				res = append(res, 0)
				i = max
				continue
			}
		}
		res = append(res, i)
		i++
	}
	return res
}

func extractCategoryAndPage(proxyPathParameter string) (string, string) {
	if len(proxyPathParameter) > 0 {
		res := regexp.MustCompile("[/]").Split(proxyPathParameter, -1)
		if len(res) > 0 {
			if res[0] == "category" {
				if len(res) > 2 && res[2] == "page" {
					return res[1], res[3]
				}
				return res[1], ""
			} else if res[0] == "page" {
				return "", res[1]
			}
		}
	}
	return "", ""
}

func getContentRange(page int, perPage int, maxContent int, data []ItemData) []ItemData {
	var list []ItemData
	mn := (page - 1) * perPage
	mx := int(math.Min(float64(page * perPage), float64(maxContent)))
	for i := mn; i < mx; i++ {
		list = append(list, data[maxContent - i - 1])
	}
	return list
}

func scan(ctx context.Context, filt expression.ConditionBuilder)(*dynamodb.ScanOutput, error)  {
	if dynamodbClient == nil {
		dynamodbClient = dynamodb.New(cfg)
	}
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, err
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(os.Getenv("TABLE_NAME")),
	}
	req := dynamodbClient.ScanRequest(input)
	res, err := req.Send(ctx)
	return res.ScanOutput, err
}

func GetDynamoDataCount(ctx context.Context, dataType int)(int, error)  {
	result, err := scan(ctx, expression.Name("item_type").Equal(expression.Value(dataType)))
	if err != nil {
		return 0, err
	}
	return len(result.Items), nil
}

func GetDynamoDataList(ctx context.Context, dataType int)([]DynamoData, error) {
	var dynamoDataList []DynamoData
	cond1 := expression.Name("item_type").Equal(expression.Value(dataType))
	cond2 := expression.Name("status").Equal(expression.Value(0))
	result, err := scan(ctx, cond1.And(cond2))
	if err != nil {
		log.Println(err)
		return dynamoDataList, err
	}
	for _, i := range result.Items {
		item := DynamoData{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Println(err)
		} else {
			dynamoDataList = append(dynamoDataList, item)
		}
	}
	sort.Slice(dynamoDataList, func(i, j int) bool { return dynamoDataList[i].Id < dynamoDataList[j].Id })
	return dynamoDataList, nil
}

func GetItemList(ctx context.Context)([]DynamoData, error) {
	return GetDynamoDataList(ctx, DataTypeItem)
}

func GetCategoryList(ctx context.Context)([]DynamoData, error) {
	return GetDynamoDataList(ctx, DataTypeCategory)
}

func GetItemCategoryMap(ctx context.Context)([]DynamoData, error) {
	return GetDynamoDataList(ctx, DataTypeItemCategory)
}

func GetSitemapDataList(ctx context.Context)([]DynamoData, error) {
	return GetDynamoDataList(ctx, DataTypeSitemap)
}

func scanContentData(ctx context.Context) ContentData {
	result, err := GetItemList(ctx)
	if err != nil {
		log.Println(err)
		return ContentData{}
	}
	var itemDataList []ItemData
	for _, i := range result {
		item := ItemData{}
		json.Unmarshal([]byte(i.Data), &item)
		itemDataList = append(itemDataList, item)
	}
	result, err = GetCategoryList(ctx)
	if err != nil {
		log.Println(err)
		return ContentData{}
	}
	var categoryNameList []string
	for _, i := range result {
		if len(i.Data) < 1 {
			log.Println("Error Category Name is nil")
		} else {
			categoryNameList = append(categoryNameList, i.Data)
		}
	}
	return ContentData{
		Const: ConstData{},
		ItemList: itemDataList,
		CategoryNameList: categoryNameList,
		CategoryItemMap: map[string][]int{},
	}
}

func scanSitemap(ctx context.Context) []SitemapData {
	result, err := GetSitemapDataList(ctx)
	if err != nil {
		log.Println(err)
		return nil
	}
	var sitemapDataList []SitemapData
	for _, i := range result {
		item := SitemapData{}
		json.Unmarshal([]byte(i.Data), &item)
		sitemapDataList = append(sitemapDataList, item)
	}
	return sitemapDataList
}

func getDefaultTemplates() *template.Template {
	funcMap := template.FuncMap {
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
		"safeurl": func(text string) template.URL { return template.URL(text) },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
	}
	return template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"templates/index.html",
		"templates/view.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/pager.html",
		"templates/button.html",
		"templates/item_list.html",
	))
}

func getSitemapTemplates() *template.Template {
	funcMap := template.FuncMap {
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
	}
	return template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/sitemap.xml"))
}

func init() {
	var err error
	cfg, err = external.LoadDefaultAWSConfig()
	cfg.Region = os.Getenv("REGION")
	if err != nil {
		log.Print(err)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
