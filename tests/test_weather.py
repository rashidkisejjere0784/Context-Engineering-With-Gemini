# tests/test_weather.py

import unittest
from unittest.mock import patch, Mock
from src.weather import api

class TestWeatherApi(unittest.TestCase):

    @patch('requests.get')
    def test_get_coordinates_success(self, mock_get):
        """
        Test that get_coordinates returns correct data on a successful API call.
        """
        mock_response = Mock()
        mock_response.json.return_value = {
            "results": [{
                "latitude": 52.52,
                "longitude": 13.41
            }]
        }
        mock_response.raise_for_status = Mock()
        mock_get.return_value = mock_response

        coords = api.get_coordinates("Berlin")
        self.assertEqual(coords, {"latitude": 52.52, "longitude": 13.41})

    @patch('requests.get')
    def test_get_coordinates_not_found(self, mock_get):
        """
        Test that get_coordinates raises CityNotFoundError for an unknown city.
        """
        mock_response = Mock()
        mock_response.json.return_value = {} # No "results" key
        mock_response.raise_for_status = Mock()
        mock_get.return_value = mock_response

        with self.assertRaises(api.CityNotFoundError):
            api.get_coordinates("InvalidCityName123")

    @patch('requests.get')
    def test_get_weather_success(self, mock_get):
        """
        Test that get_weather returns correct data on a successful API call.
        """
        mock_response = Mock()
        mock_response.json.return_value = {
            "current_weather": {
                "temperature": 18.5,
                "windspeed": 10.0
            }
        }
        mock_response.raise_for_status = Mock()
        mock_get.return_value = mock_response

        weather = api.get_weather(52.52, 13.41)
        self.assertEqual(weather, {"temperature": 18.5, "windspeed": 10.0})

if __name__ == '__main__':
    unittest.main()

