# Reservation Bot
This program will keep looking for earliest available timeslot and make an reservation on your behalf.  
## How to run
1. Download [chrome webdriver](https://chromedriver.chromium.org/downloads) corresponding to your chrome [version](https://help.zenplanner.com/hc/en-us/articles/204253654-How-to-Find-Your-Internet-Browser-Version-Number-Google-Chrome) and put the binary file in your system *PATH*. 
2. Pull this repo to your local machine and run below commands. You also have to enter your user credentials in a *config.ini* file.   
``` bash
 python3 -m venv .venv37
 source .venv37/bin/activate
 pip install -r requirements.txt
 python driver.py
```     
Program may end up making **TOO** many reservations. (manually check reservation status page)

*TODO*: rewrite in golang and add telegram bot support
