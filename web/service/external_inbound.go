// Package service provides business logic services for the 3x-ui web panel,
// including inbound/outbound management, user administration, settings, and Xray integration.
package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
)

type externalInboundSettings struct {
	Clients []model.Client `json:"clients"`
}

// ExternalInboundService handles CRUD operations for external inbounds.
type ExternalInboundService struct{}

// GetExternalInbounds retrieves all external inbounds for a user.
func (s *ExternalInboundService) GetExternalInbounds(userId int) ([]*model.ExternalInbound, error) {
	db := database.GetDB()
	var externalInbounds []*model.ExternalInbound
	err := db.Model(model.ExternalInbound{}).Where("user_id = ?", userId).Find(&externalInbounds).Error
	return externalInbounds, err
}

// GetExternalInbound retrieves a specific external inbound by ID.
func (s *ExternalInboundService) GetExternalInbound(id int) (*model.ExternalInbound, error) {
	db := database.GetDB()
	externalInbound := &model.ExternalInbound{}
	err := db.Model(model.ExternalInbound{}).Where("id = ?", id).First(externalInbound).Error
	return externalInbound, err
}

// AddExternalInbound creates a new external inbound.
func (s *ExternalInboundService) AddExternalInbound(externalInbound *model.ExternalInbound) error {
	if strings.TrimSpace(externalInbound.Name) == "" {
		return fmt.Errorf("external inbound name is required")
	}
	if strings.TrimSpace(externalInbound.Url) == "" {
		return fmt.Errorf("external inbound url is required")
	}
	db := database.GetDB()
	return db.Create(externalInbound).Error
}

// UpdateExternalInbound updates an existing external inbound.
func (s *ExternalInboundService) UpdateExternalInbound(externalInbound *model.ExternalInbound) error {
	if strings.TrimSpace(externalInbound.Name) == "" {
		return fmt.Errorf("external inbound name is required")
	}
	if strings.TrimSpace(externalInbound.Url) == "" {
		return fmt.Errorf("external inbound url is required")
	}
	db := database.GetDB()
	return db.Save(externalInbound).Error
}

// DeleteExternalInbound deletes an external inbound and its clients.
func (s *ExternalInboundService) DeleteExternalInbound(id int) error {
	db := database.GetDB()
	if err := db.Where("external_inbound_id = ?", id).Delete(&model.ExternalInboundClient{}).Error; err != nil {
		return err
	}
	return db.Delete(&model.ExternalInbound{}, id).Error
}

// GetClients retrieves external inbound clients by inbound ID.
func (s *ExternalInboundService) GetClients(externalInboundId int) ([]model.ExternalInboundClient, error) {
	db := database.GetDB()
	var clients []model.ExternalInboundClient
	err := db.Model(model.ExternalInboundClient{}).
		Where("external_inbound_id = ?", externalInboundId).
		Find(&clients).Error
	return clients, err
}

// AddExternalInboundClient appends clients to an external inbound.
func (s *ExternalInboundService) AddExternalInboundClient(id int, settings string) error {
	clients, err := parseExternalInboundClients(settings)
	if err != nil {
		return err
	}
	if len(clients) == 0 {
		return fmt.Errorf("missing external client")
	}
	client := mapExternalClient(clients[0])
	client.ExternalInboundId = id
	client.CreatedAt = nowMillisIfZero(client.CreatedAt)
	client.UpdatedAt = nowMillisIfZero(client.UpdatedAt)
	db := database.GetDB()
	return db.Create(&client).Error
}

// UpdateExternalInboundClient replaces a client in an external inbound.
func (s *ExternalInboundService) UpdateExternalInboundClient(id int, clientId string, settings string) error {
	clients, err := parseExternalInboundClients(settings)
	if err != nil {
		return err
	}
	if len(clients) == 0 {
		return fmt.Errorf("missing external client")
	}
	updated := mapExternalClient(clients[0])
	db := database.GetDB()
	var existing model.ExternalInboundClient
	err = db.Model(model.ExternalInboundClient{}).
		Where("external_inbound_id = ?", id).
		Where("client_id = ? OR email = ?", clientId, clientId).
		First(&existing).Error
	if err != nil {
		return err
	}
	updated.ExternalInboundId = id
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedAt = nowMillisIfZero(updated.UpdatedAt)
	updateData := map[string]any{
		"enable":              updated.Enable,
		"client_id":           updated.ClientID,
		"email":               updated.Email,
		"sub_id":              updated.SubID,
		"tg_id":               updated.TgID,
		"comment":             updated.Comment,
		"expiry_time":         updated.ExpiryTime,
		"reset":               updated.Reset,
		"updated_at":          updated.UpdatedAt,
		"external_inbound_id": updated.ExternalInboundId,
	}
	return db.Model(&existing).Updates(updateData).Error
}

// DeleteExternalInboundClient removes a client from an external inbound.
func (s *ExternalInboundService) DeleteExternalInboundClient(id int, clientId string) error {
	db := database.GetDB()
	return db.Where("external_inbound_id = ?", id).
		Where("client_id = ? OR email = ?", clientId, clientId).
		Delete(&model.ExternalInboundClient{}).Error
}

func parseExternalInboundClients(settings string) ([]model.Client, error) {
	payload := externalInboundSettings{}
	if strings.TrimSpace(settings) == "" {
		return []model.Client{}, nil
	}
	if err := json.Unmarshal([]byte(settings), &payload); err != nil {
		return nil, err
	}
	return payload.Clients, nil
}

func mapExternalClient(client model.Client) model.ExternalInboundClient {
	return model.ExternalInboundClient{
		Enable:     client.Enable,
		ClientID:   client.ID,
		Email:      client.Email,
		SubID:      client.SubID,
		TgID:       client.TgID,
		Comment:    client.Comment,
		ExpiryTime: client.ExpiryTime,
		Reset:      client.Reset,
		CreatedAt:  client.CreatedAt,
		UpdatedAt:  client.UpdatedAt,
	}
}

func nowMillisIfZero(value int64) int64 {
	if value == 0 {
		return time.Now().Unix() * 1000
	}
	return value
}
