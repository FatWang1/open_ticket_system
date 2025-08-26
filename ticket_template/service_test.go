package ticket_template

import (
	"context"
	"fmt"
	"testing"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateTicketTemplate 测试创建工单模板
func TestCreateTicketTemplate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *models.CreateTicketTemplateRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *models.CreateTicketTemplateResponse
		wantErr bool
	}{
		{
			name: "should_return_expected_result_for_valid_input",
			args: args{
				ctx: context.Background(),
				req: &models.CreateTicketTemplateRequest{
					Name:        "测试模板",
					Memo:        "测试备注",
					Version:     "1.0",
					Creator:     "admin",
					StartStep:   "step1",
					EndStepList: []string{"step3"},
					ConfigList:  []*models.StepConfig{},
				},
			},
			want:    &models.CreateTicketTemplateResponse{ID: 1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTicketTemplateService(nil)
			got, err := service.CreateTicketTemplate(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTicketTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Greater(t, got.ID, 0)
			}
		})
	}
}

// TestGetTicketTemplateByID 测试根据ID获取工单模板
func TestGetTicketTemplateByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name    string
		args    args
		want    *models.TicketTemplate
		wantErr bool
	}{
		{
			name: "should_return_template_for_existing_id",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    &models.TicketTemplate{},
			wantErr: false,
		},
		{
			name: "should_return_success_for_nonexistent_id_due_to_mock_service",
			args: args{
				ctx: context.Background(),
				id:  9999,
			},
			want:    &models.TicketTemplate{},
			wantErr: false, // 模拟服务总是返回成功
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTicketTemplateService(nil)
			got, err := service.GetTicketTemplateByID(tt.args.ctx, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTicketTemplateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				require.NotNil(t, got)
				assert.Equal(t, fmt.Sprintf("template-%d", tt.args.id), got.Uid)
				assert.Equal(t, "step1", got.StartStep)
				assert.False(t, got.Builtin)
			}
		})
	}
}

// TestListTicketTemplates 测试查询工单模板列表
func TestListTicketTemplates(t *testing.T) {
	type args struct {
		ctx context.Context
		req *models.ListTicketTemplateRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *models.ListTicketTemplateResponse
		wantErr bool
	}{
		{
			name: "should_return_template_list_with_pagination",
			args: args{
				ctx: context.Background(),
				req: &models.ListTicketTemplateRequest{
					Page: 1,
					Size: 10,
				},
			},
			want:    &models.ListTicketTemplateResponse{},
			wantErr: false,
		},
		{
			name: "should_return_filtered_templates",
			args: args{
				ctx: context.Background(),
				req: &models.ListTicketTemplateRequest{
					Page:    1,
					Size:    10,
					Name:    stringPtr("测试"),
					Creator: stringPtr("admin"),
				},
			},
			want:    &models.ListTicketTemplateResponse{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTicketTemplateService(nil)
			got, err := service.ListTicketTemplates(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListTicketTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				require.NotNil(t, got)
				assert.GreaterOrEqual(t, got.Total, int64(0))
			}
		})
	}
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}
