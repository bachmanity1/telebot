package webdriver

import (
	"fmt"
	"strings"
	"telebot/util"
	"time"

	sm "github.com/tebeka/selenium"
	"go.uber.org/zap"
)

const (
	port    = 9515
	baseURL = "https://www.hikorea.go.kr/memb/MembLoginR.pt"
	layout  = "2006-01-02"
)

var log *zap.SugaredLogger

func init() {
	log = util.InitLog("driver")
}

func MakeAppointment(data map[string]string) (timeslot string, err error) {
	caps := sm.Capabilities{
		"browserName": "chrome",
	}
	wd, err := sm.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		return "", err
	}
	defer wd.Quit()
	if err := wd.Get(baseURL); err != nil {
		return "", err
	}
	elem, err := wd.FindElement(sm.ByXPATH, "//input[@id='userId']")
	if err != nil {
		return "", err
	}
	elem.SendKeys(data["username"])
	elem, err = wd.FindElement(sm.ByXPATH, "//input[@id='userPasswd']")
	if err != nil {
		return "", err
	}
	elem.SendKeys(data["password"])
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_login']")
	if err != nil {
		return "", err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@id='lang_en']")
	if err != nil {
		return "", err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[contains(@href, 'resv') and @class='btn_apply']")
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
		elem, err = wd.FindElement(sm.ByXPATH, fmt.Sprintf("//input[@id='%s']", data["subbranch"]))
		if err != nil {
			return "", err
		}
		elem.Click()
		elem, err = wd.FindElement(sm.ByXPATH, "//input[@name='deskSeq']")
		if err != nil {
			return "", err
		}
		elem.Click()
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
			time.Sleep(2 * time.Second)
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
						return "", err
					}
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
			timeslot, err = elem.GetAttribute("value")
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
				log.Infow("MakeAppointment Success", "username", data["username"], "timeslot", timeslot)
				break out1
			}
		}
	}

	time.Sleep(5 * time.Second)
	return timeslot, nil
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
	t, err := time.Parse(layout, timeslot)
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

func GetSubBranches(data map[string]string) (map[string]string, error) {
	caps := sm.Capabilities{
		"browserName": "chrome",
	}
	wd, err := sm.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		return nil, err
	}
	defer wd.Quit()
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
	elem, err = wd.FindElement(sm.ByXPATH, "//a[contains(@href, 'resv') and @class='btn_apply']")
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
	subbranches := make(map[string]string)
	elems, _ := wd.FindElements(sm.ByXPATH, "//div[@id='deskSeqList']//label")
	for _, elem := range elems {
		key, _ := elem.GetAttribute("for")
		value, _ := elem.Text()
		subbranches[key] = value
	}
	return subbranches, nil
}
