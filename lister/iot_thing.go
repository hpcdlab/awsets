package lister

import (
	"fmt"

	"github.com/trek10inc/awsets/context"

	"github.com/trek10inc/awsets/resource"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/iot"
)

type AWSIoTThing struct {
}

func init() {
	i := AWSIoTThing{}
	listers = append(listers, i)
}

func (l AWSIoTThing) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.IoTThing}
}

func (l AWSIoTThing) List(ctx context.AWSetsCtx) (*resource.Group, error) {

	svc := iot.New(ctx.AWSCfg)
	rg := resource.NewGroup()
	var nextToken *string
	for {
		things, err := svc.ListThingsRequest(&iot.ListThingsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}).Send(ctx.Context)
		if err != nil {
			return rg, fmt.Errorf("failed to list iot thing: %w", err)
		}
		for _, thing := range things.Things {
			r := resource.New(ctx, resource.IoTThing, thing.ThingName, thing.ThingName, thing)
			r.AddRelation(resource.IoTThingType, thing.ThingTypeName, "")
			rg.AddResource(r)
		}
		if things.NextToken == nil {
			break
		}
		nextToken = things.NextToken
	}
	return rg, nil
}
