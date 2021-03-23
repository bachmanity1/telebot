import configparser
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
user_id = driver.find_element_by_id("userId")
user_id.send_keys(USER_ID)
user_passwd = driver.find_element_by_id("userPasswd")
user_passwd.send_keys(USER_PASSWD)
login = driver.find_element_by_class_name("btn_login")
login.click()
# go to reservation page
resv_apply = driver.find_element_by_xpath("//a[contains(@href, 'resv') and @class='btn_apply']")
resv_apply.click()
resv_apply = driver.find_element_by_class_name("btn_blue")
resv_apply.click()
resv_apply = driver.find_element_by_class_name("btn_blue_b")
resv_apply.click()

