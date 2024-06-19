package tax

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
	amount := 500.0
	expected := 5.0

	result, _ := CalculateTax(amount)

	if result != expected {
		t.Errorf("Expected %f but got %f", expected, result)
	}
}

func TestCalculateTax2(t *testing.T) {
	tax, err := CalculateTax(1000.0)
	assert.Nil(t, err)
	assert.Equal(t, 10.0, tax)

	tax, err = CalculateTax(0)
	assert.Error(t, err, "amount must be greater than 0")
	assert.Equal(t, 0.0, tax)
	assert.Contains(t, err.Error(), "greater than 0")
}

func TestCalculateTaxAndSave(t *testing.T) {
	repository := &TaxRepositoryMock{}

	repository.On("SaveTax", 10.0).Return(nil)
	repository.On("SaveTax", 0.0).Return(errors.New("error saving tax"))

	err := CalculateTaxAndSave(1000.0, repository)
	assert.Nil(t, err)

	err = CalculateTaxAndSave(0.0, repository)
	assert.Error(t, err, "error saving tax")

	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "SaveTax", 2)
}

func TestCalculateTaxBatch(t *testing.T) {
	type calcTax struct {
		amount, expected float64
	}

	table := []calcTax{
		{500.0, 5.0},
		{1000.0, 10.0},
		{1500.0, 10.0},
	}

	for _, item := range table {
		result, _ := CalculateTax(item.amount)

		if result != item.expected {
			t.Errorf("Expected %f but got %f", item.expected, result)
		}
	}
}

func BenchmarkCalculateTax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateTax(500.0)
	}
}

func FuzzCalculateTaxBatch(f *testing.F) {
	seed := []float64{-1, -2, -2.5, 500.0, 1000.0, 1501.0}

	for _, amount := range seed {
		f.Add(amount)
	}

	f.Fuzz(func(t *testing.T, amount float64) {
		result, _ := CalculateTax(amount)

		if amount <= 0 && result != 0 {
			t.Errorf("Received %f but expected 0", result)
		}

		if amount > 20000 && result != 20 {
			t.Errorf("Received %f but expected 20", result)
		}
	})
}
