from flask import Flask, jsonify
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import re
import json
import os
import discord
import random
from dotenv import load_dotenv
import asyncio
import nest_asyncio
from discord.ext import tasks
from threading import Thread
from urllib.parse import urlparse
import requests
from bs4 import BeautifulSoup
from browser_interact import FirefoxRunner

# CONSTANTS AND ENV STUFF
load_dotenv("env.txt")

USER_ID = os.getenv("USER_ID")
BOT_USER_ID = os.getenv("BOT_USER_ID")

IGNORE_USER_IDS = (USER_ID, BOT_USER_ID)

DISCORD_TOKEN = os.getenv("DISCORD_TOKEN")
DISCORD_GUILD_ID = int(os.getenv("DISCORD_GUILD_ID"))

COMMENT_PREFIX = "[COMMENT]"
LINK_PREFIX = "[LINK]"

COOKIE_PATH = "cookies.json"
PRICE_THRESHOLD = 200
CHANNEL_ID = int(os.getenv("DISCORD_CHANNEL_ID"))
GUILD = "DevDetour"

ROLE_NAME="member"

MESSAGE_TEXT = """
DEV DETOUR DID NOT MEET HIS GOAL FOR NEW YEARS RESOLUTION TODAY!!
Here's how you can make him sad. Post a message in this channel with a link to an Amazon product LESS THAN **$200** USD and a witty remark, using the following format:

```
[LINK] <your link here>
[COMMENT] <your witty remark about how I need to exercise more>
```
Don't put ANYTHING else in the message or my bot will probably skip it. Don't put quotes around the link. You MUST include a comment. If your message does not follow
these criteria, the bot will delete it from the channel - this will hopefully keep other messages out of the channel and make it easy to see and vote on links.

Example:
```
[LINK] https://www.amazon.com/ASUS-Whisper-Quiet-Radiators-lifespan-Magnetic-Levitation/dp/B0955WB2KL/?_encoding=UTF8&pd_rd_w=hmzgy&content-id=amzn1.sym.d348963e-8e0d-45d7-b79a-5a442bb25d8b&pf_rd_p=d348963e-8e0d-45d7-b79a-5a442bb25d8b&pf_rd_r=B8AAJVQ07QRZB1CHY1DX&pd_rd_wg=E63Ua&pd_rd_r=184a6365-73c1-4bf3-81c5-408242aabb0f&ref_=pd_gw_deals_ct_t1
[COMMENT] lol u suck
```

**VOTE for your favorite by reacting to the message you like with 👍. ONLY THESE REACTIONS are counted.**

When this message is *12 HOURS OLD*, the bot will buy the thing with the most votes.

"""

# MESSAGE_TEXT = "DevDetour reached his goal! Check back tomorrow... he's bound to get lazy at some point :grin:"
LOOP_DURATION_SECONDS = 60
WAIT_MINUTES_BEFORE_PUNISH = 1

# UTILITY METHODS
def get_host_from_url(url):
    parsed_url = urlparse(url)
    return parsed_url.hostname

def valid_host(url):
    return get_host_from_url(url) in ["www.amazon.com", "amazon.com"]

# FLASK STUFF
def run_flask(client, loop):
    app = Flask(__name__)
    client = client

    @app.route('/start')
    async def hello():
        # client.separate_thread(loop)
        client.activate(loop)
        return 'Hello, World!'

    @app.route('/ad_hoc')
    async def adhoc():
        client.handle_buy_adhoc(loop)
        return "ok"

    @app.route('/goal_met')
    async def goal_met():
        # client.separate_thread(loop)
        # client.activate(loop)
        msgs = ["DevDetour reached his goal! Check back tomorrow... he's bound to get lazy at some point :grin:", "🏋️‍♂️💪 DevDetour just crushed it at the gym! Sweating buckets, pushing limits, and leveling up those fitness goals like a boss! 🔥 #GrindMode #FitnessFiesta 🚀", "DevDetour conquered today's workout with sheer determination and grit, leaving no room for excuses. The grind is real, and the results are undeniable.", "Oh, bravo to DevDetour for the Herculean feat of walking for a whole half-hour! I mean, who knew putting one foot in front of the other could be such an Olympic-level challenge? 🚶‍♂️😏", "Oh, the sheer awe-inspiring spectacle! DevDetour conquered the monumental task of a light 30-minute workout. I mean, who needs to break a sweat or challenge those muscles anyway?"]
        msg = random.choice(msgs) + " You can follow along with goal progress at https://detour.dev !"
        client.send_message_sync(msg, loop)
        return msg

    @app.route('/inspect')
    async def inspect():
        return f"STATE: {client.active}, loops left: {client.loops_left}"

    @app.route('/delete_all')
    async def delete():
        client.delete_all_messages(loop)
        return f"Clearing messages!"

    app.run(host='0.0.0.0', port=5000)



# UTILITY method to get price from Amazon.
def get_price(url):
    response = requests.get(url)
    if response.status_code == 200:
        # Parse the HTML content with BeautifulSoup
        soup = BeautifulSoup(response.text, 'html.parser')

        price_span = soup.find('span', class_='a-offscreen')
        return price_span.text

    else:
        print(f"Failed to retrieve the page. Status code: {response.status_code}")
        response = requests.get(url)

    return "$1000"

# Check that price is under threshold
def price_under_threshold(price_string):
    try:
        # Extract numeric part from the price string using regex
        numeric_part = re.sub(r'[^\d.]', '', price_string)

        # Convert the numeric part to a float
        price = float(numeric_part)

        # Check if the price is less than the threshold value
        return price < PRICE_THRESHOLD

    except ValueError:
        # Handle the case where the conversion fails
        print(f"Invalid price string: {price_string}")
        return False

class CustomClient(discord.Client):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        # an attribute we can access from our task
        self.counter = 0
        
        # keep track of if we should be waiting to buy something, and how many loops we have left
        self.loops_left = 0
        self.active = False # these can be just 1 var but whatever

    async def setup_hook(self):
        # self.my_background_task.start()
        self.evaluate_time_to_buy.start()

    def change_permissions_sync(self, loop, role_name):
        print("Changing role permissions sync")
        asyncio.run_coroutine_threadsafe(self.change_permissions(role_name, True), loop)
        
    def activate(self, loop):
        if self.active:
            print("Already active! Doing nothing!")
            return
        else:
            # number of loops to wait is WAIT_MINUTES_BEFORE_PUNISH / (LOOP_DURATION_SECONDS / 60)
            self.active = True
            self.loops_left = int(WAIT_MINUTES_BEFORE_PUNISH / (LOOP_DURATION_SECONDS / 60))
            self.change_permissions_sync(loop, ROLE_NAME)
            # self.SendNotificationMessage(loop)

    # Send a message to the channel with instructions for how to make me sad.
    def SendNotificationMessage(self, loop):
        channel = self.get_channel(CHANNEL_ID)
        message = MESSAGE_TEXT
        asyncio.run_coroutine_threadsafe(channel.send(message), loop)

    # Send message to regular channel
    async def send_message(self, msg):
        channel = self.get_channel(CHANNEL_ID)
        await channel.send(msg)

    def send_message_sync(self, msg, loop):
        asyncio.run_coroutine_threadsafe(self.send_message(msg), loop)
    
    def delete_all_messages(self, loop):
        asyncio.run_coroutine_threadsafe(self.clear_channel(), loop)

    @tasks.loop(seconds=LOOP_DURATION_SECONDS)  # task runs every 60 seconds
    async def my_background_task(self):
        channel = self.get_channel(CHANNEL_ID)  # channel ID goes here
        self.counter += 1
        await channel.send(self.counter)

    @tasks.loop(seconds=LOOP_DURATION_SECONDS)
    async def evaluate_time_to_buy(self):
        print("EVALUATING TIME TO BUY")
        if not self.active:
            print("NOT ACTIVE! Not buying.")
            return
        elif self.loops_left > 0:
            self.loops_left -= 1
            if int(self.loops_left) == 1:
                channel = self.get_channel(CHANNEL_ID)  # channel ID goes here
                await channel.send(f"BUYING IN {LOOP_DURATION_SECONDS} seconds, LAST CHANCE!")
            print(f"Decremented loops left! Now: {self.loops_left}")

        if int(self.loops_left) == 0: # time to buy
            print("TIME TO BUY!!!!")
            await self.handle_buy()
            # self.active = False
            # chosenMessage = await self.get_messages_by_votes()

    def execute_buy_in_background(self, valid, loop, link):
        asyncio.run_coroutine_threadsafe(self.execute_buy(valid, link), loop)

    # TODO undo this all
    async def execute_buy(self, valid, link):
        print("RUNNING BUY")
        runner = FirefoxRunner() # TODO cleanup runner between runs?
        runner.amazon_checkout(link) # TODO validate
        await self.send_message("BOUGHT " + link)
        return

    def handle_buy_adhoc(self, loop):
        # loop = asyncio.get_running_loop() 
        valid = True
        t = Thread(target=self.execute_buy_in_background, args=[valid, loop])
        t.start()

    async def handle_buy(self):
        self.active = False
        await self.change_permissions(ROLE_NAME, False)
        print("Turned post permissions off!")

        messages_in_order = await self.get_messages_by_votes()

        valid, link, comment = False, "", ""
        # Go until we find a message that is valid
        for i in range(len(messages_in_order)):
            curr = messages_in_order[i]
            valid, link, comment = self.try_parse_message_content(curr.content)
            print(f"Valid: {valid}, link: {link}, comment: {comment}")
            if not valid:
                continue
            price = get_price(link)
            print(f"GOT PRICE: {price}")
            valid = price_under_threshold(price)
            if valid:
                message = f"CHOSE PRODUCT: {link} with price: {price}"
                break
            else:
                message = f"Product {link} too expensive ({price} vs. ${PRICE_THRESHOLD})"
                await self.send_message(message)
        
        # loooooop
        loop = asyncio.get_running_loop() 
        t = Thread(target=self.execute_buy_in_background, args=[valid, loop, link])
        t.start()
        print("Done with loop in main thread")

    # Try to parse the contents of a message
    def try_parse_message_content(self, msg: str):
        try:
            lines = msg.split("\n")
            product_link = lines[0].split(LINK_PREFIX)[1].strip()
            comment = lines[1].split(COMMENT_PREFIX)[1].strip()
            print(product_link)
            print(comment)

            if not valid_host(product_link):
                return (False, "", "")
            return (True, product_link, comment)
        except Exception as e:
            print("Failed to parse:", e)
            return (False, "", "")


    # delete all messages in channel
    async def clear_channel(self):
        messages = await self.recent_messages()
        for m in messages:
            await m.delete()

    async def get_messages_by_votes(self):
        messages = await self.recent_messages()

        window = []
        # get messages SINCE the bot's LAST message....
        # loop thru messages (in reverse chronological order. most recently sent are first.)
        # when we hit the bot's PREVIOUS start message, stop parsing.
        for m in messages:
            if m.author == self.user and m.content == MESSAGE_TEXT:
                print("BREAKING!")
                break
            window.append(m)

        print(f"Got {len(window)} messages since last evaluation")

        # Filter out bot's messages
        filtered = [m for m in window if m.author != client.user]
        print(f"Messages: {len(messages)}, Filtered: {len(filtered)}")

        sortedMessages = self.get_message_by_reactions(filtered)
        print("CHOSEN MESSAGES LENGTH: ", len(sortedMessages))

        for m in sortedMessages:
            print(f"Message content: {m.content}, {m.reactions}, {m.author}")
        
        return sortedMessages
        
            
    def get_message_by_reactions(self, messages):
        if messages is None or len(messages) == 0:
            return []

        messages = sorted(messages, reverse=True, key=lambda m: sum([r.count for r in m.reactions if r.emoji == "👍" ]))
        return messages

    def separate_thread(self, loop):
        print("running separate thread")
        channel = self.get_channel(CHANNEL_ID)  # channel ID goes here
        self.counter += 1
        asyncio.run_coroutine_threadsafe(channel.send("TEST"), loop)

    @my_background_task.before_loop
    async def before_my_task(self):
        await self.wait_until_ready()  # wait until the bot logs in

    @evaluate_time_to_buy.before_loop
    async def before_evaluate_time_to_buy(self):
        await self.wait_until_ready()  # wait until the bot logs in

    async def handle_message(self, message):
        print("HANDLING MESSAGE")
        print(message)
        if message.author.id == BOT_USER_ID or (message.author.id == USER_ID and "OVERRIDE" in message.content): # ignore
            return
        if message.channel.id == CHANNEL_ID:
            print("IT IS THE CHANNEL WE CARE ABOUT!")
            valid, link, comment = self.try_parse_message_content(message.content)
            if(valid):
                print("MESSAGE FORMAT OK!")
            else:
                print("MESSAGE FORMAT WRONG!")
                await message.delete()

    # @CustomClient.event
    async def on_message(self, message):
        await self.handle_message(message)

    # @client.event
    async def on_ready(client):
        print(f'{client.user} has connected to Discord!')
        for guild in client.guilds:
            if guild.name == GUILD:
                break
        
        print(
            f'{client.user} is connected to the following guild:\n'
            f'{guild.name}(id: {guild.id})'
        )
        
        text_channels = guild.text_channels

        # start flask server HERE, once we have loop
        loop = asyncio.get_running_loop() 
        # now THIS right here is unbelievably cursed
        t = Thread(target=run_flask, args=[client, loop])
        t.start()

        # Print the names of text channels
        channel_names = [(channel.name, channel.id) for channel in text_channels]
        print("Text channels:", channel_names)

    async def print_recent_messages(self):
        # Filter out bot's messages
        messages = await self.recent_messages()
        filtered = [m for m in messages if m.author != client.user]
        print(f"Messages: {len(messages)}, Filtered: {len(filtered)}")

        for m in filtered:
            print(m.content, m.reactions)

    async def recent_messages(self, limit: int = 500):
        channel = self.get_channel(CHANNEL_ID)

        # Get the most recent messages (limit parameter specifies the number of messages)
        messages = [ i async for i in channel.history(limit=limit)] #.flatten()
        print(f"Got {len(messages)} Recent messages: ")

        # Print the content of the most recent messages
        # for message in messages:
        #     print(f'Message content: {message.content}')

        return messages

    # @client.event
    async def change_permissions(self, role_name: str, allow=True):
        print("Changing permissions!")
        channel = self.get_channel(CHANNEL_ID)
        guild = self.get_guild(DISCORD_GUILD_ID)
        role = discord.utils.get(guild.roles, name=role_name)

        print("Channel", channel)
        print("guild", guild)
        print("role", role)

        if role:
            # Get the channel overwrite for the specific role
            overwrite = channel.overwrites_for(role)

            # Modify permissions as needed
            overwrite.update(send_messages=allow)

            # Apply the changes to the channel
            await channel.set_permissions(role, overwrite=overwrite)

            print(f'Permissions for role {role_name} changed, allow={allow}.')
        else:
            print(f'Role {role_name} not found.')


#     return True


if __name__ == '__main__':
    # Setup Discord intents
    intents = discord.Intents.default()
    intents.messages = True  # Enable message-related events
    intents.message_content = True

    client = CustomClient(intents=intents)

    # somehow need to get LOOP from CLIENT to FLASK thread.
    client.run(DISCORD_TOKEN)