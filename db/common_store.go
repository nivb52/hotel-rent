package db

import "go.mongodb.org/mongo-driver/mongo/options"

func buildPaginationOpts(p *Pagination) *options.FindOptions {
	// set max results
	if p.Limit > 50 {
		p.Limit = 50
	}
	opts := options.FindOptions{}
	opts.SetSkip(int64((p.Page - 1) * p.Limit))
	opts.SetLimit(int64(p.Limit))

	return &opts
}
