package shared

import "gopkg.in/guregu/null.v4"

type AnimeUpdateParams struct {
	ID               null.Int
	Title            null.String
	Status           null.String
	Progress         null.Int
	Score            null.Int
	ShouldParseTitle null.Bool
	IsIncrement      null.Bool
}
