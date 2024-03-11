package hetzner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Hetzner struct {
	APIKey     string
	APIBaseUrl string
}

// A DNS zone in Hetzner's API response
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DNS record in Hetzner's API response
type RecordResponse struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	ZoneID   string `json:"zone_id"`
	TTL      int    `json:"ttl"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

// Structure to update a DNS record
type RecordUpdateRequest struct {
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// API Base URL
func (h *Hetzner) apiBaseURL() string {
	if len(h.APIBaseUrl) > 0 {
		return h.APIBaseUrl
	}
	return "https://dns.hetzner.com/api/v1"
}

// Finds all zones accessible by the current API key
func (h *Hetzner) FindAllZones() ([]Zone, error) {
	url := fmt.Sprintf("%s/zones", h.apiBaseURL())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Zone{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Auth-API-Token", h.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Zone{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return []Zone{}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var zonesResponse struct {
		Zones []Zone `json:"zones"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&zonesResponse); err != nil {
		return []Zone{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return zonesResponse.Zones, nil
}

// Finds the ID of the DNS zone for a given domain name
func (h *Hetzner) FindZoneID(domainName string) (string, error) {
	zones, err := h.FindAllZones()
	if err != nil {
		return "", err
	}

	for _, zone := range zones {
		if zone.Name == domainName {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("zone for domain %s not found", domainName)
}

// Fetches all DNS records for the specified zone ID
func (h *Hetzner) FindAllRecordsForZone(zoneID string) ([]RecordResponse, error) {
	url := fmt.Sprintf("%s/records?zone_id=%s", h.apiBaseURL(), zoneID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Auth-API-Token", h.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var recordsResponse struct {
		Records []RecordResponse `json:"records"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&recordsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return recordsResponse.Records, nil
}

// Finds all records matching a specific name (i.e. _acme-challenge)
func (h *Hetzner) FindRecordsByName(zoneID string, recordName string) ([]RecordResponse, error) {
	records, err := h.FindAllRecordsForZone(zoneID)
	if err != nil {
		return nil, err
	}
	var matchingRecords []RecordResponse
	for _, record := range records {
		if record.Name == recordName {
			matchingRecords = append(matchingRecords, record)
		}
	}
	return matchingRecords, nil
}

// Finds a DNS record by a passed ID
func (h *Hetzner) FindRecordById(zoneID string, recordId string) (RecordResponse, error) {
	records, err := h.FindAllRecordsForZone(zoneID)
	if err != nil {
		return RecordResponse{}, err
	}
	for _, record := range records {
		if record.ID == recordId {
			return record, nil
		}
	}
	return RecordResponse{}, errors.New("record not found")
}

// Prints all the passed records. Used only for debugging.
func (h *Hetzner) PrintRecords(records []RecordResponse) {
	for _, record := range records {
		fmt.Printf("ID: %s, Type: %s, Name: %s, Value: %s, TTL: %d, Created: %s, Modified: %s\n",
			record.ID, record.Type, record.Name, record.Value, record.TTL, record.Created, record.Modified)
	}
}

// Updates an existing DNS record with new information
func (h *Hetzner) UpdateRecord(zoneID, recordID, recordType, recordName, recordValue string) error {
	url := fmt.Sprintf("%s/records/%s", h.apiBaseURL(), recordID)

	updatedRecord := RecordUpdateRequest{
		ID:     recordID,
		ZoneID: zoneID,
		Type:   recordType,
		Name:   recordName,
		Value:  recordValue,
	}

	requestBody, err := json.Marshal(updatedRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth-API-Token", h.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("DNS record updated successfully")
	return nil
}

// Creates a new DNS record/value pair in the specified zone
func (h *Hetzner) CreateRecord(zoneID, recordType, recordName, recordValue string) error {
	url := fmt.Sprintf("%s/records", h.apiBaseURL())

	newRecord := RecordUpdateRequest{
		ZoneID: zoneID,
		Type:   recordType,
		Name:   recordName,
		Value:  recordValue,
	}

	requestBody, err := json.Marshal(newRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth-API-Token", h.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("DNS record created successfully")
	return nil
}

// Checks if a record name/value pair already exists and updates it if it does. Otherwise, this creates a new record with the specified information.
func (h *Hetzner) CreateOrUpdateRecord(zoneID, recordType, recordName, recordValue string) error {
	records, err := h.FindRecordsByName(zoneID, recordName)
	if err != nil {
		return err
	}
	if len(records) > 0 {
		record := records[0]
		return h.UpdateRecord(zoneID, record.ID, record.Type, record.Name, recordValue)
	} else {
		return h.CreateRecord(zoneID, recordType, recordName, recordValue)
	}
}

// Deletes a DNS record given its ID
func (h *Hetzner) DeleteRecord(recordID string) error {
	url := fmt.Sprintf("%s/records/%s", h.apiBaseURL(), recordID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth-API-Token", h.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	fmt.Println("DNS record deleted successfully")
	return nil
}
