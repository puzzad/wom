package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("410a6wyrxhi89cg")
		if err != nil {
			return err
		}

		collection.Name = "games2"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("410a6wyrxhi89cg")
		if err != nil {
			return err
		}

		collection.Name = "games"

		return dao.SaveCollection(collection)
	})
}
