# ulysses
This is my New Year's resolution bot. The Golang + React part is the main business logic related to creating and evaluating goals.

I started designing it to be extensible, then quickly started cutting corners to get it done by Jan 1. should I have used Airflow? YES! but instead I built this.

The discord bot part is in the `discord/` directory, it's quite janky but mostly works (this i did in literally 1 day on Dec 30 so it's especially rough lmao)

The backend stuff (auth especially) was copy/pasted/edited from gofiber's example: https://github.com/gofiber/recipes/tree/master/auth-jwt.

I have ripped out any secrets & tokens and stuff out of the code as much as possible - I also rotated the actually dangerous ones just in case I missed any spots, so
no use looking for credentials here.

GLHF using!