package app

import (
	"context"

	"github.com/sigmonsays/jobd/api"
)

func (me *Api) Ping(context context.Context, params api.PingParams) (api.PingRes, error) {
	ret := &api.PingResponse{}

	// todo

	return ret, nil
}
