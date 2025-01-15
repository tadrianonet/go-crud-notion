package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-crud-notion/internal/config"
	"go-crud-notion/internal/entities"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type UserServiceImpl struct {
	notionAPIKey     string
	notionDatabaseID string
	httpClient       *http.Client
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		notionAPIKey:     config.NotionToken,
		notionDatabaseID: config.DatabaseID,
		httpClient:       &http.Client{},
	}
}

func (r *UserServiceImpl) Save(user *entities.User) error {
	user.ID = uuid.New().String()

	payload, err := r.buildPayload(user)
	if err != nil {
		return fmt.Errorf("failed to build payload: %w", err)
	}

	req, err := r.buildRequest("POST", config.NotionAPI+"/pages", payload)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return r.handleResponse(resp)
}

func (r *UserServiceImpl) FindByID(id string) (*entities.User, error) {
	payload := map[string]interface{}{
		"filter": map[string]interface{}{
			"property": "ID",
			"rich_text": map[string]string{
				"equals": id,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request payload: %w", err)
	}

	req, err := r.buildRequest("POST", config.NotionAPI+"/databases/"+r.notionDatabaseID+"/query", payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch user: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	results, ok := result["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}

	return r.mapResponseToUser(results[0].(map[string]interface{})), nil
}

func (r *UserServiceImpl) UpdateUser(user *entities.User) error {
	if user.ID == "" {
		return errors.New("user ID cannot be empty")
	}

	payload := map[string]interface{}{
		"properties": map[string]interface{}{
			"Name": map[string]interface{}{
				"title": []map[string]interface{}{
					{"text": map[string]string{"content": user.Name}},
				},
			},
			"Email": map[string]interface{}{
				"rich_text": []map[string]interface{}{
					{"text": map[string]string{"content": user.Email}},
				},
			},
			"Telefone": map[string]interface{}{
				"rich_text": []map[string]interface{}{
					{"text": map[string]string{"content": user.Telefone}},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create request payload: %w", err)
	}

	req, err := r.buildRequest("PATCH", config.NotionAPI+"/pages/"+user.ID, payloadBytes)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return r.handleResponse(resp)
}

func (r *UserServiceImpl) FindAll() ([]*entities.User, error) {
	req, err := r.buildRequest("POST", config.NotionAPI+"/databases/"+r.notionDatabaseID+"/query", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch users: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return r.mapResponseToUsers(result), nil
}

func (r *UserServiceImpl) buildPayload(user *entities.User) ([]byte, error) {
	payload := map[string]interface{}{
		"parent": map[string]string{
			"database_id": r.notionDatabaseID,
		},
		"properties": map[string]interface{}{
			"ID": map[string]interface{}{
				"rich_text": []map[string]interface{}{
					{"text": map[string]string{"content": user.ID}},
				},
			},
			"Name": map[string]interface{}{
				"title": []map[string]interface{}{
					{"text": map[string]string{"content": user.Name}},
				},
			},
			"Email": map[string]interface{}{
				"rich_text": []map[string]interface{}{
					{"text": map[string]string{"content": user.Email}},
				},
			},
			"Telefone": map[string]interface{}{
				"rich_text": []map[string]interface{}{
					{"text": map[string]string{"content": user.Telefone}},
				},
			},
		},
	}
	return json.Marshal(payload)
}

func (r *UserServiceImpl) buildRequest(method, url string, payload []byte) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		body = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.notionAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", config.NotionVersion)

	return req, nil
}

func (r *UserServiceImpl) handleResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP request failed: %s", string(body))
	}
	return nil
}

func (r *UserServiceImpl) mapResponseToUser(result map[string]interface{}) *entities.User {
	properties := result["properties"].(map[string]interface{})
	return &entities.User{
		ID:       r.getRichTextContent(properties["ID"]),
		Name:     r.getTitleContent(properties["Name"]),
		Email:    r.getRichTextContent(properties["Email"]),
		Telefone: r.getRichTextContent(properties["Telefone"]),
		PageId:   result["id"].(string),
	}
}

func (r *UserServiceImpl) mapResponseToUsers(result map[string]interface{}) []*entities.User {
	results := result["results"].([]interface{})
	users := make([]*entities.User, 0, len(results))

	for _, item := range results {
		user := r.mapResponseToUser(item.(map[string]interface{}))
		users = append(users, user)
	}

	return users
}

func (r *UserServiceImpl) DeleteUserByPageID(pageID string) error {
	if pageID == "" {
		return errors.New("page ID is required")
	}

	payload := map[string]bool{
		"archived": true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create request payload: %w", err)
	}

	req, err := r.buildRequest("PATCH", config.NotionAPI+"/pages/"+pageID, payloadBytes)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete user: %s", string(body))
	}

	return nil
}

func (r *UserServiceImpl) getRichTextContent(field interface{}) string {
	richText := field.(map[string]interface{})["rich_text"].([]interface{})
	if len(richText) == 0 {
		return ""
	}
	return richText[0].(map[string]interface{})["text"].(map[string]interface{})["content"].(string)
}

func (r *UserServiceImpl) getTitleContent(field interface{}) string {
	title := field.(map[string]interface{})["title"].([]interface{})
	if len(title) == 0 {
		return ""
	}
	return title[0].(map[string]interface{})["text"].(map[string]interface{})["content"].(string)
}
