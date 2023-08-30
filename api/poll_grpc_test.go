package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_In(t *testing.T) {
	assert.True(t, In("only-for-migration-a99", automatedOkList...))
	assert.True(t, In("only-for-migration-a1", automatedOkList...))
	assert.True(t, In("only-for-migration", automatedOkList...))
	assert.True(t, In("only-for-migration-", automatedOkList...))
	assert.False(t, In("mulbarton-a2", automatedOkList...))
	assert.False(t, In("glenvista2-a24", automatedOkList...))

}
