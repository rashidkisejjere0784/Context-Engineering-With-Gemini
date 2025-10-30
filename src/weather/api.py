# src/weather/api.py

import requests

# Define constants for the API endpoints
GEOCODING_API_URL = "https://geocoding-api.open-meteo.com/v1/search"
WEATHER_API_URL = "https://api.open-meteo.com/v1/forecast"

class CityNotFoundError(Exception):
    """Custom exception for when a city cannot be found."""
    pass

def get_coordinates(city_name: str) -> dict:
    """
    Fetches the latitude and longitude for a given city name.

    Args:
        city_name: The name of the city to look up.

    Returns:
        A dictionary containing 'latitude' and 'longitude'.

    Raises:
        CityNotFoundError: If the city cannot be found.
        requests.exceptions.RequestException: For network-related errors.
    """
    params = {"name": city_name, "count": 1}
    response = requests.get(GEOCODING_API_URL, params=params)
    response.raise_for_status()  # Raises an HTTPError for bad responses (4xx or 5xx)

    data = response.json()
    if not data.get("results"):
        raise CityNotFoundError(f"Could not find coordinates for city '{city_name}'.")

    # Return the coordinates from the first result
    result = data["results"][0]
    return {
        "latitude": result["latitude"],
        "longitude": result["longitude"]
    }

def get_weather(lat: float, lon: float) -> dict:
    """
    Fetches the current weather for a given latitude and longitude.

    Args:
        lat: The latitude.
        lon: The longitude.

    Returns:
        A dictionary containing the current weather information.

    Raises:
        requests.exceptions.RequestException: For network-related errors.
    """
    params = {
        "latitude": lat,
        "longitude": lon,
        "current_weather": "true"
    }
    response = requests.get(WEATHER_API_URL, params=params)
    response.raise_for_status()

    return response.json().get("current_weather", {})

