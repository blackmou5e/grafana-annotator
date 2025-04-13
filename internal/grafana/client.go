package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/blackmou5e/grafana-annotator/pkg/errors"
)

type GrafanaClient interface {
	FetchDashboards(ctx context.Context) ([]Dashboard, error)
	CreateAnnotation(ctx context.Context, annotation Annotation) error
}
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

var _ GrafanaClient = (*Client)(nil)

type Dashboard struct {
	ID          int      `json:"id"`
	UID         string   `json:"uid"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	IsStarred   bool     `json:"isStarred"`
	FolderID    int      `json:"folderId"`
	FolderUID   string   `json:"folderUid"`
	FolderTitle string   `json:"folderTitle"`
	FolderURL   string   `json:"folderUrl"`
}

type Annotation struct {
	DashboardUID string   `json:"dashboardUID,omitempty"`
	PanelID      int      `json:"panelId,omitempty"`
	TimeStart    int64    `json:"time"`
	TimeEnd      int64    `json:"timeEnd,omitempty"`
	Tags         []string `json:"tags"`
	Text         string   `json:"text"`
}

func NewClient(baseURL, token string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) FetchDashboards(ctx context.Context) ([]Dashboard, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/search?type=dash-db", c.baseURL), nil)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrGrafanaAPI, "Failed to create request", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrGrafanaAPI, "Failed to fetch dashboards", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewAppError(errors.ErrGrafanaAPI,
			fmt.Sprintf("Unexpected status code: %d", resp.StatusCode), nil)
	}

	var dashboards []Dashboard
	if err := json.NewDecoder(resp.Body).Decode(&dashboards); err != nil {
		return nil, errors.NewAppError(errors.ErrGrafanaAPI, "Failed to decode response", err)
	}

	return dashboards, nil
}

func (c *Client) CreateAnnotation(ctx context.Context, annotation Annotation) error {
	jsonData, err := json.Marshal(annotation)
	if err != nil {
		return errors.NewAppError(errors.ErrInternal, "Failed to marshal annotation", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/annotations", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.NewAppError(errors.ErrGrafanaAPI, "Failed to create request", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.NewAppError(errors.ErrGrafanaAPI, "Failed to create annotation", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.NewAppError(errors.ErrGrafanaAPI,
			fmt.Sprintf("Failed to create annotation. Status: %d, Body: %s",
				resp.StatusCode, string(body)), nil)
	}

	return nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")
}
