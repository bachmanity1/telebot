package webdriver

import (
	"fmt"
	"telebot/handler"
	"telebot/util"
	"time"

	"github.com/spf13/viper"
	sm "github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go.uber.org/zap"
)

const baseURL = "https://www.hikorea.go.kr/memb/MembLoginR.pt"

var log *zap.SugaredLogger

func InitData(config *viper.Viper) error {
	log = util.InitLog("driver")
	host := config.GetString("driver_host")
	port := config.GetString("driver_port")
	apiPrefix := config.GetString("driver_api_prefix")
	driverURL := fmt.Sprintf("http://%s:%s%s", host, port, apiPrefix)
	log.Infow("InitDriver", "driverURL", driverURL)

	caps := sm.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--window-size=1920,1080",
			"--no-sandbox",
			"--disable-extensions",
			"--disable-gpu",
			"--dns-prefetch-disable",
			"--shm-size=2g",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := sm.NewRemote(caps, driverURL)
	if err != nil {
		return err
	}
	defer wd.Quit()
	wd.SetPageLoadTimeout(2 * time.Second)
	wd.SetImplicitWaitTimeout(2 * time.Second)

	userID := config.GetString("hikorea_user_id")
	userPasswd := config.GetString("hikorea_user_passwd")
	boothes, err := getBoothes(wd, userID, userPasswd)
	if err != nil {
		return err
	}
	handler.MakeBoothMarkup(boothes)
	return nil
}

func getBoothes(wd sm.WebDriver, userID, userPasswd string) (map[string][]string, error) {
	boothes := make(map[string][]string)
	if err := wd.Get(baseURL); err != nil {
		return nil, err
	}
	elem, err := wd.FindElement(sm.ByXPATH, "//input[@id='userId']")
	if err != nil {
		return nil, err
	}
	elem.SendKeys(userID)
	elem, err = wd.FindElement(sm.ByXPATH, "//input[@id='userPasswd']")
	if err != nil {
		return nil, err
	}
	elem.SendKeys(userPasswd)
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

	branches, err := wd.FindElements(sm.ByXPATH, "//select[@id='orgnCd']//option")
	if err != nil {
		return nil, err
	}
	for _, branch := range branches[1:] {
		branch.Click()
		branchName, err := branch.GetAttribute("value")
		if err != nil {
			return nil, err
		}
		boothz := make([]string, 0)
		elems, err := wd.FindElements(sm.ByXPATH, "//div[@id='deskSeqList']//label")
		if err != nil {
			return nil, err
		}
		for _, elem := range elems {
			key, _ := elem.GetAttribute("for")
			value, _ := elem.Text()
			boothz = append(boothz, key, value)
		}
		boothes[branchName] = boothz
		log.Debugw("getBoothes", "branch", branch, "boothes", boothz)
	}
	return boothes, nil
}
