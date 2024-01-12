package utils

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"os"
)

type RspData struct {
	Data      string `json:"Data"`
	RequestID string `json:"RequestId"`
}

type Data struct {
	Result struct {
		PositiveProb float64 `json:"positive_prob"`
		NegativeProb float64 `json:"negative_prob"`
		NeutralProb  float64 `json:"neutral_prob"`
		Sentiment    string  `json:"sentiment"`
	} `json:"result"`

	Success  bool   `json:"success"`
	TracerID string `json:"tracerId"`
}

func SentimentAnalysis(text string) Data {

	AccessKeyId := os.Getenv("NLP_AK")
	AccessKeySecret := os.Getenv("NLP_SK")
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", AccessKeyId, AccessKeySecret)
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Domain = "alinlp.cn-hangzhou.aliyuncs.com"
	request.Version = "2020-06-29"
	// 因为是RPC接口，因此需指定ApiName(Action)
	request.ApiName = "GetSaChGeneral"
	request.QueryParams["ServiceCode"] = "alinlp"
	request.QueryParams["Text"] = text
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		code := response.GetHttpStatus()
		fmt.Println(err.Error(), response.GetHttpContentString(), code)
	}

	res := response.GetHttpContentString()
	fmt.Println(res)

	var data RspData

	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println(err.Error())
	}

	// string to map
	var result Data
	err = json.Unmarshal([]byte(data.Data), &result)
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}
