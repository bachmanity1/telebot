package webdriver

import (
	"fmt"
	"time"

	sm "github.com/tebeka/selenium"
)

const (
	port    = 9515
	baseURL = "https://www.hikorea.go.kr/memb/MembLoginR.pt"
)

func MakeAppointment(data map[string]string) error {
	caps := sm.Capabilities{
		"browserName": "chrome",
	}
	wd, err := sm.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		return err
	}
	// Navigate to the simple playground interface.
	if err := wd.Get(baseURL); err != nil {
		return err
	}
	elem, err := wd.FindElement(sm.ByXPATH, "//input[@id='userId']")
	if err != nil {
		return err
	}
	elem.SendKeys(data["username"])
	elem, err = wd.FindElement(sm.ByXPATH, "//input[@id='userPasswd']")
	if err != nil {
		return err
	}
	elem.SendKeys(data["password"])
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_login']")
	if err != nil {
		return err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[contains(@href, 'resv') and @class='btn_apply']")
	if err != nil {
		return err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//button[@class='btn_blue']")
	if err != nil {
		return err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//a[@class='btn_blue_b']")
	if err != nil {
		return err
	}
	elem.Click()
	elem, err = wd.FindElement(sm.ByXPATH, "//option[@value='1270700']")
	if err != nil {
		return err
	}
	elem.Click()

	time.Sleep(5 * time.Second)
	defer wd.Quit()
	return err
}
