package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	// "github.com/NeurologiaLogic/gochat/graph/model"
	"github.com/NeurologiaLogic/gochat/database"
)
type Resolver struct{
	db *database.DB
}

func ConfigureResolver() Config {
	return Config{
		 Resolvers: &Resolver{
				db: database.GetDB(),
		 },
	}
}