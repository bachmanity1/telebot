package scraper

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

func InitRequest() {
	client = resty.New()
	client.SetDebug(true)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(1 * time.Minute)
	client.SetHeaders(map[string]string{
		"Host":                      "www.hikorea.go.kr",
		"Connection":                "keep-alive",
		"Cache-Control":             "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"Origin":                    "https://www.hikorea.go.kr",
		"Content-Type":              "application/x-www-form-urlencoded",
		"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Sec-GPC":                   "1",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-User":            "?1",
		"Sec-Fetch-Dest":            "document",
		"Referer":                   "https://www.hikorea.go.kr/resv/ResvIdntC.pt",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "en-US,en;q=0.9",
		"Cookie":                    "WMONID=ZYMMcIU8ZWG; JSESSIONID=uo2ad0fiar7f6hJz47bEaMnWlGnYc1JrBAS1beKiyoVbaYYeJ1x9APDd4GgLvRGh.amV1c19kb21haW4vaGlrb3JlYS1lZ292MQ==",
	})
}

func MakeAppointment(requestData map[string]string) (receipt string, err error) {
	resp, err := client.R().
		SetFormData(map[string]string{
			"userId":          "hikorea_2",
			"operDeskCnt":     "7",
			"targetSeq":       "39",
			"resvDt":          "20210430",
			"selBusiTypeList": "F01",
			"orgnCd":          "1270667",
			"deskSeq":         "702",
			"visiPurp":        "AA",
			"resvTime1":       "1540_1550",
			"resvNm":          "Nelson Bighetti",
			"selBusiType1_1":  "F01",
			"resvPasswd":      "1111",
			"resvYmd":         "2021-04-30 15:40~15:50",
			"visiPurpTxt":     "",
			"TRAN_TYPE":       "ComSubmit",
		}).Post("https://www.hikorea.go.kr/resv/ResvC.pt")
	if err != nil {
		log.Debugw("MakeAppointment", "error", err)
		return "", err
	}
	fmt.Println(resp.RawBody())
	return "", nil
}
