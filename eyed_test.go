package eyed

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var dogId *EYEDType = RegisterType("dog", "dog")
var catId *EYEDType = RegisterType("cat", "cat")
var fishId *EYEDType = RegisterType("fish", "fi")

func TestEYEDType(t *testing.T) {

	t.Run("create new IDs", func(t *testing.T) {
		t.Parallel()

		id := dogId.New()
		idStr := id.String()
		assert.Equal(t, "dog", id.idType.name)
		assert.Len(t, idStr, 31)
		assert.Equal(t, "dog_", idStr[0:4])

		id = catId.New()
		idStr = id.String()
		assert.Equal(t, "cat", id.idType.name)
		assert.Len(t, idStr, 31)
		assert.Equal(t, "cat_", idStr[0:4])

		id = fishId.New()
		idStr = id.String()
		assert.Equal(t, "fish", id.idType.name)
		assert.Len(t, idStr, 30)
		assert.Equal(t, "fi_", idStr[0:3])
	})

	t.Run("Getter methods", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "fish", fishId.Name())
		assert.Equal(t, "fi", fishId.Prefix())
	})

	t.Run("Is()", func(t *testing.T) {
		t.Parallel()

		assert.True(t, dogId.Is(dogId.New().String()))
		assert.False(t, dogId.Is(fishId.New().String()))
		assert.False(t, dogId.Is("bad"))
	})
}

func TestGetType(t *testing.T) {
	id := dogId.New().String()
	idType, found := GetType(id)
	assert.True(t, found)
	assert.Equal(t, "dog", idType.name)

	id = "nounderscore"
	idType, found = GetType(id)
	assert.False(t, found)
	assert.Nil(t, idType)

	id = "too_many_underscores"
	idType, found = GetType(id)
	assert.False(t, found)
	assert.Nil(t, idType)

	id = "bad_" + ksuid.New().String()
	idType, found = GetType(id)
	assert.False(t, found)
	assert.Nil(t, idType)
}

func TestParse(t *testing.T) {

	id := dogId.New().String()
	parsedId, ok := Parse(id)
	assert.True(t, ok)
	assert.Equal(t, id, parsedId.String())

	id = "nounderscore"
	_, ok = Parse(id)
	assert.False(t, ok)

	id = "too_many_underscores"
	_, ok = Parse(id)
	assert.False(t, ok)

	id = "bad_" + ksuid.New().String()
	_, ok = Parse(id)
	assert.False(t, ok)

	id = "fi_invalidksuid"
	_, ok = Parse(id)
	assert.False(t, ok)
}

func TestEYED(t *testing.T) {
	id := dogId.New()
	assert.Equal(t, "dog", id.Type().Name())
	assert.Len(t, id.Ksuid().String(), 27)
}
