package validator

import (
	"testing"

	"github.com/FatWang1/open_ticket_system/internal/models"
)

func TestValidateCreateTicketRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *models.CreateTicketRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &models.CreateTicketRequest{
				Name:       "Test Ticket",
				Creator:    "testuser",
				TemplateID: 1,
				Memo:       "Test memo",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			request: &models.CreateTicketRequest{
				Creator:    "testuser",
				TemplateID: 1,
			},
			wantErr: true,
		},
		{
			name: "missing creator",
			request: &models.CreateTicketRequest{
				Name:       "Test Ticket",
				TemplateID: 1,
			},
			wantErr: true,
		},
		{
			name: "missing template_id",
			request: &models.CreateTicketRequest{
				Name:    "Test Ticket",
				Creator: "testuser",
			},
			wantErr: true,
		},
		{
			name: "empty name",
			request: &models.CreateTicketRequest{
				Name:       "",
				Creator:    "testuser",
				TemplateID: 1,
			},
			wantErr: true,
		},
		{
			name: "empty creator",
			request: &models.CreateTicketRequest{
				Name:       "Test Ticket",
				Creator:    "",
				TemplateID: 1,
			},
			wantErr: true,
		},
		{
			name: "zero template_id",
			request: &models.CreateTicketRequest{
				Name:       "Test Ticket",
				Creator:    "testuser",
				TemplateID: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateTicketRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreateTicketRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateApprovalRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *models.ApprovalRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &models.ApprovalRequest{
				ID:           1,
				ApprovalUser: "approver",
				Operation:    "approve",
				NextStep:     "next_step",
			},
			wantErr: false,
		},
		{
			name: "invalid operation",
			request: &models.ApprovalRequest{
				ID:           1,
				ApprovalUser: "approver",
				Operation:    "invalid",
				NextStep:     "next_step",
			},
			wantErr: true,
		},
		{
			name: "missing id",
			request: &models.ApprovalRequest{
				ApprovalUser: "approver",
				Operation:    "approve",
				NextStep:     "next_step",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateApprovalRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateApprovalRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
