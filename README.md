# Telebot

## Description
Have you ever found yourself in a situation where you have to urgently visit the Immigration Office, but all the nearing dates are unavailable? Then this bot is what you need! This bot will keep checking HIKOREA website for earliest available dates and make a reservation on your behalf. It only needs your full name(**exactly** as it appears in your ARC)!

Live version of this bot is can be found in Telegram **@hikoreanelsonbot**\
**Disclaimer**: USE THIS BOT AT YOUR OWN RISK

## Run in Local Envinronment
1. Download [chrome webdriver](https://chromedriver.chromium.org/downloads) corresponding to your chrome [version](https://help.zenplanner.com/hc/en-us/articles/204253654-How-to-Find-Your-Internet-Browser-Version-Number-Google-Chrome) and run the binary file. 
2. In a separate terminal build and run this bot.   
``` bash
 make build
 bin/telebot
```     

## Run in Docker
``` bash
 docker-compose up
```

## Run Tests 
```bash
 make test
```
