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

type UserAnimeInfo struct {
	AnimeInfo
	UserID     int64     `db:"user_id"`
	Status     string    `db:"status"`
	Rewatches  uint16    `db:"rewatches"`
	Progress   uint16    `db:"progress"`
	Score      uint8     `db:"score"`
	StartDate  null.Time `db:"start_date"`
	FinishDate null.Time `db:"finish_date"`
}
type AnimeInfo struct {
	SourceID uint8  `db:"source_id"` // anime site ID
	ID       int64  `db:"id"`        // anime ID on source site
	Title    string `db:"title"`     // TODO: use a struct with jap/romaji/eng titles maybe
	Episodes uint16 `db:"episodes"`  // total amount of eps in the series
}
