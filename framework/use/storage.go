package use

import (
	"embed"
	"github.com/evolidev/evoli/framework/filesystem"
)

func init() {
	myStores = NewCollection[string, filesystem.Store]()
}

var embedFs embed.FS

var myStores *Collection[string, filesystem.Store]

func Storage(stores ...string) filesystem.Store {
	var store string
	if len(stores) == 0 {
		cnf := Config("storage")
		cnf.SetDefault("default", "embed")
		store = cnf.Get("default").Value().(string)
	} else {
		store = stores[0]
	}

	if myStores.Has(store) {
		return myStores.Get(store)
	}

	myStores.Add(store, filesystem.NewEmbedFS(embedFs))

	return Storage(stores...)
}

func Embed(toEmbed embed.FS) {
	embedFs = toEmbed
	//fs.GlobFS()
}
