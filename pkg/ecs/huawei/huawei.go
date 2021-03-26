package huawei

import (
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	emodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/model"
)

const (
	//endpoint     = "https://ecs.ap-southeast-3.myhuaweicloud.com"
	//VpcEndpoint  = "https://vpc.ap-southeast-3.myhuaweicloud.com"
	//projectID    = "06b275f705800f262f3bc014ffcdbde1"
	defaultCount = 1
)

type HClient struct {
	Count            int32
	Ak               string
	Sk               string
	EcsEndpoint      string
	VpcEndpoint      string
	ProjectId        string
	AvailabilityZone string
	EcsClient        *ecs.EcsClient
	EipClient        *eip.EipClient
}

func GetDefaultHAuth(ak, sk, projectId, AvailabilityZone string) *HClient {
	n := len(AvailabilityZone)
	ecsEndpoint := fmt.Sprintf("https://ecs.%s.myhuaweicloud.com", AvailabilityZone[:n-1])
	vpcEndpoint := fmt.Sprintf("https://vpc.%s.myhuaweicloud.com", AvailabilityZone[:n-1])
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithProjectId(projectId).
		Build()
	return &HClient{
		EcsEndpoint:      ecsEndpoint,
		VpcEndpoint:      vpcEndpoint,
		ProjectId:        projectId,
		AvailabilityZone: AvailabilityZone,
		Ak:               ak,
		Sk:               sk,
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

func (h *HClient) Show(serverid string) {

	client := h.EcsClient

	request := &model.ShowServerRequest{}
	request.ServerId = serverid

	response, err := client.ShowServer(request)
	if err == nil {
		date, _ := json.MarshalIndent(response.Server, "", "    ")
		fmt.Println(string(date))
	} else {
		fmt.Println(err)
	}
}

func (h *HClient) GenerateEipServer(count, sizePostPaidServerEipBandwidth, sizePostPaidServerRootVolume int32, eip bool, FlavorRef, ImageRef, Vpcid, SubnetId, adminPass, keyName string) []string {

	client := h.EcsClient
	request := &model.CreatePostPaidServersRequest{}
	var listPostPaidServerNicNicsPostPaidServer = []model.PostPaidServerNic{
		{
			SubnetId: SubnetId,
		},
	}
	var listPostPaidServerTagServerTagsPostPaidServer = []model.PostPaidServerTag{
		{
			Key:   "test",
			Value: "sealos",
		},
	}
	publicipPostPaidServer := &model.PostPaidServerPublicip{}
	if eip {
		chargemodePostPaidServerEipBandwidth := "traffic"
		bandwidthPostPaidServerEip := &model.PostPaidServerEipBandwidth{
			Size:       &sizePostPaidServerEipBandwidth,
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
	countPostPaidServer := count
	isAutoRenamePostPaidServer := false
	keyNamePostPaidServer := keyName
	adminPassPostPaidServer := adminPass
	rootVolumePostPaidServer := &model.PostPaidServerRootVolume{
		Volumetype: model.GetPostPaidServerRootVolumeVolumetypeEnum().SSD,
		Size:       &sizePostPaidServerRootVolume,
	}
	serverCreatePostPaidServersRequestBody := &model.PostPaidServer{
		AvailabilityZone: h.AvailabilityZone,
		FlavorRef:        FlavorRef,
		ImageRef:         ImageRef,
		Name:             "sealos",
		Nics:             listPostPaidServerNicNicsPostPaidServer,
		Publicip:         publicipPostPaidServer,
		RootVolume:       rootVolumePostPaidServer,
		ServerTags:       &listPostPaidServerTagServerTagsPostPaidServer,
		Vpcid:            Vpcid,
		KeyName:          &keyNamePostPaidServer,
		AdminPass:        &adminPassPostPaidServer,
		IsAutoRename:     &isAutoRenamePostPaidServer,
		Count:            &countPostPaidServer,
	}
	request.Body = &model.CreatePostPaidServersRequestBody{
		Server: serverCreatePostPaidServersRequestBody,
	}

	response, err := client.CreatePostPaidServers(request)

	if err == nil {
		date, _ := json.MarshalIndent(response, "", "    ")
		fmt.Println(string(date))
		return *response.ServerIds
	} else {
		fmt.Println(err)
		return nil
	}
}

func (h *HClient) DeleteServer(serverId string, delPublicIp bool) {

	client := h.EcsClient

	request := &model.DeleteServersRequest{}
	var listServerIdServersDeleteServersRequestBody = []model.ServerId{
		{
			Id: serverId,
		},
	}

	request.Body = &model.DeleteServersRequestBody{
		DeletePublicip: &delPublicIp,
		Servers:        listServerIdServersDeleteServersRequestBody,
	}

	response, err := client.DeleteServers(request)

	if err == nil {
		date, _ := json.MarshalIndent(response, "", "    ")
		fmt.Println(string(date))
	} else {
		fmt.Println(err)
	}
}

func (h *HClient) ListServer() {
	client := h.EcsClient
	request := &model.ListServersDetailsRequest{}
	response, err := client.ListServersDetails(request)

	if err == nil {
		//fmt.Printf("%+v\n", response.Servers)
		date, _ := json.MarshalIndent(response.Servers, "", "    ")
		fmt.Println(string(date))
	} else {
		fmt.Println(err)
	}
}
func (h *HClient) ListIps() error {

	auth := basic.NewCredentialsBuilder().
		WithAk(h.Ak).
		WithSk(h.Sk).
		WithProjectId("").
		Build()

	client := eip.NewEipClient(
		eip.EipClientBuilder().
			WithEndpoint(h.VpcEndpoint).
			WithCredential(auth).
			Build())

	request := &emodel.NeutronListFloatingIpsRequest{}

	response, err := client.NeutronListFloatingIps(request)

	if err == nil {
		date, _ := json.MarshalIndent(response.Floatingips, "", "    ")
		fmt.Println(string(date))
	}
	return err
}

func (h *HClient) DeleteIp(ipId string) error {
	auth := basic.NewCredentialsBuilder().
		WithAk(h.Ak).
		WithSk(h.Sk).
		WithProjectId("").
		Build()

	client := eip.NewEipClient(
		eip.EipClientBuilder().
			WithEndpoint(h.VpcEndpoint).
			WithCredential(auth).
			Build())

	request := &emodel.NeutronDeleteFloatingIpRequest{}
	request.FloatingipId = ipId

	response, err := client.NeutronDeleteFloatingIp(request)

	if err == nil {
		fmt.Printf("%+v\n", response)
	}
	return err
}
