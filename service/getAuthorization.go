package service

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"pointMall/setting"
)

func GetAuthorization() string {
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
	params := url.Values{}
	params.Set("client_id", "965666444")
	params.Set("response_type", "code")
	params.Set("redirect_uri", "https://dgtmall.dgtis.com/mall/#/auth")
	Url, err := url.Parse("https://dgt.dgtis.com/oneportal/oauth2api/authorizeMall.if")
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
	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))
	// fmt.Println(resp.Header)
	// fmt.Println(resp.Header.Get("Authorization"))
	// var res result
	// _ = json.Unmarshal(body, &res)
	// fmt.Printf("%#v", res)
	token := resp.Header.Get("Authorization")
	return token
}
