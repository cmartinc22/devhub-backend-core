package models

type Metric struct {
	Kind     string `db:"kind"`
	Quantity string `db:"quantity"`
}
