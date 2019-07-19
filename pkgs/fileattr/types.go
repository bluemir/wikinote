package fileattr

type FileAttr struct {
	Path  string // 1024
	Key   string // 256
	Value string // 2048
}

type Store interface {
	//Query(query string) WhereClause
	Where(*Options) WhereClause
	SortBy(OrderType, OrderDirection) SortByClause
	Limit(int) LimitClause
	Find() ([]FileAttr, error)

	Path(string) PathClause
}
type WhereClause interface {
	SortBy(OrderType, OrderDirection) SortByClause
	Limit(int) LimitClause
	Find() ([]FileAttr, error)

	Get() (*FileAttr, error)
}
type SortByClause interface {
	Limit(int) LimitClause
	Find() ([]FileAttr, error)
}
type LimitClause interface {
	Find() ([]FileAttr, error)
}
type OrderType int
type OrderDirection int

const (
	OrderTypeKey OrderType = iota
	OrderTypeValue

	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

// using by plugins
type PathClause interface {
	Get(key string) (string, error)
	Set(key, value string) error
	All() (map[string]string, error)
}
