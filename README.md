# Reservation Bot

## How to run code
First, make sure to download [chrome webdriver](https://chromedriver.chromium.org/downloads) corresponding to your chrome [version](https://help.zenplanner.com/hc/en-us/articles/204253654-How-to-Find-Your-Internet-Browser-Version-Number-Google-Chrome) and put the binary file in your system *PATH*, then enter your user credentials in a *config.ini* file, pull this repo to your local machine and run below commands. The program   
will keep looking for earliest available timeslot and make an reservation on your behalf. (it may end up making TOO many reservations)
``` 
    python3 -m venv .venv37
    source .venv37/bin/activate
    pip install -r requirements.txt
    python driver.py
```     
