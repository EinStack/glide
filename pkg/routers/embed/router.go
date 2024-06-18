package embed

import (
	"context"
	"github.com/EinStack/glide/pkg/api/schemas"
	"github.com/EinStack/glide/pkg/routers"
	"github.com/EinStack/glide/pkg/routers/retry"
	"github.com/EinStack/glide/pkg/telemetry"
	"go.uber.org/zap"
)

type EmbeddingRouter struct {
	routerID routers.RouterID
	Config   *LangRouterConfig
	retry    *retry.ExpRetry
	tel      *telemetry.Telemetry
	logger   *zap.Logger
}

func (r *routers.LangRouter) Embed(ctx context.Context, req *schemas.EmbedRequest) (*schemas.EmbedResponse, error) {

}
