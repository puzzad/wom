package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("410a6wyrxhi89cg")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("410a6wyrxhi89cg")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.games.id ?= id || @request.auth.id = user.id || @request.headers.x_pocketbase_game = code")

		return dao.SaveCollection(collection)
	})
}
