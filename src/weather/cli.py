# src/weather/cli.py

import argparse
from . import api

def main():
    """
    The main entry point for the Weather CLI tool.
    Parses arguments, fetches weather, and prints the result.
    """
    parser = argparse.ArgumentParser(
        description="Get the current weather for a specified city."
    )
    parser.add_argument(
        "--city",
        type=str,
        required=True,
        help="The name of the city to get the weather for."
    )
    args = parser.parse_args()
    city = args.city

    try:
        # 1. Get coordinates for the city
        print(f"Fetching coordinates for {city}...")
        coords = api.get_coordinates(city)
        lat, lon = coords["latitude"], coords["longitude"]

        # 2. Get the weather for the coordinates
        print(f"Fetching weather for coordinates ({lat}, {lon})...")
        weather = api.get_weather(lat, lon)

        if not weather:
            print(f"Could not retrieve current weather for {city}.")
            return

        # 3. Print the result
        temp = weather.get("temperature")
        unit = "°C" # Assuming Celsius as default from API
        print("-" * 30)
        print(f"The current temperature in {city} is {temp}{unit}.")
        print("-" * 30)

    except api.CityNotFoundError as e:
        print(f"Error: {e}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

if __name__ == "__main__":
    main()
