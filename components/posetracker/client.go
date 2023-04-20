package posetracker

import (
	"context"

	"github.com/edaniels/golog"
	pb "go.viam.com/api/component/posetracker/v1"
	"go.viam.com/utils/protoutils"
	"go.viam.com/utils/rpc"

	rprotoutils "go.viam.com/rdk/protoutils"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
)

// client implements PoseTrackerServiceClient.
type client struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	name   string
	client pb.PoseTrackerServiceClient
	logger golog.Logger
}

// NewClientFromConn constructs a new Client from connection passed in.
func NewClientFromConn(ctx context.Context, conn rpc.ClientConn, name resource.Name, logger golog.Logger) (PoseTracker, error) {
	c := pb.NewPoseTrackerServiceClient(conn)
	return &client{
		Named:  name.AsNamed(),
		name:   name.ShortNameForClient(),
		client: c,
		logger: logger,
	}, nil
}

func (c *client) Poses(
	ctx context.Context, bodyNames []string, extra map[string]interface{},
) (BodyToPoseInFrame, error) {
	ext, err := protoutils.StructToStructPb(extra)
	if err != nil {
		return nil, err
	}
	req := &pb.GetPosesRequest{
		Name:      c.name,
		BodyNames: bodyNames,
		Extra:     ext,
	}
	resp, err := c.client.GetPoses(ctx, req)
	if err != nil {
		return nil, err
	}
	result := BodyToPoseInFrame{}
	for key, pf := range resp.GetBodyPoses() {
		result[key] = referenceframe.ProtobufToPoseInFrame(pf)
	}
	return result, nil
}

func (c *client) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	return Readings(ctx, c)
}

func (c *client) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return rprotoutils.DoFromResourceClient(ctx, c.client, c.name, cmd)
}
