package collection

import (
	"encoding/json"
	"fmt"
	"go.k6.io/k6/js/modules"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

func init() {
	modules.Register("k6/x/collection", new(COLLECTION))
}

// COLLECTION is the k6 extension
type COLLECTION struct {
	items map[string]Object
}

type Object struct {
	Name           string
	ObjectContents map[string]ObjectContent
}

type ObjectContent struct {
	Name string
	Size int64
	Data []byte
}

func (c *COLLECTION) CreateCollection(collectionPath string) {
	fmt.Println("Creating collection!")
	c.items = make(map[string]Object)
	err := filepath.WalkDir(collectionPath,
		func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				//fmt.Println("Found file: ", path)
				_, ok := c.items[filepath.Dir(path)]
				if ok {
					c.items[filepath.Dir(path)].ObjectContents[path] = ObjectContent{
						Name: filepath.Base(path),
					}
				} else {
					oc := make(map[string]ObjectContent)
					oc[path] = ObjectContent{
						Name: filepath.Base(path),
					}
					c.items[filepath.Dir(path)] = Object{
						Name:           filepath.Dir(path),
						ObjectContents: oc,
					}
				}

			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

}

func (c *COLLECTION) GetRandomItem() Object {
	key := randMapKey(c.items)
	c.fillData(c.items[key])
	return c.items[key]
}

func (c *COLLECTION) fillData(o Object) {
	collectionJSON, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Got an object %s\n", string(collectionJSON))
	for k, v := range o.ObjectContents {
		dat, err := os.ReadFile(k)
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}
		v.Data = dat
		v.Size = int64(len(dat))
		c.items[filepath.Dir(k)].ObjectContents[k] = v
	}
}

func randMapKey(m map[string]Object) string {
	mapKeys := make([]string, 0, len(m)) // pre-allocate exact size
	for key := range m {
		mapKeys = append(mapKeys, key)
	}
	return mapKeys[rand.Intn(len(mapKeys))]
}
