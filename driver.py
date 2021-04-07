import configparser
import time
from dateutil import parser
from datetime import datetime, timedelta
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.alert import Alert
from selenium.common.exceptions import NoSuchWindowException

config = configparser.ConfigParser()
config.read('config.ini')

URL = "https://www.hikorea.go.kr/memb/MembLoginR.pt"
USER_ID = config['user']['id']
USER_PASSWD = config['user']['passwd']
USER_PHONE_NUMBER = config['user']['phone'].split('-')

chrome_options = Options()
#chrome_options.add_argument("--disable-extensions")
#chrome_options.add_argument("--disable-gpu")
#chrome_options.add_argument("--no-sandbox") # linux only
chrome_options.add_argument("--headless")
# chrome_options.headless = True # also works
driver = webdriver.Chrome(options=chrome_options)
earliest_date = datetime.now() + timedelta(days=300)
try:
    # go to login page
    driver.get(URL)
    # enter login credentials
    driver.find_element_by_id("userId").send_keys(USER_ID)
    driver.find_element_by_id("userPasswd").send_keys(USER_PASSWD)
    driver.find_element_by_class_name("btn_login").click()
    # go to reservation page
    driver.find_element_by_xpath("//a[contains(@href, 'resv') and @class='btn_apply']").click()
    driver.find_element_by_class_name("btn_blue").click()
    driver.find_element_by_class_name("btn_blue_b").click()
   
    main_window = driver.window_handles[0]
    blank = True
    # find earliest available date and time
    while True:
        if blank:
            driver.find_element_by_xpath("//input[@name='deskSeq']").click()
            driver.find_element_by_xpath("//input[@name='selBusiType1_1']").click()
            driver.find_element_by_xpath("//select[@id='mobileTelNo1']").send_keys(USER_PHONE_NUMBER[0])
            driver.find_element_by_xpath("//input[@id='mobileTelNo2']").send_keys(USER_PHONE_NUMBER[1])
            driver.find_element_by_xpath("//input[@id='mobileTelNo3']").send_keys(USER_PHONE_NUMBER[2])
        
        driver.find_element_by_id("resvYmdSelect").click()
        popup_window = driver.window_handles[1]
        driver.switch_to.window(popup_window)
        time.sleep(2)
        try:
            while True:
                dates_len = len(driver.find_elements_by_xpath("//table[@class='ui-datepicker-calendar']//a"))
                for i in range(dates_len):
                    date = driver.find_elements_by_xpath("//table[@class='ui-datepicker-calendar']//a")[i]
                    date.click()
                    times = driver.find_elements_by_xpath("//div[@class='select_time_table']//a")
                    for t in times:
                        t.click()
                        try:
                            alert = Alert(driver)
                            alert.accept()
                        except:
                            pass
                            # print(t.text + " : no alert message")
                driver.find_element_by_xpath("//a[@class='ui-datepicker-next ui-corner-all']").click()
        except NoSuchWindowException:
            driver.switch_to.window(main_window)
        timeslot = driver.find_element_by_id("resvYmd").get_attribute('value')
        print("found available timeslot: ", timeslot)
        tempdate = parser.parse(timeslot.split()[0])
        if tempdate < earliest_date:
            earliest_date = tempdate 
            driver.find_element_by_class_name("btn_blue").click()
            driver.find_element_by_class_name("btn_blue_b").click()
            print("succesfully made an resvervation for: ", timeslot)
            driver.back()
            blank = True
        else:
            blank = False
except Exception as e:
    print("something went wrong... ", str(e))
finally:
    # close all windows
    time.sleep(5)
    driver.quit()
