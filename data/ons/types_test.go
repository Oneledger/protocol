/*

 */

package ons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName_IsSub(t *testing.T) {

	name := "abc.xyzzz.ol"
	assert.True(t, GetNameFromString(name).IsSub())
}

func TestName_GetParentName(t *testing.T) {

	name := "abc.xyzzz.ol"
	_, err := GetNameFromString(name).GetParentName()
	assert.NoError(t, err)
}
