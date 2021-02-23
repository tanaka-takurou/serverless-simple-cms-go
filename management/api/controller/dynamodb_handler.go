package controller

import (
	"os"
	"fmt"
	"log"
	"time"
	"sort"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

var dynamodbClient *dynamodb.Client
const LastmodLayout string = "2006-01-02"

func scan(ctx context.Context, filt expression.ConditionBuilder)(*dynamodb.ScanOutput, error)  {
	if dynamodbClient == nil {
		dynamodbClient = dynamodb.NewFromConfig(cfg)
	}
	proj := expression.NamesList(expression.Name("id"), expression.Name("data"), expression.Name("item_type"), expression.Name("status"), expression.Name("created"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(os.Getenv("ITEM_TABLE_NAME")),
	}
	res, err := dynamodbClient.Scan(ctx, input)
	return res, err
}

func put(ctx context.Context, av map[string]types.AttributeValue) error {
	if dynamodbClient == nil {
		dynamodbClient = dynamodb.NewFromConfig(cfg)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("ITEM_TABLE_NAME")),
	}
	_, err := dynamodbClient.PutItem(ctx, input)
	return err
}

func update(ctx context.Context, an map[string]string, av map[string]types.AttributeValue, key map[string]types.AttributeValue, updateExpression string) error {
	if dynamodbClient == nil {
		dynamodbClient = dynamodb.NewFromConfig(cfg)
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: an,
		ExpressionAttributeValues: av,
		TableName: aws.String(os.Getenv("ITEM_TABLE_NAME")),
		Key: key,
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String(updateExpression),
	}

	_, err := dynamodbClient.UpdateItem(ctx, input)
	return err
}

func GetDynamoDataCount(ctx context.Context, dataType int)(int, error)  {
	result, err := scan(ctx, expression.Name("item_type").Equal(expression.Value(dataType)))
	if err != nil {
		return 0, err
	}
	return len(result.Items), nil
}

func SetConst(ctx context.Context, title string, headImage string, jsFileName string, cssFileName string) error {
	t := time.Now()
	data_, _ := json.Marshal(ConstData{
		Title: title,
		HeadImage: headImage,
		JsFileName: jsFileName,
		CssFileName: cssFileName,
	})
	item := DynamoData {
		Id: 1,
		Data: string(data_),
		Type: DataTypeConst,
		Status: 0,
		Created: t.Format(Layout),
	}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	err = put(ctx, av)
	if err != nil {
		return err
	}
	return nil
}

func SetSampleData(ctx context.Context, content ContentData) error {
	dataCount, err := GetDynamoDataCount(ctx, DataTypeItem)
	if err != nil {
		log.Println(err)
		return err
	} else if dataCount > 0 {
		return fmt.Errorf("Error: %s", "Data is already set.")
	}
	err = putItemListData(ctx, content.ItemList)
	if err != nil {
		log.Println(err)
		return err
	}
	err = putCategoryNameListData(ctx, content.CategoryNameList)
	if err != nil {
		log.Println(err)
		return err
	}
	err = putCategoryItemMapData(ctx, content.CategoryItemMap)
	if err != nil {
		log.Println(err)
		return err
	}
	err = putSitemapDataList(ctx, content.SitemapDataList)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func putItemData(ctx context.Context, itemId int, data ItemData) error {
	t := time.Now()
	data_, err := json.Marshal(data)
	if err != nil {
		return err
	}
	item := DynamoData {
		Id: itemId,
		Data: string(data_),
		Type: DataTypeItem,
		Status: 0,
		Created: t.Format(Layout),
	}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	err = put(ctx, av)
	if err != nil {
		return err
	}
	return nil
}

func putItemListData(ctx context.Context, data []ItemData) error {
	dataCount, err := GetDynamoDataCount(ctx, DataTypeItem)
	if err != nil {
		log.Println(err)
		return err
	}
	for i, v := range data {
		err := putItemData(ctx, dataCount + i + 1, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func putCategoryData(ctx context.Context, itemId int, data string) error {
	t := time.Now()
	item := DynamoData {
		Id: itemId,
		Data: data,
		Type: DataTypeCategory,
		Status: 0,
		Created: t.Format(Layout),
	}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	err = put(ctx, av)
	if err != nil {
		return err
	}
	return nil
}

func putCategoryNameListData(ctx context.Context, data []string) error {
	dataCount, err := GetDynamoDataCount(ctx, DataTypeCategory)
	if err != nil {
		log.Println(err)
		return err
	}
	for i, v := range data {
		err := putCategoryData(ctx, dataCount + i + 1, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func putCategoryItemData(ctx context.Context, itemId int, data KVSData) error {
	t := time.Now()
	data_, err := json.Marshal(data)
	if err != nil {
		return err
	}
	item := DynamoData {
		Id: itemId,
		Data: string(data_),
		Type: DataTypeItemCategory,
		Status: 0,
		Created: t.Format(Layout),
	}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	err = put(ctx, av)
	if err != nil {
		return err
	}
	return nil
}

func putCategoryItemMapData(ctx context.Context, data map[string][]int) error {
	dataCount, err := GetDynamoDataCount(ctx, DataTypeItemCategory)
	if err != nil {
		log.Println(err)
		return err
	}
	i := 0
	for k, v := range data {
		v_, _ := json.Marshal(v)
		err := putCategoryItemData(ctx, dataCount + i + 1, KVSData{K: k, V: string(v_)})
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func putSitemapData(ctx context.Context, itemId int, data SitemapData) error {
	t := time.Now()
	data_, _ := json.Marshal(data)
	item := DynamoData {
		Id: itemId,
		Data: string(data_),
		Type: DataTypeSitemap,
		Status: 0,
		Created: t.Format(Layout),
	}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	err = put(ctx, av)
	if err != nil {
		return err
	}
	return nil
}

func putSitemapDataList(ctx context.Context, data []SitemapData) error {
	dataCount, err := GetDynamoDataCount(ctx, DataTypeSitemap)
	if err != nil {
		log.Println(err)
		return err
	}
	for i, v := range data {
		err := putSitemapData(ctx, dataCount + i + 1, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDynamoDataList(ctx context.Context, dataType int)([]DynamoData, error) {
	dynamoDataList := []DynamoData{}
	cond1 := expression.Name("item_type").Equal(expression.Value(dataType))
	cond2 := expression.Name("status").Equal(expression.Value(0))
	result, err := scan(ctx, cond1.And(cond2))
	if err != nil {
		log.Println(err)
		return dynamoDataList, err
	}
	for _, i := range result.Items {
		item := DynamoData{}
		err = attributevalue.UnmarshalMap(i, &item)
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

func GetSingleDynamoData(ctx context.Context, dataType int, itemId int)(DynamoData, error) {
	var dynamoData DynamoData
	cond1 := expression.Name("item_type").Equal(expression.Value(dataType))
	cond2 := expression.Name("id").Equal(expression.Value(itemId))
	result, err := scan(ctx, cond1.And(cond2))
	if err != nil {
		log.Println(err)
		return dynamoData, err
	}
	if len(result.Items) > 0 {
		item := DynamoData{}
		err = attributevalue.UnmarshalMap(result.Items[0], &item)
		if err != nil {
			log.Println(err)
		} else {
			dynamoData = item
		}
	}
	return dynamoData, nil
}

func GetSingleConst(ctx context.Context)(DynamoData, error) {
	return GetSingleDynamoData(ctx, DataTypeConst, 1)
}

func GetSingleItem(ctx context.Context, itemId int)(DynamoData, error) {
	return GetSingleDynamoData(ctx, DataTypeItem, itemId)
}

func GetSingleCategory(ctx context.Context, itemId int)(DynamoData, error) {
	return GetSingleDynamoData(ctx, DataTypeCategory, itemId)
}

func GetSingleItemCategoryMap(ctx context.Context, itemId int)(DynamoData, error) {
	return GetSingleDynamoData(ctx, DataTypeItemCategory, itemId)
}

func GetSingleSitemapData(ctx context.Context, itemId int)(DynamoData, error) {
	return GetSingleDynamoData(ctx, DataTypeSitemap, itemId)
}

func UpdateDynamoData(ctx context.Context, dataType int, itemId int, data string) error {
	an := map[string]string{
		"#d": "data",
	}
	item := struct {NewData string `dynamodbav:":new_d"`}{data}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	item_ := struct {
		Id       int `dynamodbav:"id"`
		ItemType int `dynamodbav:"item_type"`
	}{
		Id: itemId,
		ItemType: dataType,
	}
	key, err := attributevalue.MarshalMap(item_)
	if err != nil {
		return err
	}
	updateExpression := "set #d = :new_d"
	return update(ctx, an, av, key, updateExpression)
}

func UpdateConst(ctx context.Context, data string) error {
	return UpdateDynamoData(ctx, DataTypeConst, 1, data)
}

func UpdateItem(ctx context.Context, itemId int, data string) error {
	return UpdateDynamoData(ctx, DataTypeItem, itemId, data)
}

func UpdateCategory(ctx context.Context, itemId int, data string) error {
	return UpdateDynamoData(ctx, DataTypeCategory, itemId, data)
}

func UpdateItemCategoryMap(ctx context.Context, itemId int, data string) error {
	return UpdateDynamoData(ctx, DataTypeItemCategory, itemId, data)
}

func UpdateSitemapData(ctx context.Context, itemId int, data string) error {
	return UpdateDynamoData(ctx, DataTypeSitemap, itemId, data)
}

func SetCategoryList(ctx context.Context, itemId int, categoryNames []string)([]int, error) {
	var err error
	var categoryIdList []int
	categoryList, err := GetCategoryList(ctx)
	if err != nil {
		log.Println(err)
		return categoryIdList, err
	}
	for _, v := range categoryNames {
		if len(v) < 1 {
			continue
		}
		categoryId := 0
		for _, w := range categoryList {
			if v == w.Data {
				categoryId = w.Id
				break
			}
		}
		if categoryId == 0 {
			// New category
			newCategoryId, err := SetNewCategory(ctx, v)
			if err != nil {
				log.Println(err)
				continue
			}
			err = SetNewItemCategory(ctx, itemId, v)
			if err != nil {
				log.Println(err)
				continue
			}
			err = SetNewSitemap(ctx, v)
			if err != nil {
				log.Println(err)
				continue
			}
			categoryIdList = append(categoryIdList, newCategoryId)
		} else {
			// Already exist category
			err = AppendItemCategory(ctx, itemId, categoryId)
			if err != nil {
				log.Println(err)
				continue
			}
			categoryIdList = append(categoryIdList, categoryId)
		}
	}
	return categoryIdList, nil
}

func SetNewCategory(ctx context.Context, categoryName string)(int, error) {
	dataCount, err := GetDynamoDataCount(ctx, DataTypeCategory)
	if err != nil {
		return 0, err
	}
	err = putCategoryData(ctx, dataCount + 1, categoryName)
	if err != nil {
		return 0, err
	}
	return dataCount + 1, nil
}

func SetNewItemCategory(ctx context.Context, itemId int, categoryName string) error {
	dataCount, err := GetDynamoDataCount(ctx, DataTypeItemCategory)
	if err != nil {
		return err
	}
	v, _ := json.Marshal([]int{itemId})
	err = putCategoryItemData(ctx, dataCount + 1, KVSData{K: categoryName, V: string(v)})
	if err != nil {
		return err
	}
	return nil
}

func AppendItemCategory(ctx context.Context, itemId int, categoryId int) error {
	itemCategoryMap, err := GetSingleItemCategoryMap(ctx, categoryId)
	if err != nil {
		log.Println(err)
		return err
	}
	var kvs KVSData
	json.Unmarshal([]byte(itemCategoryMap.Data), &kvs)
	var itemIdList []int
	json.Unmarshal([]byte(kvs.V), &itemIdList)

	isContain := false
	for _, v := range itemIdList {
		if v == itemId {
			isContain = true
			break
		}
	}
	if !isContain {
		itemIdList = append(itemIdList, itemId)
	}
	v, _ := json.Marshal(itemIdList)
	data := KVSData{
		K: kvs.K,
		V: string(v),
	}
	data_, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return UpdateItemCategoryMap(ctx, categoryId, string(data_))
}

func DeleteItemCategory(ctx context.Context, itemId int, categoryId int) error {
	itemCategoryMap, err := GetSingleItemCategoryMap(ctx, categoryId)
	if err != nil {
		log.Println(err)
		return err
	}
	var kvs KVSData
	json.Unmarshal([]byte(itemCategoryMap.Data), &kvs)
	var itemIdList []int
	json.Unmarshal([]byte(kvs.V), &itemIdList)

	var itemIdList_ []int
	for i, v := range itemIdList {
		if v == itemId {
			itemIdList_ = append(itemIdList[:i], itemIdList[i+1:]...)
			break
		}
	}
	v, _ := json.Marshal(itemIdList_)
	data := KVSData{
		K: kvs.K,
		V: string(v),
	}
	data_, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return UpdateItemCategoryMap(ctx, categoryId, string(data_))
}

func SetNewSitemap(ctx context.Context, categoryName string) error {
	t := time.Now()
	dataCount, err := GetDynamoDataCount(ctx, DataTypeSitemap)
	if err != nil {
		log.Println(err)
		return err
	}
	i := 0
	if dataCount == 0 {
		sitemapData := SitemapData{
			Loc: "",
			Lastmod: t.Format(LastmodLayout),
			Changefreq: "monthly",
			Priority: "1.0",
		}
		err = putSitemapData(ctx, 1, sitemapData)
		if err != nil {
			return err
		}
		i = 1
	}
	sitemapData := SitemapData{
		Loc: "category/" + categoryName,
		Lastmod: t.Format(LastmodLayout),
		Changefreq: "monthly",
		Priority: "0.9",
	}
	err = putSitemapData(ctx, dataCount + i + 1, sitemapData)
	if err != nil {
		return err
	}
	return nil
}

func AddItem(ctx context.Context, title string, description string, image string, categoryNames []string) error {
	itemCount, err := GetDynamoDataCount(ctx, DataTypeItem)
	if err != nil {
		log.Println(err)
		return err
	}
	categoryIds, err := SetCategoryList(ctx, itemCount + 1, categoryNames)
	if err != nil {
		log.Println(err)
		return err
	}
	return putItemData(ctx, itemCount + 1,
		ItemData{
			Title: title,
			Description: description,
			Image: image,
			CategoryIds: categoryIds,
		},
	)
}

func FixItem(ctx context.Context, id int, title string, description string, image string, categoryNames []string, oldCategoryIds []int) error {
	categoryIds, err := SetCategoryList(ctx, id, categoryNames)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, v := range oldCategoryIds {
		deletedId := v
		for _, w := range categoryIds {
			if v == w {
				deletedId = 0
				break
			}
		}
		if deletedId > 0 {
			DeleteItemCategory(ctx, id, deletedId)
		}
	}
	data := ItemData{
		Title: title,
		Description: description,
		Image: image,
		CategoryIds: categoryIds,
	}
	data_, err := json.Marshal(data)
	return UpdateItem(ctx, id, string(data_))
}

func SetJsFileName(ctx context.Context, fileName string) error {
	var err error
	dataCount, err := GetDynamoDataCount(ctx, DataTypeConst)
	if err != nil {
		log.Println(err)
		return err
	} else if dataCount > 0 {
		constDynamoData, err := GetSingleConst(ctx)
		if err != nil {
			log.Println(err)
			return err
		}

		var constData ConstData
		json.Unmarshal([]byte(constDynamoData.Data), &constData)
		constData.JsFileName = fileName
		constData_, err := json.Marshal(constData)

		err = UpdateConst(ctx, string(constData_))
	} else {
		err = SetConst(ctx, "", "", fileName, "")
	}
	return err
}

func SetCssFileName(ctx context.Context, fileName string) error {
	var err error
	dataCount, err := GetDynamoDataCount(ctx, DataTypeConst)
	if err != nil {
		log.Println(err)
		return err
	} else if dataCount > 0 {
		constDynamoData, err := GetSingleConst(ctx)
		if err != nil {
			log.Println(err)
			return err
		}

		var constData ConstData
		json.Unmarshal([]byte(constDynamoData.Data), &constData)
		constData.CssFileName = fileName
		constData_, err := json.Marshal(constData)

		err = UpdateConst(ctx, string(constData_))
	} else {
		err = SetConst(ctx, "", "", "", fileName)
	}
	return err
}

func GetDynamoDBData(ctx context.Context)(string, interface{}, error) {
	result, err := scan(ctx, expression.NotEqual(expression.Name("status"), expression.Value(-1)))
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	var tableContents []interface{}
	for _, i := range result.Items {
		var item interface{}
		err := attributevalue.UnmarshalMap(i, &item)
		if err != nil {
			log.Print(err)
		} else {
			tableContents = append(tableContents, item)
		}
	}
	return os.Getenv("ITEM_TABLE_NAME"), tableContents, nil
}
