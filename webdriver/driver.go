package webdriver

import (
	"errors"
	"fmt"
	"strings"
	"telebot/util"
	"time"

	"github.com/spf13/viper"
	sm "github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go.uber.org/zap"
)

const (
	baseURL          = "https://www.hikorea.go.kr/memb/MembLoginR.pt"
	applyTimeLayout  = "2006-01-02"
	cancelTimeLayout = "2006. 01. 02    15 : 04"
)

var (
	log       *zap.SugaredLogger
	driverURL string
)

func InitDriver(config *viper.Viper) {
	log = util.InitLog("driver")
	host := config.GetString("driver_host")
	port := config.GetString("driver_port")
	apiPrefix := config.GetString("driver_api_prefix")
	driverURL = fmt.Sprintf("http://%s:%s%s", host, port, apiPrefix)
	log.Infow("InitDriver", "driverURL", driverURL)
}

func login(data map[string]string) (wd sm.WebDriver, err error) {
	caps := sm.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-extensions",
			"--disable-gpu",
			"--shm-size=2g",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err = sm.NewRemote(caps, driverURL)
	if err != nil {
		return nil, err
	}
	wd.SetPageLoadTimeout(2 * time.Second)
	wd.SetImplicitWaitTimeout(2 * time.Second)
	if err := wd.Get(baseURL); err != nil {
		return nil, err
	}
	elem, err := wd.FindElement(sm.ByXPATH, "//input[@id='userId']")
	if err != nil {
		return nil, err
	}
	elem.SendKeys(data["username"])
	elem, err = wd.FindElement(sm.ByXPATH, "//input[@id='userPasswd']")
	if err != nil {
		return nil, err
	}
	elem.SendKeys(data["password"])
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_login']")
	if err != nil {
		return nil, err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@id='lang_en']")
	if err != nil {
		return nil, err
	}
	elem.Click()
	return wd, nil
}

func MakeAppointment(data map[string]string, done chan bool) (receipt string, err error) {
	wd, err := login(data)
	if err != nil {
		return "", err
	}
	defer wd.Quit()
	go func() {
		<-done
		if wd != nil {
			log.Debugw("MakeAppointment", "remove zombie process", wd.SessionID())
			wd.Quit()
		}
	}()
	elem, err := wd.FindElement(sm.ByXPATH, "//a[contains(@href, 'resv') and @class='btn_apply']")
	if err != nil {
		return "", err

	}
	elem.Click()
out1:
	for {
		elem, err = wd.FindElement(sm.ByXPATH, "//button[@class='btn_blue']")
		if err != nil {
			return "", err
		}
		elem.Click()
		elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_blue_b']")
		if err != nil {
			return "", err
		}
		elem.Click()
		elem, err = wd.FindElement(sm.ByXPATH, fmt.Sprintf("//option[@value='%s']", data["branch"]))
		if err != nil {
			return "", err
		}
		elem.Click()
		elem, err = wd.FindElement(sm.ByXPATH, fmt.Sprintf("//input[@id='%s']", data["booth"]))
		if err != nil {
			return "", err
		}
		elem.Click()
		validBooth := false
		elem, err = wd.FindElement(sm.ByXPATH, fmt.Sprintf("//input[@value='%s']", data["purpose"]))
		if err != nil {
			return "", err
		}
		elem.Click()
		if phone := getPhoneNumber(data["phone"]); phone != nil {
			elem, err = wd.FindElement(sm.ByXPATH, "//select[@id='mobileTelNo1']")
			if err != nil {
				return "", err
			}
			elem.SendKeys(phone[0])
			elem, err = wd.FindElement(sm.ByXPATH, "//input[@id='mobileTelNo2']")
			if err != nil {
				return "", err
			}
			elem.SendKeys(phone[1])
			elem, err = wd.FindElement(sm.ByXPATH, "//input[@id='mobileTelNo3']")
			if err != nil {
				return "", err
			}
			elem.SendKeys(phone[2])
		}

		for {
			elem, err = wd.FindElement(sm.ByXPATH, "//a[@id='resvYmdSelect']")
			if err != nil {
				return "", err
			}
			elem.Click()
			windows, err := wd.WindowHandles()
			if err != nil {
				return "", err
			}
			if err := wd.SwitchWindow(windows[1]); err != nil {
				return "", err
			}
		out2:
			for {
				dates, err := wd.FindElements(sm.ByXPATH, "//table[@class='ui-datepicker-calendar']//a")
				if err != nil {
					return "", err
				}
				for i := 0; i < len(dates); i++ {
					dates, err = wd.FindElements(sm.ByXPATH, "//table[@class='ui-datepicker-calendar']//a")
					if err != nil {
						return "", err
					}
					dates[i].Click()
					timeslotes, err := wd.FindElements(sm.ByXPATH, "//div[@class='select_time_table']//a")
					if err != nil {
						log.Errorw("Make Appointment", "error", err)
						if validBooth {
							wd.AcceptAlert()
							continue
						}
						return "", err
					}
					validBooth = true
					time.Sleep(500 * time.Millisecond)
					for _, timeslot := range timeslotes {
						if err := timeslot.Click(); err != nil {
							if err := wd.SwitchWindow(windows[0]); err != nil {
								return "", err
							}
							break out2
						}
						wd.AcceptAlert()
					}
				}
				elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='ui-datepicker-next ui-corner-all']")
				if err != nil {
					return "", err
				}
				elem.Click()
			}
			elem, err = wd.FindElement(sm.ByXPATH, "//input[@id='resvYmd']")
			if err != nil {
				return "", err
			}
			timeslot, err := elem.GetAttribute("value")
			if err != nil {
				return "", err
			}
			log.Debugw("MakeAppointment Found Timeslot", "username", data["username"], "timeslot", timeslot)
			if isValidTimeslot(data["prevtimeslot"], timeslot) {
				elem, err = wd.FindElement(sm.ByXPATH, "//button[@class='btn_blue']")
				if err != nil {
					return "", err
				}
				elem.Click()
				elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_blue_b']")
				if err != nil {
					return "", err
				}
				elem.Click()
				elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_blue_b']")
				if err == nil {
					log.Debugw("MakeAppointment Same Day Error", "username", data["username"], "timeslot", timeslot)
					elem.Click()
					continue out1
				}
				data["prevtimeslot"] = timeslot
				trs, err := wd.FindElements(sm.ByXPATH, "//tr")
				if err != nil {
					return "", err
				}
				var sb strings.Builder
				for _, tr := range trs {
					th, _ := tr.FindElement(sm.ByTagName, "th")
					td, _ := tr.FindElement(sm.ByTagName, "td")
					if val, _ := td.Text(); val != "" {
						key, _ := th.Text()
						sb.WriteString(key)
						sb.WriteString(": ")
						sb.WriteString(val)
						sb.WriteByte('\n')
					}
				}
				receipt = strings.Trim(sb.String(), "\n")
				log.Infow("MakeAppointment Success", "timeslot", timeslot, "receipt", receipt)
				break out1
			}
		}
	}

	return receipt, nil
}

func GetBoothes(data map[string]string) ([]string, error) {
	wd, err := login(data)
	if err != nil {
		return nil, err
	}
	defer wd.Quit()
	elem, err := wd.FindElement(sm.ByXPATH, "//a[contains(@href, 'resv') and @class='btn_apply']")
	if err != nil {
		return nil, err

	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//button[@class='btn_blue']")
	if err != nil {
		return nil, err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_blue_b']")
	if err != nil {
		return nil, err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, fmt.Sprintf("//option[@value='%s']", data["branch"]))
	if err != nil {
		return nil, err
	}
	elem.Click()
	boothes := make([]string, 0)
	elems, _ := wd.FindElements(sm.ByXPATH, "//div[@id='deskSeqList']//label")
	for _, elem := range elems {
		key, _ := elem.GetAttribute("for")
		value, _ := elem.Text()
		boothes = append(boothes, key, value)
	}
	return boothes, nil
}

func CancelPrevAppointment(data map[string]string) {
	prev, err := cancelPrevAppointment(data)
	if err != nil {
		log.Errorw("CancelPrevAppointment", "error", err)
		return
	}
	log.Debugw("CancelPrevAppointment", "prev", prev)
}

func cancelPrevAppointment(data map[string]string) (prev string, err error) {
	wd, err := login(data)
	if err != nil {
		return "", err
	}
	defer wd.Quit()
	elem, err := wd.FindElement(sm.ByXPATH, "//a[@title='Reserve Visit Status']")
	if err != nil {
		return "", err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//option[@value='RS']")
	if err != nil {
		return "", err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//button[@class='btn_search']")
	if err != nil {
		return "", err
	}
	elem.Click()
	trs, err := wd.FindElements(sm.ByXPATH, "//div[@class='grp_table scroll_x']//tbody//tr")
	if err != nil {
		return "", err
	}
	if len(trs) < 2 {
		return "", errors.New("no previous appointment exists")
	}
	var prevA sm.WebElement
	for _, tr := range trs {
		tds, err := tr.FindElements(sm.ByTagName, "td")
		if err != nil {
			return "", err
		}
		a, err := tds[1].FindElement(sm.ByTagName, "a")
		if err != nil {
			return "", err
		}
		curr, _ := a.Text()
		if isLater(curr, prev) {
			prevA = a
			prev = curr
		}
	}
	prevA.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@id='btn_cencelResv']")
	if err != nil {
		return "", err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_blue_b']")
	if err != nil {
		return "", err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_blue_b']")
	if err != nil {
		return "", err
	}
	elem.Click()
	return prev, nil
}

func isValidTimeslot(prev, next string) bool {
	nexttime, ok := parse(next)
	if !ok {
		return false
	}
	prevtime, ok := parse(prev)
	if !ok {
		return true
	}
	return prevtime.After(nexttime)
}

func parse(timeslot string) (time.Time, bool) {
	timeslot = strings.TrimSpace(timeslot)
	if timeslot == "" {
		return time.Time{}, false
	}
	timeslot = strings.Fields(timeslot)[0]
	t, err := time.Parse(applyTimeLayout, timeslot)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func getPhoneNumber(input string) []string {
	n := 3
	number := make([]string, 0)
	temp := make([]byte, 0)
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' {
			temp = append(temp, input[i])
			if len(temp) == n {
				number = append(number, string(temp))
				temp = make([]byte, 0)
				n = 4
			}
		} else if input[i] == '-' {
			continue
		} else {
			return nil
		}
	}
	if len(temp) != 0 || len(number) != 3 {
		return nil
	}
	return number
}

func isLater(curr, prev string) bool {
	if prev == "" {
		return true
	}
	prevtime, _ := time.Parse(cancelTimeLayout, prev)
	currtime, _ := time.Parse(cancelTimeLayout, curr)
	return currtime.After(prevtime)
}
