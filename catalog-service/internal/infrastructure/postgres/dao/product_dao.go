package dao

type ProductRow struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
	CategoryID  int64   `db:"category_id"`
}
