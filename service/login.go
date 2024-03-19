package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"pointMall/setting"
	"strings"
)

func Login(userName string, passwd string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("stapped after %d redirects", len(via))
			}
			return nil
		},
		Jar: setting.GCurCookieJar,
	}
	pwd := passwd
	hash := md5.Sum([]byte(pwd))
	md5Str := hex.EncodeToString(hash[:])
	upperMd5 := strings.ToUpper(md5Str)
	params := url.Values{}
	params.Set("loginId", userName)
	params.Set("loginPwd", upperMd5)
	params.Set("pwd", "")
	//fmt.Println(params)
	data := params.Encode()
	//fmt.Println(data)
	// bytesData, _ := json.Marshal(data)
	requested, _ := http.NewRequest("POST", "https://dgt.dgtis.com/oneportal/loginSubmit", strings.NewReader(data))
	requested.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (HTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0")
	requested.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Authorization", "3")
	respond, _ := client.Do(requested)
	// body, _ := io.ReadAll(resp.Body)
	respondCookie := respond.Cookies()
	//全局保存
	setting.GCurCookies = setting.GCurCookieJar.Cookies(requested.URL)
	fmt.Println("Cookie:", respondCookie)
	// fmt.Println(string(body))
	// fmt.Println(resp.Request.URL)

}
