import chromedriver_autoinstaller
import time

from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.support import expected_conditions as EC

URL = "http://127.0.0.1:8000"
DYNAMIC_TABLE_URL = "http://127.0.0.1:8000/%2Fdatatb/product/add"

# ========================= [DEBUGGING FUNCTION] =========================
def CheckElementExistInDOM(url: str, selector_val: str, selector_type) -> bool:
    chromedriver_autoinstaller.install()

    driver = webdriver.Chrome()
    driver.get(url)
    elem = driver.find_element(selector_type, selector_val)
    driver.close()

    if elem is not None:
        print(f"Found element with class '{selector_val}'!")
        return True
    else: 
        print(f"Element with class '{selector_val}' not found")
        return False

# checks for element in a DOM within an iframe element
# this function will try to find the latest updated table entry
def FindLatestTableEntry(url: str) -> None:
    chromedriver_autoinstaller.install()

    driver = webdriver.Chrome()
    driver.get(url)
    time.sleep(2) # wait for iframe to load

    # switch driver to iframe
    driver.switch_to.frame(driver.find_element(By.TAG_NAME, "iframe"))

    # find next button to click to latest page
    pagination_list = driver.find_element(By.CLASS_NAME, "pagination.justify-content-center")
    pagination_links = pagination_list.find_elements(By.TAG_NAME, "li")
    last_link = pagination_links[-1]
    next_link = last_link.find_element(By.TAG_NAME, "a")
    actions = ActionChains(driver)
    actions.move_to_element(next_link).click().perform()
    time.sleep(1) # wait for page to switch

    # find latest table entry
    tbody = driver.find_element(By.TAG_NAME, "tbody")
    rows = tbody.find_elements(By.TAG_NAME, "tr")
    cells = rows[-1].find_elements(By.TAG_NAME, "td")
    print(f"Latest entry: [{cells[0].text},{cells[1].text},{cells[2].text},{cells[3].text}]")
    with open("selenium_output.txt", "w") as file:
        for i in range(4):
            file.write(cells[i].text + ",")
    driver.close()

if __name__ == '__main__':
    FindLatestTableEntry(DYNAMIC_TABLE_URL)
