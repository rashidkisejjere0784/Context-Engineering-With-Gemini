# Product Requirements Prompt (PRP) - Weather CLI
## 1. Overview
- **Feature Name:** Weather CLI Tool

- **Objective:** Build a command-line tool that fetches and displays the current weather for a user-specified city.

- **Why:** To provide users with a quick and easy way to check the weather from their terminal without needing a browser.

## 2. Success Criteria
- [x] The code runs without errors.

- [x] All new unit tests pass.

- [x] The CLI tool correctly accepts a city name as an argument.

- [x] The tool successfully fetches data from the Open-Meteo API.

- [x] The tool prints the current temperature in a clear, human-readable format.

- [x] The tool handles errors gracefully if the city is not found or the API fails.

- [x] The code adheres to the project standards defined in `GEMINI.md`.

## 3. Context & Resources
### 📚 External Documentation:
- **Resource:** Open-Meteo Geocoding API
   - **Purpose:** To convert a city name into latitude and longitude. This is the first required step.

- **Resource:** Open-Meteo Weather API
   - **Purpose:** To get the current weather using the coordinates from the Geocoding API. We only need the `current_weather` parameter.

- **Resource:** Python `requests` Library
   - **Purpose:** For making the HTTP GET requests to the APIs.

- **Resource:** Python `argparse` Library
   - **Purpose:** To build the command-line interface.

### 💻 Internal Codebase Patterns:
- **File:** N/A
  - **Reason:** This is the first feature, so we are establishing the patterns now.

### ⚠️ Known Pitfalls:
- The process requires two separate API calls: first to geocode the city, then to get the weather. The second call depends on the success of the first.

- The Geocoding API can return multiple results for a city name. We should always use the first result.

- The final weather API call must include `current_weather=true` in the URL parameters to get the data we need.

## 4. Implementation Blueprint
### Proposed File Structure:
```
   src/
   └── weather/
   ├── __init__.py      (new)
   ├── cli.py           (new, handles argparse)
   └── api.py           (new, handles all Open-Meteo logic)
   tests/
   └── test_weather.py      (new)
```

### Task Breakdown:
**Task 1: API Logic (`src/weather/api.py`)**

- Create a function `get_coordinates(city_name: str) -> dict`.

   - It should call the geocoding API.

   - It should handle the case where no results are found and raise an exception.

   - It should return the first result's `latitude` and `longitude`.

- Create a function `get_weather(lat: float, lon: float) -> dict`.

   - It should call the weather forecast API with the coordinates.

   - It should return the `current_weather` object from the response.

**Task 2: CLI Logic (`src/weather/cli.py`)**

- Create a `main()` function.

- Use `argparse` to define one required argument: city.

- Call the API functions from `api.py` in sequence.

- Implement `try...except` blocks to catch errors from the API module and print user-friendly messages.

- Print the final temperature to the console like: `The current temperature in London is 15.0°C`.

## 5. Validation Plan
### Unit Tests (`tests/test_weather.py`):
- `test_get_coordinates_success()`: Mock the requests.get call and ensure it returns correct coordinates for a valid city.

- `test_get_coordinates_not_found()`: Mock the API returning an empty result and ensure the function raises a custom exception.

- `test_get_weather_success()`: Mock the weather API call and ensure it returns the expected weather data.

- It is not necessary to test the `cli.py` directly with unit tests; we will test it manually.

### Manual Test Command:
```
python -m src.weather.cli --city "Berlin"
```

**Expected Output:**
```
The current temperature in Berlin is 18.5°C.
```
(Note: The exact temperature will vary.)

**Manual Test for Error:**
```
python -m src.weather.cli --city "InvalidCityName123"
```
**Expected Output:**
```
Error: Could not find coordinates for city "InvalidCityName123".
```