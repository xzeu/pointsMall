package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SignIn(token string) {
	//执行签到
	sginUrl := "https://xiaoyou.dgtis.com/admin/mall-sign/sgin"

	req, _ := http.NewRequest("POST", sginUrl, nil)

	req.Header.Add("Authorization", token)
	req.Header.Add("Sec-Ch-Ua", "\"Microsoft Edge\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", "macOS")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	respondStr, _ := io.ReadAll(res.Body)
	var ask respondJson
	err := json.Unmarshal(respondStr, &ask)
	if err != nil {
		panic(err)
	}
	// fmt.Println(res)
	fmt.Println(ask.Errmsg)
}
