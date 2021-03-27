package huawei

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	"github.com/sealyun/cloud-kernel/pkg/vars"
)

const (
	projectID = "06b275f707800f272feec0147ee22b32"
	zone      = "ap-southeast-1a"
)

type HClient struct {
	Zone      string
	EcsClient *ecs.EcsClient
	EipClient *eip.EipClient
}

func NewClientWithAccessKey(ak, sk string) *HClient {
	n := len(zone)
	ecsEndpoint := fmt.Sprintf("https://ecs.%s.myhuaweicloud.com", zone[:n-1])
	vpcEndpoint := fmt.Sprintf("https://vpc.%s.myhuaweicloud.com", zone[:n-1])
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithProjectId(projectID).
		Build()
	return &HClient{
		Zone: zone,
		EcsClient: ecs.NewEcsClient(
			ecs.EcsClientBuilder().
				WithEndpoint(ecsEndpoint).
				WithCredential(auth).
				Build()),
		EipClient: eip.NewEipClient(
			eip.EipClientBuilder().
				WithEndpoint(vpcEndpoint).
				WithCredential(auth).
				Build()),
	}
}

func (h *HClient) Describe(serverId string) (*model.ShowServerResponse, error) {

	client := h.EcsClient

	request := &model.ShowServerRequest{}
	request.ServerId = serverId

	return client.ShowServer(request)
}

func (h *HClient) RunInstances(amount int, dryRun bool, bandwidthOut bool) ([]string, error) {

	client := h.EcsClient
	request := &model.CreatePostPaidServersRequest{}
	var listPostPaidServerNicNicsPostPaidServer = []model.PostPaidServerNic{
		{
			SubnetId: "179d3994-1a14-44df-bfb1-bb1e5495bb45",
		},
	}
	var listPostPaidServerTagServerTagsPostPaidServer = []model.PostPaidServerTag{
		{
			Key:   "test",
			Value: "sealos",
		},
	}
	publicipPostPaidServer := &model.PostPaidServerPublicip{}
	if bandwidthOut {
		chargemodePostPaidServerEipBandwidth := "traffic"
		var serverEipBandwidth int32 = 100
		bandwidthPostPaidServerEip := &model.PostPaidServerEipBandwidth{
			Size:       &serverEipBandwidth,
			Sharetype:  model.GetPostPaidServerEipBandwidthSharetypeEnum().PER,
			Chargemode: &chargemodePostPaidServerEipBandwidth,
		}
		eipPostPaidServerPublicip := &model.PostPaidServerEip{
			Iptype:    "5_bgp",
			Bandwidth: bandwidthPostPaidServerEip,
		}
		publicipPostPaidServer = &model.PostPaidServerPublicip{
			Eip: eipPostPaidServerPublicip,
		}
	}
	countPostPaidServer := int32(amount)
	isAutoRenamePostPaidServer := false
	adminPassPostPaidServer := vars.EcsPassword
	var serverRootVolume int32 = 40
	rootVolumePostPaidServer := &model.PostPaidServerRootVolume{
		Volumetype: model.GetPostPaidServerRootVolumeVolumetypeEnum().SSD,
		Size:       &serverRootVolume,
	}
	serverCreatePostPaidServersRequestBody := &model.PostPaidServer{
		AvailabilityZone: h.Zone,
		FlavorRef:        "kc1.large.2",
		ImageRef:         "04678140-fcc1-465d-ba36-3a2b19d155f9",
		Name:             "sealos",
		Nics:             listPostPaidServerNicNicsPostPaidServer,
		Publicip:         publicipPostPaidServer,
		RootVolume:       rootVolumePostPaidServer,
		ServerTags:       &listPostPaidServerTagServerTagsPostPaidServer,
		Vpcid:            "085d5e2a-8339-4780-bb6b-8b418959fc9a",
		AdminPass:        &adminPassPostPaidServer,
		IsAutoRename:     &isAutoRenamePostPaidServer,
		Count:            &countPostPaidServer,
	}
	request.Body = &model.CreatePostPaidServersRequestBody{
		Server: serverCreatePostPaidServersRequestBody,
		DryRun: &dryRun,
	}
	response, err := client.CreatePostPaidServers(request)
	if err == nil {
		return *response.ServerIds, nil
	} else {
		return nil, err
	}
}

func (h *HClient) DeleteInstances(serverId []string, delPublicIp bool) (*model.DeleteServersResponse, error) {

	client := h.EcsClient

	request := &model.DeleteServersRequest{}
	var listServerIdServersDeleteServersRequestBody = make([]model.ServerId, 0)
	for _, v := range serverId {
		listServerIdServersDeleteServersRequestBody = append(listServerIdServersDeleteServersRequestBody, model.ServerId{Id: v})
	}
	request.Body = &model.DeleteServersRequestBody{
		DeletePublicip: &delPublicIp,
		Servers:        listServerIdServersDeleteServersRequestBody,
	}

	return client.DeleteServers(request)
}
