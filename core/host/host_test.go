package host

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefreshHost(t *testing.T) {
	RefreshHost()
	assert.NotEqual(t, 0, len(PrivateIPv4.Load().([]string)))
	assert.NotEqual(t, 0, len(Name.Load().(string)))
}
