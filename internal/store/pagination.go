package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit              int      `json:"limit" validate:"gte=1,lte=20"`
	Offset             int      `json:"offset" validate:"gte=0"`
	Sort               string   `json:"sort" validate:"oneof=asc desc"`
	Tags               []string `json:"tags" validate:"max=5"`
	TitleContentSearch string   `json:"search" validate:"max=100"`
	Since              string   `json:"since"`
	Until              string   `json:"until"`
}

func (pf *PaginatedFeedQuery) Parse(r *http.Request) error {
	rq := r.URL.Query()

	limit := rq.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		pf.Limit = l
	}

	offset := rq.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return err
		}
		pf.Offset = o
	}

	sort := rq.Get("sort")
	if sort != "" {
		pf.Sort = sort
	}

	tags := rq.Get("tags")
	if tags != "" {
		pf.Tags = strings.Split(tags, ",")
	} else {
		pf.Tags = []string{}
	}

	search := rq.Get("search")
	if search != "" {
		pf.TitleContentSearch = search
	}

	since := rq.Get("since")
	if since != "" {
		pf.Since = parseTime(since)
	}

	until := rq.Get("until")
	if until != "" {
		pf.Until = parseTime(until)
	}

	return nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
