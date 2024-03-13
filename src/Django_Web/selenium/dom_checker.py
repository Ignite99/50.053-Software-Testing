import chromedriver_autoinstaller
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

url = "http://127.0.0.1:8000"

chromedriver_autoinstaller.install()
driver = webdriver.Chrome()
print("Started driver")
driver.get(url)
elem = driver.find_element(By.CLASS_NAME, "mb-4")

if elem is not None:
    print("Found element with class 'mb-4' found!")
else: 
    print("Element with class 'mb-4' not found")

driver.close()