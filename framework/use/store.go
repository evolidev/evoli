package use

import (
	"embed"
	"github.com/evolidev/evoli/framework/filesystem"
	"io/fs"
	"os"
)

func init() {
	myStores = NewCollection[string, filesystem.Store]()
}

var myStores *Collection[string, filesystem.Store]

func Store(stores ...string) filesystem.Store {
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

	myStores.Add(store, filesystem.NewFS(getStorage()))

	return Store(stores...)
}

func Embed(toEmbed embed.FS) {
	if myStores.Has("embed") {
		return
	}

	myStores.Add("embed", filesystem.NewFS(toEmbed))
}

func getStorage() fs.FS {
	return os.DirFS(BasePath())
}
