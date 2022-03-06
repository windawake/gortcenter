package models

type ResetTransactSql struct {
	Id             int    `db:"id"`
	RequestId      string `db:"request_id"`
	TransactId     string `db:"transact_id"`
	ChainId        string `db:"chain_id"`
	TransactStatus int    `db:"transact_status"`
	Connection     string `db:"connection"`
	Sql            string `db:"sql"`
	Result         int    `db:"result"`
	CheckResult    int    `db:"check_result"`
	CreatedAt      string `db:"created_at"`
}
