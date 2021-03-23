import configparser
import time
from selenium import webdriver
from selenium.webdriver.common.alert import Alert

config = configparser.ConfigParser()
config.read('config.ini')

URL = "https://www.hikorea.go.kr/memb/MembLoginR.pt"
USER_ID = config['user']['id']
USER_PASSWD = config['user']['passwd']


# go to login page
driver = webdriver.Chrome()
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
dates = driver.find_elements_by_xpath("//table[@class='ui-datepicker-calendar']//a")
for date in dates:
    print("day: ", date.text)
    # date.click()


# apply_btn = driver.find_element_by_class_name("btn_blue")
# apply_btn.click()

#close
time.sleep(5)
driver.quit()



