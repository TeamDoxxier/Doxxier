package Foundation

import (
	"testing"

	"doxxier.tech/doxxier/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestDoxxierManager_CreateDoxxier(t *testing.T) {
	// Initialize a DoxxierManager
	dm := InitialiseDoxxierManager()

	// Call CreateDoxxier method
	createdDoxxier := dm.CreateDoxxier()

	// Check if the returned Doxxier is the same as the one in DoxxierManager
	if !compareDoxxiers(createdDoxxier, dm.Doxxier) {
		t.Errorf("Doxxier in DoxxierManager and created Doxxier are not the same")
	}
}

func compareDoxxiers(d1, d2 models.Doxxier) bool {
	// Implement comparison logic here
	return d1.Id == d2.Id // Example comparison based on ID
}

func TestDoxxierManager_AddPart(t *testing.T) {
	dm := NewDoxxierManager()
	part := models.NewDoxxierPart()
	part.Id = "part1"

	assert.Equal(t, dm.Doxxier.Parts[len(dm.Doxxier.Parts)-1], part.Id)
}
func TestNewDoxxierManager(t *testing.T) {
	dm := NewDoxxierManager()

	assert.NotNil(t, dm, "DoxxierManager should not be nil")
	assert.NotNil(t, dm.Doxxier, "Doxxier should not be nil")
	assert.Equal(t, 0, len(dm.Doxxier.Parts), "Doxxier should have no parts initially")
}
