from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC


# Idk hwo dafuq i run my webdriver on wsl wtf is this HAHHAHAHHAHA
webdriver_path = "/path/to/your/webdriver"


url = "http://127.0.0.1:8000"
driver = webdriver.Chrome(executable_path=webdriver_path)
driver.get(url)

try:
    element = WebDriverWait(driver, 10).until(
        EC.presence_of_element_located((By.ID, "your_element_id_here"))
    )
    print("Element found!")
except:
    print("Element not found or took too long to load.")


driver.quit()