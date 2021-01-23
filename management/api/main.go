package main

import (
	"fmt"
	"log"
	"context"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tanaka-takurou/serverless-simple-cms-go/management/api/controller"
)

type UserResponse struct {
	Name     string `json:"name"`
	Token    string `json:"token"`
	ImgUrl   string `json:"imgurl"`
}

type ContentResponse struct {
	Const            controller.DynamoData   `json:"constData"`
	ItemList         []controller.DynamoData `json:"itemList"`
	CategoryList     []controller.DynamoData `json:"categoryList"`
	CategoryItemList []controller.DynamoData `json:"categoryItemMap"`
}

type DataResponse struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	controller.InitConfig(ctx)
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "login" :
			err = checkParameters(d, []string{"name","pass"})
			if err == nil {
				t, e := controller.Login(ctx, d["name"], d["pass"])
				if e == nil {
					jsonBytes, _ = json.Marshal(UserResponse{Name: d["name"], Token: t, ImgUrl: ""})
				} else {
					err = e
				}
			}
		case "get_user" :
			n, e := getUser(ctx, d)
			if e == nil {
				jsonBytes, _ = json.Marshal(UserResponse{Name: n, Token: "", ImgUrl: ""})
			} else {
				err = e
			}
		case "change_pass" :
			err = checkParameters(d, []string{"token","pass","newpass"})
			if err == nil {
				err = controller.ChangePass(ctx, d["token"], d["pass"], d["newpass"])
				jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: ""})
			}
		case "logout" :
			err = checkParameters(d, []string{"token"})
			if err == nil {
				err = controller.Logout(ctx, d["token"])
				if err == nil {
					jsonBytes, _ = json.Marshal(UserResponse{Name: "ok", Token: "Logout", ImgUrl: ""})
				}
			}
		case "sign_up" :
			err = checkParameters(d, []string{"name","pass","mail"})
			if err == nil {
				err = controller.Signup(ctx, d["name"], d["pass"], d["mail"])
				jsonBytes, _ = json.Marshal(UserResponse{Name: d["name"], Token: "", ImgUrl: ""})
			}
		case "confirm_signup" :
			err = checkParameters(d, []string{"name","code"})
			if err == nil {
				err = controller.ConfirmSignup(ctx, d["name"], d["code"])
				jsonBytes, _ = json.Marshal(UserResponse{Name: d["name"], Token: "", ImgUrl: ""})
			}
		case "upload_img" :
			_, err = getUser(ctx, d)
			if err == nil {
				err := checkParameters(d, []string{"filename","filedata"})
				if err == nil {
					imgUrl, _ := controller.UploadImage(ctx, d["filename"], d["filedata"])
					jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: imgUrl})
				}
			}
		case "get_const" :
			_, err = getUser(ctx, d)
			if err == nil {
				var constData controller.DynamoData
				constData, err = controller.GetSingleConst(ctx)
				jsonBytes, _ = json.Marshal(ContentResponse{
					Const: constData,
					ItemList: []controller.DynamoData{},
					CategoryList: []controller.DynamoData{},
					CategoryItemList: []controller.DynamoData{},
				})
			}
		case "set_const" :
			_, err = getUser(ctx, d)
			if err == nil {
				err := checkParameters(d, []string{"title", "image"})
				if err == nil {
					err = controller.SetConst(ctx, d["title"], d["image"], "", "")
					jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: ""})
				}
			}
		case "set_sample" :
			_, err = getUser(ctx, d)
			if err == nil {
				err = controller.SetSampleData(ctx)
				jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: ""})
			}
		case "get_item_category_list" :
			_, err = getUser(ctx, d)
			if err == nil {
				var itemList []controller.DynamoData
				itemList, err = controller.GetItemList(ctx)
				if err == nil {
					var categoryList []controller.DynamoData
					categoryList, err = controller.GetCategoryList(ctx)
					jsonBytes, _ = json.Marshal(ContentResponse{
						Const: controller.DynamoData{},
						ItemList: itemList,
						CategoryList: categoryList,
						CategoryItemList: []controller.DynamoData{},
					})
				}
			}
		case "get_category_list" :
			_, err = getUser(ctx, d)
			if err == nil {
				var categoryList []controller.DynamoData
				categoryList, err = controller.GetCategoryList(ctx)
				jsonBytes, _ = json.Marshal(ContentResponse{
					Const: controller.DynamoData{},
					ItemList: []controller.DynamoData{},
					CategoryList: categoryList,
					CategoryItemList: []controller.DynamoData{},
				})
			}
		case "add_item" :
			_, err = getUser(ctx, d)
			if err == nil {
				err := checkParameters(d, []string{"title","description","image","categories"})
				if err == nil {
					categoryNames := make([]string, 3)
					json.Unmarshal([]byte(d["categories"]), &categoryNames)
					err = controller.AddItem(ctx, d["title"], d["description"], d["image"], categoryNames)
					jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: ""})
				}
			}
		case "fix_item" :
			_, err = getUser(ctx, d)
			if err == nil {
				err := checkParameters(d, []string{"id","title","description","image","categories","old_categories"})
				if err == nil {
					itemId, e := strconv.Atoi(d["id"])
					if e == nil {
						categoryNames := make([]string, 3)
						json.Unmarshal([]byte(d["categories"]), &categoryNames)
						oldCategoryIds := make([]int, 3)
						json.Unmarshal([]byte(d["old_categories"]), &oldCategoryIds)
						err = controller.FixItem(ctx, itemId, d["title"], d["description"], d["image"], categoryNames, oldCategoryIds)
						jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: ""})
					} else {
						err = e
					}
				}
			}
		case "fix_js" :
			_, err = getUser(ctx, d)
			if err == nil {
				err := checkParameters(d, []string{"js_text"})
				if err == nil {
					err = controller.FixJS(ctx, d["js_text"])
					jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: ""})
				}
			}
		case "fix_css" :
			_, err = getUser(ctx, d)
			if err == nil {
				err := checkParameters(d, []string{"css_text"})
				if err == nil {
					err = controller.FixCSS(ctx, d["css_text"])
					jsonBytes, _ = json.Marshal(UserResponse{Name: "", Token: "", ImgUrl: ""})
				}
			}
		case "get_dynamodb_data" :
			_, err = getUser(ctx, d)
			if err == nil {
				name := ""
				var data interface{}
				name, data, err = controller.GetDynamoDBData(ctx)
				if err == nil {
					jsonBytes, _ = json.Marshal(DataResponse{Name: name, Data: data})
				}
			}
		case "get_s3_data" :
			_, err = getUser(ctx, d)
			if err == nil {
				name := ""
				var data interface{}
				name, data, err = controller.GetS3Data(ctx)
				if err == nil {
					jsonBytes, _ = json.Marshal(DataResponse{Name: name, Data: data})
				}
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		jsonBytes, _ = json.Marshal(ErrorResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body: string(jsonBytes),
		}, nil
	}
	responseBody := ""
	if len(jsonBytes) > 0 {
		responseBody = string(jsonBytes)
	}
	return Response {
		StatusCode: http.StatusOK,
		Body: responseBody,
	}, nil
}

func getUser(ctx context.Context, d map[string]string)(string, error) {
	if t, ok := d["token"]; ok {
		return controller.GetUser(ctx, t)
	}
	return "", fmt.Errorf("Error: %s", "No token.")
}

func checkParameters(d map[string]string, targets []string) error {
	for _, v := range targets {
		w, ok := d[v]
		if !ok || len(w) < 1 {
			return fmt.Errorf("Error: %s is nil or empty.", v)
		}
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
