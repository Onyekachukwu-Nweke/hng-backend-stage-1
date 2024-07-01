# hng-backend-stage-1

This project is a basic web server written in Go. It provides an API endpoint that returns the client's IP address, location, and a greeting message with the current temperature in the client's city. The location is determined using the IP2Location API, and the temperature is fetched from the OpenWeather API. The server is deployed on DigitalOcean App Platform.

## Features

- API endpoint to greet the visitor and provide location-based temperature.
- Uses IP2Location API for geolocation.
- Uses OpenWeather API for fetching temperature data.
- Dockerized for easy deployment.

## API Endpoint

### [GET] `/api/hello?visitor_name="Owerri"`

#### Response:
```json
{
  "client_ip": "127.0.0.1",
  "location": "Owerri",
  "greeting": "Hello, Onyekachukwu! The temperature is 24.85 degrees Celsius in Owerri"
}

