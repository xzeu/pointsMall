package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"pointMall/encrypt"
	"pointMall/setting"
)

type Answers struct {
	Answers    []map[string]interface{} `json:"answers"`
	Id         int                      `json:"id"`
	SortId     int                      `json:"sortId"`
	Sort       string                   `json:"sort"`
	Ask        string                   `json:"ask"`
	Answer     []map[string]interface{} `json:"answer"`
	UserAnswer string                   `json:"userAnswer"`
}

type Questions struct {
	RoomId    string    `json:"roomId"`
	Questions []Answers `json:"questions"`
}

type Resp struct {
	Answers []interface{} `json:"answers"`
	RoomId  string        `json:"roomId"`
}

type respondJson struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	Data   string `json:"data"`
	Errors string `json:"errors"`
}

func httpget(token string, address string, params url.Values) []byte {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// fmt.Println("via", len(via))
			// fmt.Println("via", via[0].Header)
			// if len(via) >= 3 {
			//  return fmt.Errorf("stapped after %d redirects", len(via))
			// }
			// return nil
			return errors.New("forbidden redirect")
		},
		Jar: setting.GCurCookieJar,
	}

	Url, err := url.Parse(address)
	if err != nil {
		panic(err)
	}
	// params.Set("name", "xzeu")
	// params.Set("age", "23")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	// fmt.Println(urlPath) // https://httpbin.org/get?age=23&name=zhaofan
	req, err := http.NewRequest("GET", urlPath, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)
	if err != nil {
		panic(err)
	}
	resp, _ := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))
	// fmt.Println(resp.Header)
	// fmt.Println(resp.Header.Get("Authorization"))
	// var res result
	// jsonstr := json.Unmarshal(body, &res)
	// // fmt.Printf("%#v", res)
	// fmt.Println(string(body))
	return body
}

func DailyQA(token string) (err error, askJson []byte) {
	asklist := httpget(token, "https://xiaoyou.dgtis.com/admin/answer/room", nil)
	var ask respondJson
	err = json.Unmarshal(asklist, &ask)
	if err != nil {
		panic(err)
	}
	if ask.Errno != -1 {

		// fmt.Println(ask.Data)
		res, err := base64.StdEncoding.DecodeString(ask.Data)
		if err != nil {
			log.Fatal(err)
		}
		str := encrypt.NewAes.EcbDecrypt(res)
		// fmt.Println( )}
		var questions Questions
		err = json.Unmarshal(str, &questions)
		// fmt.Println(questions)
		roomId := questions.RoomId
		// fmt.Println(roomId)
		var respStr Resp
		respStr.RoomId = roomId
		for _, item := range questions.Questions {
			// fmt.Println(item.Ask)
			// fmt.Println(item.Answers)
			for _, i := range item.Answers {
				// fmt.Println(i["answer"], i["right"])
				if i["right"] == true {
					respStr.Answers = append(respStr.Answers, i["answer"])
				}
			}

		}
		respJson, err := json.Marshal(&respStr)
		if err != nil {
			return err, nil
		}
		fmt.Println(string(respJson))
		return nil, respJson
	}
	return errors.New(ask.Errmsg), nil
}
