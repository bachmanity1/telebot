package scraper

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"telebot/util"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jpillora/backoff"
)

const (
	successMessage   = "방문예약신청완료"
	workdayStartTime = 9 * time.Hour
	dayLength        = 24 * time.Hour
	workdayLength    = 9 * time.Hour
	slotLength       = 10 * time.Minute
	windowLength     = 30 * 24 * time.Hour
	coronaTime       = 50
)

var client *resty.Client

func InitRequest() {
	log = util.InitLog("scraper")
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
	})
}

func MakeAppointment(req map[string]string, done chan bool) bool {
	b := &backoff.Backoff{
		Min:    1 * time.Minute,
		Max:    1 * time.Hour,
		Factor: 2,
		Jitter: true,
	}
	startDate, endDate := getDateWindow(req["resvDt"])
	for {
		select {
		case <-done:
			return false
		default:
			date := startDate
			for date.Before(endDate) {
				startTime, endTime := date, date.Add(workdayLength)
				for startTime.Before(endTime) {
					from, to := startTime, startTime.Add(slotLength)
					startTime = to
					if notValidDate(from) {
						continue
					}
					req["resvDt"] = from.Format("20060102")
					x := from.Format("1504")
					y := to.Format("1504")
					req["resvTime1"] = fmt.Sprintf("%s_%s", x, y)
					x = from.Format("2006-01-02 15:04")
					y = to.Format("15:04")
					req["resvYmd"] = fmt.Sprintf("%s~%s", x, y)

					if ok := sendRequest(req); ok {
						log.Debugw("MakeAppointment Success", "date", req["resvYmd"])
						return true
					}
				}
				date = date.Add(dayLength)
			}
			time.Sleep(b.Duration())
		}
	}
}

func sendRequest(req map[string]string) bool {
	response, err := client.R().
		SetFormData(map[string]string{
			"userId":          "hikorea_2",
			"resvDt":          req["resvDt"],
			"selBusiTypeList": "F01",
			"orgnCd":          req["branch"],
			"deskSeq":         req["booth"],
			"visiPurp":        "AA",
			"resvTime1":       req["resvTime1"],
			"resvNm":          req["name"],
			"selBusiType1_1":  "F01",
			"mobileTelNo1":    req["phone1"],
			"mobileTelNo2":    req["phone2"],
			"mobileTelNo3":    req["phone3"],
			"resvPasswd":      "1111",
			"resvYmd":         req["resvYmd"],
			"TRAN_TYPE":       "ComSubmit",
		}).Post("https://www.hikorea.go.kr/resv/ResvC.pt")
	if err != nil || response.StatusCode() != http.StatusOK {
		log.Debugw("MakeAppointment", "error", err)
		return false
	}
	body := string(response.Body())
	return strings.Contains(body, successMessage)
}

func getDateWindow(prevDate string) (time.Time, time.Time) {
	now := time.Now().Add(dayLength)
	layout := "20060102"
	startDate, _ := time.Parse(layout, now.Format(layout))
	startDate = startDate.Add(workdayStartTime)
	endDate, err := time.Parse(layout, prevDate)
	if err != nil {
		endDate = startDate.Add(windowLength)
	}
	return startDate, endDate
}

func notValidDate(date time.Time) bool {
	day := date.Weekday()
	minute := date.Minute()
	return minute == coronaTime || day == time.Saturday || day == time.Sunday
}
