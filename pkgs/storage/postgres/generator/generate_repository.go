package generator

import (
	"gitea/pcp-inariam/inariam/pkgs/storage/postgres/entites"

	"gorm.io/gen"
	"gorm.io/gorm"
)

type Querier interface {
	// ADD your Dynamic Queries here.
}

func GenerateServices(db *gorm.DB, outPath string) {
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		// generation mode
		Mode: gen.WithoutContext |
			gen.WithDefaultQuery |
			gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.ApplyBasic(
		entites.Users{},
		entites.Groups{},
		entites.Roles{},
		entites.Permissions{},
		entites.Teams{},
		entites.Accounts{},
	)

	g.ApplyInterface(func(Querier) {}, entites.Users{})

	// Generate the code
	g.Execute()
}
