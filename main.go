package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Notification struct {
	NotifID          int    `json:"notifId"`
	ECValidationDate string `json:"ecValidationDate"`
	Reference        string `json:"reference"`
	NotifyingCountry struct {
		OrganizationName string `json:"organizationName"`
		ISOCode          string `json:"isoCode"`
	} `json:"notifyingCountry"`
	Subject         string `json:"subject"`
	ProductCategory struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"productCategory"`
	ProductType struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"productType"`
	NotificationClassification struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"notificationClassification"`
	RiskDecision struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"riskDecision"`
	Published       bool `json:"published"`
	OriginCountries []struct {
		OrganizationName string `json:"organizationName"`
		ISOCode          string `json:"isoCode"`
	} `json:"originCountries"`
}

type Response struct {
	Notifications []Notification `json:"notifications"`
	TotalPages    int            `json:"totalPages"`
	TotalElements int            `json:"totalElements"`
}

func main() {
	// Set up the request body
	requestBody := map[string]interface{}{
		"parameters": map[string]interface{}{
			"pageNumber":   1,
			"itemsPerPage": 25,
		},
		"notificationReference":      nil,
		"subject":                    nil,
		"ecValidDateFrom":            time.Now().AddDate(0, 0, -1).Format("02-01-2006 15:04:05"),
		"ecValidDateTo":              time.Now().AddDate(0, 0, -1).Format("02-01-2006 15:04:05"),
		"notifyingCountry":           []int{2},
		"originCountry":              nil,
		"distributionCountry":        nil,
		"notificationType":           nil,
		"notificationClassification": nil,
		"notificationBasis":          nil,
		"productCategory":            nil,
		"actionTaken":                nil,
		"hazardCategory":             nil,
		"riskDecision":               nil,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return
	}

	// Create a request object
	req, err := http.NewRequest("POST", "https://webgate.ec.europa.eu/rasff-window/backend/public/notification/search/consolidated/", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse the response body into a Response struct
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing response body:", err)
		return
	}

	// Print the parsed response
	fmt.Println("Notifications:")
	for _, notification := range response.Notifications {
		fmt.Printf("Notification ID: %d\n", notification.NotifID)
		fmt.Printf("EC Validation Date: %s\n", notification.ECValidationDate)
		fmt.Printf("Reference: %s\n", notification.Reference)
		fmt.Printf("Notifying Country: %s (%s)\n", notification.NotifyingCountry.OrganizationName, notification.NotifyingCountry.ISOCode)
		fmt.Printf("Subject: %s\n", notification.Subject)
		fmt.Printf("Product Category: %s\n", notification.ProductCategory.Description)
		fmt.Printf("Product Type: %s\n", notification.ProductType.Description)
		fmt.Printf("Notification Classification: %s\n", notification.NotificationClassification.Description)
		fmt.Printf("Risk Decision: %s\n", notification.RiskDecision.Description)
		fmt.Printf("Published: %t\n", notification.Published)
		fmt.Printf("Origin Countries:\n")
		for _, originCountry := range notification.OriginCountries {
			fmt.Printf("  - %s (%s)\n", originCountry.OrganizationName, originCountry.ISOCode)
		}
		fmt.Println()
	}

	fmt.Printf("Total Pages: %d\n", response.TotalPages)
	fmt.Printf("Total Elements: %d\n", response.TotalElements)
}
