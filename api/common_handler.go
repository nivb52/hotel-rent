package api

import "github.com/nivb52/hotel-rent/types"

func newResourceResp(data any, n int64, p int) *types.ResourceResp {
	return &types.ResourceResp{
		Data:  data,
		Total: n,
		Page:  p,
	}
}
