package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chandler-WQ/experiment/db"
)

func TestMain(m *testing.M) {
	db.MustInitDb()
	m.Run()
}

func TestSegmentCreate(t *testing.T) {
	err := SegmentCreate()
	assert.Nil(t, err)
}
