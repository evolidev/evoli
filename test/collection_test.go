package test

import (
	"github.com/evolidev/evoli/framework/use"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddItemToCollectionShouldHasToValue(t *testing.T) {
	collection := use.NewCollection[string, string]()

	collection.Add("foo", "bar")

	assert.True(t, collection.Has("foo"))
}

func TestGetItemFromCollectionShouldReturnDesiredOne(t *testing.T) {
	collection := use.NewCollection[string, string]()

	collection.Add("foo", "bar")

	assert.Exactly(t, "bar", collection.Get("foo"))
}

func TestGetItemFormCollectionWhichNotExistsShouldReturnNil(t *testing.T) {
	collection := use.NewCollection[string, string]()

	assert.Empty(t, collection.Get("some_not_existing_key"))
}

func TestCountShouldReturnAmountOfElementsInCollection(t *testing.T) {
	collection := use.NewCollection[string, string]()

	assert.Zero(t, collection.Count())

	collection.Add("foo", "bar")

	assert.Exactly(t, 1, collection.Count())
}

func TestAddAndRemoveItemFromCollection(t *testing.T) {
	collection := use.NewCollection[string, string]()

	collection.Add("foo", "bar")

	assert.True(t, collection.Has("foo"))

	collection.Remove("foo")

	assert.False(t, collection.Has("foo"))
}
