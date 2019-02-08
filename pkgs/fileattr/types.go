package fileattr

type FileAttr struct {
	Path      string
	Namespace string
	Key       string
	Value     string
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
	Get(namespaceKey string) (string, error)
	Set(namespaceKey, value string) error
	All(namespace string) (map[string]string, error)
}
