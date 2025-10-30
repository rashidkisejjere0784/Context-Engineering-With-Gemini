import unittest
from unittest.mock import patch
from src.weather_cli import weather
from src.weather_cli import cli
from src.weather_cli.weather import GeocodingResult, WeatherData, CurrentWeather
import requests
import argparse

class TestWeatherCLI(unittest.TestCase):

    def test_get_coordinates_success(self):
        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = {
                "results": [{"latitude": 51.5, "longitude": 0.1}]
            }
            result = weather.get_coordinates("London")
            self.assertEqual(result.latitude, 51.5)
            self.assertEqual(result.longitude, 0.1)

    def test_get_coordinates_city_not_found(self):
        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = {"results": None}
            with self.assertRaises(ValueError):
                weather.get_coordinates("NonExistentCity")

    def test_get_weather_success(self):
        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = {
                "current_weather": {"temperature": 15.5}
            }
            result = weather.get_weather(51.5, 0.1)
            self.assertEqual(result.current_weather.temperature, 15.5)

    def test_get_weather_api_error(self):
         with patch('requests.get') as mock_get:
            mock_get.side_effect = requests.exceptions.RequestException("API Error")
            with self.assertRaises(Exception) as context:
                weather.get_weather(51.5, 0.1)
            self.assertTrue("API request failed" in str(context.exception))

    @patch('src.weather_cli.cli.get_coordinates')
    @patch('src.weather_cli.cli.get_weather')
    @patch('argparse.ArgumentParser.parse_args', return_value=argparse.Namespace(city="London"))
    @patch('sys.stdout.write')
    def test_cli_valid_city(self, mock_stdout, mock_argparse, mock_get_weather, mock_get_coordinates):
        mock_get_coordinates.return_value = GeocodingResult(latitude=51.5, longitude=0.1)
        mock_get_weather.return_value = WeatherData(current_weather=CurrentWeather(temperature=20.0))

        cli.main()
        mock_stdout.assert_called_with("The current temperature in London is 20.0°C")

    @patch('src.weather_cli.cli.get_coordinates', side_effect=ValueError("City not found"))
    @patch('argparse.ArgumentParser.parse_args', return_value=argparse.Namespace(city="NonExistentCity"))
    @patch('sys.stdout.write')
    def test_cli_invalid_city(self, mock_stdout, mock_argparse, mock_get_coordinates):
        cli.main()
        mock_stdout.assert_called_with("Error: City not found")

if __name__ == '__main__':
    unittest.main()
