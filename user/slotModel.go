package user

// import (
// 	"database/sql"
// )

// // Slot represents a parking slot with optional booking details.
// // sql.Null* types are used to encode NULLs cleanly to JSON for the frontend.
// // The JSON tags intentionally match snake_case to align with DB columns.
// type Slot struct {
// 	ID        int            `db:"id" json:"id"`
// 	UserID    sql.NullString `db:"user_id" json:"user_id"`
// 	Email     sql.NullString `db:"email" json:"email"`
// 	StartTime sql.NullTime   `db:"start_time" json:"start_time"`
// 	EndTime   sql.NullTime   `db:"end_time" json:"end_time"`
// 	Notified  bool           `db:"notified" json:"notified"`
// }
