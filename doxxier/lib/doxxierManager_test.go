package lib

import (
	"testing"

	"doxxier.tech/doxxier/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestDoxxierManager_CreateDoxxier(t *testing.T) {
	// Initialize a DoxxierManager
	dm := NewDoxxierManager(
		(func(string) {}),
	)

	// Check if the returned Doxxier is the same as the one in DoxxierManager
	if !compareDoxxiers(*dm.doxxier, *dm.GetDoxxier()) {
		t.Errorf("Doxxier in DoxxierManager and created Doxxier are not the same")
	}
}

func compareDoxxiers(d1, d2 models.Doxxier) bool {
	// Implement comparison logic here
	return d1.Id == d2.Id // Example comparison based on ID
}

func TestDoxxierManager_AddPart(t *testing.T) {
	dm := NewDoxxierManager(func(string) {})
	part := models.NewDoxxierPart()
	part.Id = "part1"

	assert.Equal(t, dm.GetDoxxier().Parts[len(dm.GetDoxxier().Parts)-1], part.Id)
}
func TestNewDoxxierManager(t *testing.T) {
	dm := NewDoxxierManager(func(s string) {})

	assert.NotNil(t, dm, "DoxxierManager should not be nil")
	assert.NotNil(t, dm.GetDoxxier(), "Doxxier should not be nil")
	assert.Equal(t, 0, len(dm.GetDoxxier().Parts), "Doxxier should have no parts initially")
}
