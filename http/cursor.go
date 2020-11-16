package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type APICursor struct {
	Total        int `default:"0"`
	Size         int `default:"500"`
	From         int `default:"0"`
	OrderBy      string
	OrderReverse bool
}

func (cur *APICursor) Sql(prefix string) string {
	if len(cur.OrderBy) == 0 {
		return ""
	}
	sql := "ORDER BY " + prefix + cur.OrderBy
	if cur.OrderReverse {
		sql += " DESC"
	} else {
		sql += " ASC"
	}
	return sql
}

func (cur *APICursor) SqlOrderBy(prefix string) string {
	if len(cur.OrderBy) == 0 {
		return ""
	}
	sql := prefix + cur.OrderBy
	if cur.OrderReverse {
		sql += " ASC"
	} else {
		sql += " DESC"
	}
	return sql
}

func (cur *APICursor) ApplyToResponse(w http.ResponseWriter) {
	if cur == nil {
		return
	}
	w.Header().Set("X-Page-Offset", fmt.Sprintf("%v", cur.From))
	w.Header().Set("X-Page-Limit", fmt.Sprintf("%v", cur.Size))
	w.Header().Set("X-Page-Total", fmt.Sprintf("%v", cur.Total))
}

func ResolveCursor(r *http.Request, defSize int, defSort string) *APICursor {
	c := &APICursor{
		Total:        0,
		Size:         defSize,
		From:         0,
		OrderBy:      defSort,
		OrderReverse: false,
	}
	sv := r.URL.Query().Get("limit")
	if len(sv) == 0 {
		sv = r.Header.Get("X-Page-Limit")
	}
	v, err := strconv.Atoi(sv)
	if err == nil && v > 0 {
		c.Size = v
	} else {
		c.Size = defSize
	}
	sv = r.URL.Query().Get("offset")
	if len(sv) == 0 {
		sv = r.Header.Get("X-Page-Offset")
	}
	v, err = strconv.Atoi(sv)
	if err == nil && v > 0 {
		c.From = v
	}
	sr := strings.TrimSpace(r.URL.Query().Get("sort"))
	if len(sr) == 0 {
		sr = strings.TrimSpace(r.Header.Get("X-Sort-By"))
	}
	if len(sr) == 0 {
		sr = strings.TrimSpace(defSort)
	}
	if len(sr) > 0 {
		if sr[0] == '+' ||  sr[0] == '-' {
			c.OrderReverse = sr[0] == '-'
			sr = sr[1:]
		}
		c.OrderBy = sr
	}
	return c
}
