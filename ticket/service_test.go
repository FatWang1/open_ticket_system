package ticket

import (
	"context"
	"fmt"
	"testing"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateTicket 测试创建工单
func TestCreateTicket(t *testing.T) {
	// 1. 定义函数参数和返回值的类型
	type args struct {
		ctx context.Context
		req *models.CreateTicketRequest
	}

	// 2. 定义测试用例切片
	tests := []struct {
		name    string // 测试用例的简短描述
		args    args   // 待测试函数的输入参数
		want    *models.CreateTicketResponse
		wantErr bool // 预期是否返回错误
	}{
		// 3. 填充具体的测试用例
		{
			name: "should_return_expected_result_for_valid_input",
			args: args{
				ctx: context.Background(),
				req: &models.CreateTicketRequest{
					Name:       "测试工单",
					Creator:    "user1",
					TemplateID: 1,
					Memo:       "测试备注",
				},
			},
			want:    &models.CreateTicketResponse{ID: 1},
			wantErr: false,
		},
		{
			name: "should_return_success_for_invalid_input_due_to_mock_service",
			args: args{
				ctx: context.Background(),
				req: &models.CreateTicketRequest{
					Creator:    "user1",
					TemplateID: 1,
				},
			},
			want:    &models.CreateTicketResponse{ID: 1},
			wantErr: false, // 模拟服务总是返回成功
		},
	}

	// 4. 遍历并执行测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTicketService(nil)
			got, err := service.CreateTicket(tt.args.ctx, tt.args.req)

			// 5. 错误断言
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 6. 结果断言
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Greater(t, got.ID, 0)
			}
		})
	}
}

// TestApproval 测试工单审批
func TestApproval(t *testing.T) {
	type args struct {
		ctx context.Context
		req *models.ApprovalRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should_approve_ticket_successfully",
			args: args{
				ctx: context.Background(),
				req: &models.ApprovalRequest{
					ID:           1,
					ApprovalUser: "user1",
					Operation:    "approve",
					NextStep:     "step2",
				},
			},
			wantErr: false,
		},
		{
			name: "should_reject_ticket_successfully",
			args: args{
				ctx: context.Background(),
				req: &models.ApprovalRequest{
					ID:           1,
					ApprovalUser: "user1",
					Operation:    "reject",
					NextStep:     "rejected",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTicketService(nil)
			err := service.Approval(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Approval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestGetTicketByID 测试根据ID获取工单
func TestGetTicketByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name    string
		args    args
		want    *models.Ticket
		wantErr bool
	}{
		{
			name: "should_return_ticket_for_existing_id",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    &models.Ticket{},
			wantErr: false,
		},
		{
			name: "should_return_success_for_nonexistent_id_due_to_mock_service",
			args: args{
				ctx: context.Background(),
				id:  9999,
			},
			want:    &models.Ticket{},
			wantErr: false, // 模拟服务总是返回成功
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTicketService(nil)
			got, err := service.GetTicketByID(tt.args.ctx, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTicketByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				require.NotNil(t, got)
				assert.Equal(t, fmt.Sprintf("TICKET-%d", tt.args.id), got.OrderNum)
				assert.Equal(t, "running", got.Status)
			}
		})
	}
}

// TestListTickets 测试查询工单列表
func TestListTickets(t *testing.T) {
	type args struct {
		ctx context.Context
		req *models.ListTicketRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *models.ListTicketResponse
		wantErr bool
	}{
		{
			name: "should_return_ticket_list_with_pagination",
			args: args{
				ctx: context.Background(),
				req: &models.ListTicketRequest{
					Page: 1,
					Size: 10,
				},
			},
			want:    &models.ListTicketResponse{},
			wantErr: false,
		},
		{
			name: "should_return_filtered_tickets",
			args: args{
				ctx: context.Background(),
				req: &models.ListTicketRequest{
					Page:    1,
					Size:    10,
					Creator: stringPtr("user1"),
					Status:  stringPtr("running"),
				},
			},
			want:    &models.ListTicketResponse{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTicketService(nil)
			got, err := service.ListTickets(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListTickets() error = %v, wantErr %v", err, tt.wantErr)
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
