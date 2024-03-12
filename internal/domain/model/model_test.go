package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMigrationInfo_FromMigrationFileInfo(t *testing.T) {
	info := MigrationFileInfo{Version: "vetsion_1"}

	actual := new(MigrationInfo).FromMigrationFileInfo(info, time.Now().Add(-10*time.Second), "something error")

	assert.Equal(t, actual.ExecutionTime, uint64(10))
	assert.Equal(t, actual.Version, info.Version)
	assert.Equal(t, actual.Error, "something error")
}
