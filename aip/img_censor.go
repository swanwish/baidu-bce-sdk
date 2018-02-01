package aip

import (
	"fmt"

	"net/http"

	"encoding/json"

	"github.com/swanwish/baidu-bce-sdk/common"
	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

const (
	IMG_CENSOR_URL = "https://aip.baidubce.com/api/v1/solution/direct/img_censor"

	SCENE_ANTIPORN = "antiporn"
	//SCENE_OCR        = "ocr"
	//SCENE_POLITICIAN = "politician"
	//SCENE_TERROR     = "terror"
	//SCENE_WEBIMAGE   = "webimage"
	//SCENE_DISGUST    = "disgust"
	//SCENE_WATERMARK  = "watermark"
	//SCENE_QUALITY    = "quality"
)

type ImgCensorParam struct {
	Image  string   `json:"image,omitempty"`
	ImgUrl string   `json:"imgUrl"`
	Scenes []string `json:"scenes"`
}

type ImgCensorResponse struct {
	Result    ImgCensorResult `json:"result"`
	LogId     int64           `json:"log_id"`
	ErrorCode string          `json:"error_code"`
	ErrorMsg  string          `json:"error_msg"`
}

type ImgCensorResult struct {
	Antiporn CheckResult `json:"antiporn"`
}

type CheckResult struct {
	Result                []CheckResultItem `json:"result"`
	Conclusion            string            `json:"conclusion"`
	LogId                 int64             `json:"log_id"`
	ConfidenceCoefficient string            `json:"confidence_coefficient"`
	ResultNum             int64             `json:"result_num"`
}

type CheckResultItem struct {
	Probability float64 `json:"probability"`
	ClassName   string  `json:"class_name"`
}

func (client *Client) CheckPornImg(imgUrl string) error {
	accessToken, err := client.GetAccessToken()
	if err != nil {
		logs.Errorf("Failed to get access token, the error is %#v", err)
		return err
	}

	//accessToken := "24.ec5d2660db05665ecb22bae42cf7f43d.2592000.1519962917.282335-10733319"
	checkUrl := fmt.Sprintf("%s?access_token=%s", IMG_CENSOR_URL, accessToken)
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	param := ImgCensorParam{
		ImgUrl: imgUrl,
		Scenes: []string{SCENE_ANTIPORN},
	}
	body, err := json.Marshal(param)
	if err != nil {
		logs.Errorf("Failed to marshal param, the error is %#v", err)
		return err
	}
	logs.Debugf("The body is %s", string(body))
	status, content, err := utils.PostUrlContent(checkUrl, body, nil)
	if err != nil {
		logs.Errorf("Failed to check image, the error is %#v", err)
		return err
	}
	if status != http.StatusOK {
		logs.Errorf("Invalid status %d", status)
		return common.ErrInvalidStatus
	}
	imgCensorResponse := ImgCensorResponse{}
	if err = json.Unmarshal(content, &imgCensorResponse); err != nil {
		logs.Errorf("Failed to unmarshal img censor response")
		return err
	}
	if imgCensorResponse.ErrorCode != "" {
		logs.Errorf("The error message is %s", imgCensorResponse.ErrorMsg)
		return common.ErrInternalError
	}
	logs.Debugf("The antiporn result is %#v", imgCensorResponse.Result.Antiporn.Result)
	//logs.Debugf("The content is %s", string(content))
	return nil
}
