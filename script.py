#!/usr/bin/python
import random
import time
import requests

url = "http://13.114.55.155/api/luckyBetBenchMark"

luckyNumber = random.randint(0,9)
amount = random.randint(1,5)

payload = "address=IOST2g5LzaXkjAwpxCnCm29HK69wdbyRKbfG4BQQT7Yuqk57bgTFkY&privKey=319xGCaLZP5D4sAVCEX4LDAMgzaZ3LJiXgCVxB8y1igTmUCkHj6DJRCH4C8myor1P3rZHttFneApzznHqvqqTpiu&betAmount={0}00000000&luckyNumber={1}&gcaptcha=ahaha".format(amount, luckyNumber)


headers = {
    'Content-Type': "application/x-www-form-urlencoded",
    'Cache-Control': "no-cache",
    'Postman-Token': "6b809870-8be4-4ceb-8d7c-ad9905a24386"
}

while True:
    response = requests.request("POST", url, data=payload, headers=headers)
    print(response.text)
    time.sleep(6)
