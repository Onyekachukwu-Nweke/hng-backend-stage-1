package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Response struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

type IP2LocationResponse struct {
	City      string  `json:"city_name"`
	Region    string  `json:"region_name"`
	Country   string  `json:"country_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env files")
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/hello", helloHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, router)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	visitorName := r.URL.Query().Get("visitor_name")
	if visitorName == "" {
		http.Error(w, "visitor_name query parameter is required", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)
	location, err := getLocation(clientIP)
	if err != nil {
		http.Error(w, "Error getting location", http.StatusInternalServerError)
		return
	}

	temperature := getTemperature()
	greeting := fmt.Sprintf("Hello, %s! The temperature is %d degrees Celsius in %s", visitorName, temperature, location.City)

	response := Response{
		ClientIP: clientIP,
		Location: location.City,
		Greeting: greeting,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getClientIP(r *http.Request) string {
	real_ip := r.Header.Get("X-Forwarded-For")
	if real_ip != "" {
		// X-Forwarded-For can contain multiple IPs, the first one is the client IP
		ips := strings.Split(real_ip, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if clientIP != "" {
				return clientIP
			}
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "127.0.0.1"
	}
	return userIP.String()
}

func getLocation(ip string) (*IP2LocationResponse, error) {
	apiKey := os.Getenv("IP2LOCATION_API_KEY")
	url := fmt.Sprintf("https://api.ip2location.io/?key=%s&ip=%s&format=json", apiKey, ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var location IP2LocationResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func getTemperature() int {
	// Mocking the temperature for simplicity.
	// In a real application, you would call a weather API here.
	return 11
}
