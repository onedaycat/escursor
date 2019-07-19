package escursor

import (
	"encoding/base64"
)

//go:generate msgp
type Cursorfields []interface{}

type sortField struct {
	key     string
	options map[string]interface{}
}

//msgp:ignore QueryCursorBuilder
type QueryCursorBuilder struct {
	sortFields []sortField
	size       int
	query      map[string]interface{}
	token      string
	tokenData  []interface{}
}

func NewBuilder() *QueryCursorBuilder {
	return &QueryCursorBuilder{
		sortFields: make([]sortField, 0, 5),
		size:       11,
	}
}

func (c *QueryCursorBuilder) Token(token string) *QueryCursorBuilder {
	c.token = token

	return c
}

func (c *QueryCursorBuilder) Size(size int) *QueryCursorBuilder {
	c.size = size + 1

	return c
}

func (c *QueryCursorBuilder) Sort(key string, options map[string]interface{}) *QueryCursorBuilder {
	c.sortFields = append(c.sortFields, sortField{
		key:     key,
		options: options,
	})

	return c
}

func (c *QueryCursorBuilder) Query(query map[string]interface{}) *QueryCursorBuilder {
	c.query = query

	return c
}

func (c *QueryCursorBuilder) Build() map[string]interface{} {
	if c.query == nil {
		c.query = make(map[string]interface{}, 10)
	}
	c.query["size"] = c.size

	if c.token != "" {
		cfs := parseToken(c.token)
		if len(cfs) > 0 {
			searchAfter := make([]interface{}, len(cfs))
			for i, cf := range cfs {
				searchAfter[i] = cf
			}
			c.query["search_after"] = searchAfter
		}
	}

	if len(c.sortFields) > 0 {
		sorts := make([]map[string]interface{}, len(c.sortFields))
		for i, sf := range c.sortFields {
			sorts[i] = map[string]interface{}{
				sf.key: sf.options,
			}
		}

		c.query["sort"] = sorts
	}

	return c.query
}

func parseToken(token string) Cursorfields {
	cfByte, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil
	}

	cfs := &Cursorfields{}
	if _, err = cfs.UnmarshalMsg(cfByte); err != nil {
		return nil
	}

	return *cfs
}
