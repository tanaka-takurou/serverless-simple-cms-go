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
	Category       string        `json:"category"`
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
	Const              ConstData        `json:"constData"`
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
	baseTitle := "Example Site "
	tmp := template.New("tmp")
	contentType := "text/html"
	var dat PageData
	dat.URL = "https://" + request.RequestContext.DomainName + "/"
	p := request.PathParameters
	q := request.QueryStringParameters
	r := request.Resource
	buf := new(bytes.Buffer)
	fw := io.Writer(buf)
	if r == "/sitemap" || p["proxy"] == "sitemap" {
		dat.SitemapHeadTag = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
		contentType = "application/xml"
		dat.SitemapList = scanSitemap(ctx)
		tmp = getSitemapTemplates()
	} else {
		dat.Content = scanContentData(ctx)
		if len(dat.Content.Const.Title) > 0 {
			baseTitle = dat.Content.Const.Title
		}
		category, page := extractCategoryAndPage(p["proxy"])
		if len(page) == 0 {
			page = q["page"]
		}
		if len(category) == 0 {
			category = q["category"]
		}
		maxContentPerPage := 10
		tmp = getDefaultTemplates()
		pageNumber := 1
		if contains(dat.Content.CategoryNameList, category) {
			baseTitle = baseTitle + " " + category
			maxPage := int(math.Ceil(float64(len(dat.Content.CategoryItemMap[category]))/float64(maxContentPerPage)))
			if len(page) > 0 {
				pageNumber, _ = strconv.Atoi(page)
			}
			if pageNumber > 1 && pageNumber <= maxPage {
				dat = getPageDataTargetCategory(dat, baseTitle + " " + page, pageNumber, maxPage, maxContentPerPage, category)
			} else {
				dat = getPageDataTargetCategory(dat, baseTitle, 1, maxPage, maxContentPerPage, category)
			}

		} else {
			maxPage := int(math.Ceil(float64(len(dat.Content.ItemList))/float64(maxContentPerPage)))
			if len(page) > 0 {
				pageNumber, _ = strconv.Atoi(page)
			}
			if pageNumber > 1 && pageNumber <= maxPage {
				dat = getPageData(dat, baseTitle + " " + page, pageNumber, maxPage, maxContentPerPage)
			} else {
				dat = getPageData(dat, baseTitle, 1, maxPage, maxContentPerPage)
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

func getContentRangeTargetCategory(page int, perPage int, targetItemIdList []int, data []ItemData) []ItemData {
	var list []ItemData
	maxContent := len(targetItemIdList)
	mn := (page - 1) * perPage
	mx := int(math.Min(float64(page * perPage), float64(maxContent)))
	for i := mn; i < mx; i++ {
		list = append(list, data[targetItemIdList[maxContent - i - 1] - 1])
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

func GetConst(ctx context.Context)(DynamoData, error) {
	var dynamoData DynamoData
	cond1 := expression.Name("item_type").Equal(expression.Value(DataTypeConst))
	cond2 := expression.Name("id").Equal(expression.Value(1))
	result, err := scan(ctx, cond1.And(cond2))
	if err != nil {
		log.Println(err)
		return dynamoData, err
	}
	if len(result.Items) > 0 {
		item := DynamoData{}
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
		if err != nil {
			return dynamoData, err
		} else {
			dynamoData = item
		}
	}
	return dynamoData, nil
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
	result, err := scan(ctx, expression.NotEqual(expression.Name("status"), expression.Value(-1)))
	if err != nil {
		log.Println(err)
		return ContentData{}
	}
	constData := ConstData{}
	itemDataList := []ItemData{}
	categoryNameList := []string{}
	itemCategoryMap := make(map[string][]int)
	for _, v := range result.Items {
		dynamoData := DynamoData{}
		err = dynamodbattribute.UnmarshalMap(v, &dynamoData)
		if err != nil {
			log.Println(err)
		} else {
			switch {
			case dynamoData.Type == DataTypeConst:
				json.Unmarshal([]byte(dynamoData.Data), &constData)
			case dynamoData.Type == DataTypeItem:
				item := ItemData{}
				json.Unmarshal([]byte(dynamoData.Data), &item)
				itemDataList = append(itemDataList, item)
			case dynamoData.Type == DataTypeCategory:
				categoryNameList = append(categoryNameList, dynamoData.Data)
			case dynamoData.Type == DataTypeItemCategory:
				kvs := KVSData{}
				itemIdList := []int{}
				json.Unmarshal([]byte(dynamoData.Data), &kvs)
				json.Unmarshal([]byte(kvs.V), &itemIdList)
				itemCategoryMap[kvs.K] = itemIdList
			default:
			}
		}
	}
	return ContentData{
		Const: constData,
		ItemList: itemDataList,
		CategoryNameList: categoryNameList,
		CategoryItemMap: itemCategoryMap,
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

func getPageData(dat PageData, title string, pageNumber int, maxPage int, maxContentPerPage int) PageData {
	dat.Title = title
	dat.Page = pageNumber
	dat.PageList = getPageList(pageNumber, maxPage)
	dat.Content.ItemList = getContentRange(pageNumber, maxContentPerPage, len(dat.Content.ItemList), dat.Content.ItemList)
	return dat
}

func getPageDataTargetCategory(dat PageData, title string, pageNumber int, maxPage int, maxContentPerPage int, category string) PageData {
	dat.Title = title
	dat.Page = pageNumber
	dat.PageList = getPageList(pageNumber, maxPage)
	tmpItemList := getContentRangeTargetCategory(pageNumber, maxContentPerPage, dat.Content.CategoryItemMap[category], dat.Content.ItemList)
	dat.Content.ItemList = tmpItemList
	dat.Category = category
	return dat
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
