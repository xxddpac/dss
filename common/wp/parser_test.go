package wp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParserWeakPasswordList(t *testing.T) {
	assert.NotEqual(t, 0, len(WeakPasswordList))
}
