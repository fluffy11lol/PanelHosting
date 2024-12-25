package billing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTariff_Fields(t *testing.T) {
	tariff := &Tariff{
		ID:    "1",
		Name:  "Standard",
		SSD:   100,
		CPU:   4,
		RAM:   16,
		Price: 2999,
	}

	assert.Equal(t, "1", tariff.GetID(), "ID should be '1'")
	assert.Equal(t, "Standard", tariff.GetName(), "Name should be 'Standard'")
	assert.Equal(t, int64(100), tariff.GetSSD(), "SSD should be 100")
	assert.Equal(t, int64(4), tariff.GetCPU(), "CPU should be 4")
	assert.Equal(t, int64(16), tariff.GetRAM(), "RAM should be 16")
	assert.Equal(t, int64(2999), tariff.GetPrice(), "Price should be 2999")
}

func TestGetTariffsResponse(t *testing.T) {
	tariff1 := &Tariff{
		ID:    "1",
		Name:  "Standard",
		SSD:   100,
		CPU:   4,
		RAM:   16,
		Price: 2999,
	}
	tariff2 := &Tariff{
		ID:    "2",
		Name:  "Premium",
		SSD:   200,
		CPU:   8,
		RAM:   32,
		Price: 4999,
	}

	response := &GetTariffsResponse{
		Tariffs: []*Tariff{tariff1, tariff2},
	}

	assert.Len(t, response.GetTariffs(), 2, "There should be 2 tariffs in the response")
	assert.Equal(t, "1", response.GetTariffs()[0].GetID(), "First tariff ID should be '1'")
	assert.Equal(t, "Premium", response.GetTariffs()[1].GetName(), "Second tariff name should be 'Premium'")
}

func TestTariff_EmptyFields(t *testing.T) {
	tariff := &Tariff{}

	assert.Equal(t, "", tariff.GetID(), "ID should be an empty string")
	assert.Equal(t, "", tariff.GetName(), "Name should be an empty string")
	assert.Equal(t, int64(0), tariff.GetSSD(), "SSD should be 0")
	assert.Equal(t, int64(0), tariff.GetCPU(), "CPU should be 0")
	assert.Equal(t, int64(0), tariff.GetRAM(), "RAM should be 0")
	assert.Equal(t, int64(0), tariff.GetPrice(), "Price should be 0")
}
