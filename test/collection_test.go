package test

import (
	"github.com/evolidev/evoli/framework/use"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddItemToCollectionShouldHasToValue(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()

	collection.Add("foo", "bar")

	assert.True(t, collection.Has("foo"))
}

func TestGetItemFromCollectionShouldReturnDesiredOne(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()

	collection.Add("foo", "bar")

	assert.Exactly(t, "bar", collection.Get("foo"))
}

func TestGetItemFormCollectionWhichNotExistsShouldReturnNil(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()

	assert.Empty(t, collection.Get("some_not_existing_key"))
}

func TestCountShouldReturnAmountOfElementsInCollection(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()

	assert.Zero(t, collection.Count())

	collection.Add("foo", "bar")

	assert.Exactly(t, 1, collection.Count())
}

func TestAddAndRemoveItemFromCollection(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()

	collection.Add("foo", "bar")

	assert.True(t, collection.Has("foo"))

	collection.Remove("foo")

	assert.False(t, collection.Has("foo"))
}

func TestNextShouldReturnNextItem(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()
	collection.Add("foo", "bar")
	collection.Add("foo2", "bar2")

	assert.Exactly(t, "bar", collection.Next())
	assert.Exactly(t, "bar2", collection.Next())
}

func TestPreviousShouldReturnPreviousItem(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()
	collection.Add("foo", "bar")
	collection.Add("foo2", "bar2")
	collection.Next()
	collection.Next()

	assert.Exactly(t, "bar", collection.Previous())
}

func TestFirstShouldReturnFirstItem(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()
	collection.Add("foo", "bar")
	collection.Add("foo2", "bar2")
	collection.Add("foo3", "bar3")
	collection.Next()
	collection.Next()

	assert.Exactly(t, "bar", collection.First())
}

func TestKeyShouldReturnCurrentKeyItem(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()
	collection.Add("foo", "bar")

	assert.Exactly(t, "foo", collection.Key())
}

func TestCurrentShouldReturnCurrentValue(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()
	collection.Add("foo", "bar")

	assert.Exactly(t, "bar", collection.Current())
}

func TestLastShouldReturnLastElement(t *testing.T) {
	t.Parallel()
	collection := use.NewCollection[string, string]()
	collection.Add("foo", "bar")
	collection.Add("foo2", "bar2")
	collection.Add("foo3", "bar3")
	collection.Next()
	collection.Next()

	assert.Exactly(t, "bar3", collection.Last())
}

func TestHasNext(t *testing.T) {
	t.Parallel()
	t.Run("Has next should return true if there is a next element", func(t *testing.T) {
		collection := use.NewCollection[string, string]()
		collection.Add("foo", "bar")

		assert.True(t, collection.HasNext())
	})

	t.Run("Has next should return false if there is no next element", func(t *testing.T) {
		collection := use.NewCollection[string, string]()
		collection.Add("foo", "bar")
		collection.Next()

		assert.False(t, collection.HasNext())
	})
}

func TestHasPrevious(t *testing.T) {
	t.Parallel()
	t.Run("Has previous should return true if there is a previous element", func(t *testing.T) {
		collection := use.NewCollection[string, string]()
		collection.Add("foo", "bar")
		collection.Add("foo", "bar")
		collection.Next()
		collection.Next()

		assert.True(t, collection.HasPrevious())
	})

	t.Run("Has previous should return false if there is no previous element", func(t *testing.T) {
		collection := use.NewCollection[string, string]()
		collection.Add("foo", "bar")

		assert.False(t, collection.HasPrevious())
	})
}

func TestMerge(t *testing.T) {
	t.Parallel()
	t.Run("Merge should merge to collection together", func(t *testing.T) {
		collection := use.NewCollection[string, string]()
		collection.Add("foo", "bar")

		collection2 := use.NewCollection[string, string]()
		collection2.Add("foo2", "bar2")

		collection.Merge(collection2)

		assert.True(t, collection.Has("foo"))
		assert.True(t, collection.Has("foo2"))
	})
}
