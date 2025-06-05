package dao

type CategoryRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
