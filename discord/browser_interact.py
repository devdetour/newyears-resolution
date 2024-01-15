from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import json
import time
import re
import os
from dotenv import load_dotenv

COOKIE_PATH = "cookies.json"

def save_cookie(driver, path):
    with open(path, 'w') as filehandler:
        json.dump(driver.get_cookies(), filehandler)

def load_cookie(driver, path):
    with open(path, 'r') as cookiesfile:
        cookies = json.load(cookiesfile)
    for cookie in cookies:
        driver.add_cookie(cookie)

def highlight_element(driver, element):
    if element != None:
        # Use JavaScript to change the background color to red
        driver.execute_script("arguments[0].style.backgroundColor = 'red';", element)
        # Wait for a moment to see the highlighted element (you can adjust the sleep time)
        driver.implicitly_wait(2)

def test():
    return "Test from webdriver"

class FirefoxRunner():
    # Start and setup cookies for amazon domain
    def __init__(self):
        self.driver = webdriver.Firefox()
        DEFAULT_URL = "https://www.amazon.com/"
        # url = "https://www.example.com/path/to/page"
        self.driver.get(DEFAULT_URL)
        load_cookie(self.driver, COOKIE_PATH)

        # addProductToCart(driver, link)
        #     self.driver = 
    
    def dispose(self):
        try:
            self.driver.close()
        except:
            pass

    def amazon_checkout(self, url: str):
        try:
            driver = self.driver
            self.driver.get(url)
            time.sleep(2)

            # add to cart
            button = driver.find_element(By.ID, "submit.add-to-cart")
            highlight_element(driver, button)
            button.click()

            # navigate to checkout
            # wait a bit just in case
            atc_selector = (By.ID, "desktop-ptc-button-celWidget")
            WebDriverWait(driver, 3).until(EC.visibility_of_element_located(atc_selector))
            time.sleep(2)

            checkout_widget = driver.find_element(atc_selector[0], atc_selector[1])
            checkout_button = checkout_widget.find_element(By.CLASS_NAME, "a-button-inner")
            highlight_element(driver, checkout_button)
            checkout_button.click()

            # proceed to checkout
            checkout_selector = (By.ID, "submitOrderButtonId")
            WebDriverWait(driver, 3).until(EC.visibility_of_element_located(checkout_selector))
            time.sleep(2)
            purchase_widget = driver.find_element(checkout_selector[0], checkout_selector[1])

            purchase_selector = (By.CLASS_NAME, "a-button-inner")
            WebDriverWait(driver, 3).until(EC.presence_of_element_located(purchase_selector))
            time.sleep(2)
            purchase_button = purchase_widget.find_element(purchase_selector[0], purchase_selector[1])
            highlight_element(driver, purchase_button)
            checkout_button.click() # DANGER DANGER THIS WILL CHECKOUT
            return True
        except Exception as e:
            print("AMAZON CHECKOUT FAILED!")
            print(e)
            return False
