package gormdb

import(
	"fmt"
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/snpavlov/gorm-probe/domain"

)

// https://github.com/go-gormigrate/gormigrate
// https://gorm.io/docs/query.html

type GormDBContext struct {
}

func (dctx GormDBContext) Open(connection string, dbschema string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{ 
				TablePrefix: dbschema + ".",
			},
	})
	if err != nil {
		return nil, fmt.Errorf("can't open database! Error: %v", err)
	}

	return db, nil
}

func (dctx GormDBContext) Migrate(connection string, dbschema string) error {
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{ 
				TablePrefix: dbschema + ".",
				SingularTable: true, 
				NoLowerCase: true,
			},
	})
	if err != nil {
		return fmt.Errorf("can't open database! Error: %v", err)
	}

	// Migrate the schema
    err = db.AutoMigrate(&domain.Play{}, &domain.Actor{}, &domain.Showing{})
	if err != nil {
		return fmt.Errorf("can't migrate database! Error: %v", err)
	}

	return nil
}
