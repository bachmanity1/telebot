import configparser
import time
from selenium import webdriver
from selenium.webdriver.common.alert import Alert
from selenium.webdriver.remote import switch_to
from selenium.common.exceptions import NoSuchWindowException
config = configparser.ConfigParser()
config.read('config.ini')

URL = "https://www.hikorea.go.kr/memb/MembLoginR.pt"
USER_ID = config['user']['id']
USER_PASSWD = config['user']['passwd']

driver = webdriver.Chrome()
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
    # choose time
    main_window = driver.window_handles[0]
    driver.find_element_by_xpath("//input[@name='deskSeq']").click()
    driver.find_element_by_xpath("//input[@name='selBusiType1_1']").click()
    driver.find_element_by_id("resvYmdSelect").click()
    popup_window = driver.window_handles[1]
    driver.switch_to.window(popup_window)
    time.sleep(2)
    # find earliest available date and time
    try:
        while True:
            dates_len = len(driver.find_elements_by_xpath("//table[@class='ui-datepicker-calendar']//a"))
            for i in range(dates_len):
                date = driver.find_elements_by_xpath("//table[@class='ui-datepicker-calendar']//a")[i]
                print("day: ", date.text)
                date.click()
                times = driver.find_elements_by_xpath("//div[@class='select_time_table']//a")
                for t in times:
                    t.click()
                    try:
                        alert = Alert(driver)
                        alert.accept()
                    except:
                        print(t.text + " : no alert message")
            driver.find_element_by_xpath("//a[@class='ui-datepicker-next ui-corner-all']").click()
    except NoSuchWindowException as e:
        print("found available time slot")
        driver.switch_to.window(main_window)
    print("window title: ", driver.title)
    driver.find_element_by_class_name("btn_blue").click()
    driver.find_element_by_class_name("btn_blue_b").click()
    print("succesfully made an resvervation")
except Exception as e:
    print("something went wrong... ", str(e))
finally:
    # close all windows
    time.sleep(5)
    driver.quit()

# apply_btn = driver.find_element_by_class_name("btn_blue")
# apply_btn.click()




