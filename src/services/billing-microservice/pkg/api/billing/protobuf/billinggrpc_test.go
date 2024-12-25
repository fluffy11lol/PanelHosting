package billing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetTariffsSuccess(t *testing.T) {
	mockServer := &mockUserServiceServer{}

	ctx := context.Background()
	req := &emptypb.Empty{}

	resp, _ := mockServer.GetTariffs(ctx, req)

	assert.Equal(t, 1, len(resp.Tariffs))
	assert.Equal(t, "1", resp.Tariffs[0].ID)
	assert.NotEqual(t, "Basic", resp.Tariffs[0].Name)
	assert.Equal(t, int64(100), resp.Tariffs[0].SSD)
}

func TestGetTariffsError(t *testing.T) {
	mockServer := &mockUserServiceServer{}

	ctx := context.Background()

	resp, _ := mockServer.GetTariffs(ctx, nil)

	assert.NotNil(t, resp)
}
