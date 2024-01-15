# run once to setup cookies. login, then give input to save login cookies.
from browser_interact import FirefoxRunner, save_cookie, COOKIE_PATH

runner = FirefoxRunner()

input()
save_cookie(runner.driver, COOKIE_PATH)