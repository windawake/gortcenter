package models

const STATUS_START = 0
const STATUS_COMMIT = 1
const STATUS_ROLLBACK = 2

const ACTION_START = 0
const ACTION_END = 1
const ACTION_PREPARE = 2
const ACTION_PREPARE_COMMIT = 3
const ACTION_PREPARE_ROLLBACK = 4
const ACTION_COMMIT = 5
const ACTION_ROLLBACK = 6

type ResetTransact struct {
	Id               int    `db:"id"`
	TransactId       string `db:"transact_id"`
	TransactRollback string `db:"transact_rollback"`
	Action           int    `db:"action"`
	XidsInfo         string `db:"xids_info"`
	CreatedAt        string `db:"created_at"`
}
