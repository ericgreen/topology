package main

import (
	"github.com/SpirentOrion/httprouter"
	"github.com/SpirentOrion/luddite"
	"golang.org/x/net/context"
	"net/http"
)

type APIError struct {
	ErrorCode    int    `json:"errorCode,required" description:"Error Code"`
	ErrorMessage string `json:"errorMessage,required" description:"Error Message"`
}

type DirWithFallback struct {
	d        http.Dir
	fallback string
}

func (df *DirWithFallback) Open(name string) (f http.File, err error) {
	f, err = df.d.Open(name)
	if err != nil {
		f, err = df.d.Open(df.fallback)
	}
	return
}

func NewContentHandler(path, fallback string) http.Handler {
	fs := &DirWithFallback{
		d:        http.Dir(path),
		fallback: fallback,
	}
	return http.FileServer(fs)
}

func CloudTopology(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	cloudList := cloudGetCloudList()
	cloudInfo := cloudList[0]

	topologyTitle := cloudInfo.Name + " Topology"

	hypervisorList := cloudGetHypervisorList(cloudInfo)
	networkList := cloudGetNetworkList(cloudInfo)

	var nodeList []TopologyNode
	var linkList []TopologyLink
	var nodeSetList []TopologyNodeSet

	nodeId := 0
	linkId := 0

	cloudProps := make(map[string]interface{})
	cloudViews := make(map[string]string)
	cloudProps["authUrl"] = cloudInfo.AuthUrl
	cloudProps["user"] = cloudInfo.User
	cloudProps["password"] = cloudInfo.Password
	cloudProps["tenant"] = cloudInfo.Tenant
	cloudProps["provider"] = cloudInfo.Provider
	cloudNodeId := nodeId
	cloudNode := TopologyNode{
		ID:         cloudNodeId,
		Name:       cloudInfo.Name,
		DeviceType: "cloud",
		Color:      "#0000FF",
		Props:      cloudProps,
		Views:		cloudViews,
	}
	nodeList = append(nodeList, cloudNode)
	nodeId++

	for i, hypervisor := range hypervisorList {
		hypervisorNodeProps := make(map[string]interface{})
		hypervisorNodeViews := make(map[string]string)
		hypervisorLinkProps := make(map[string]interface{})
		hypervisorNodeProps["id"] = hypervisor.ID
		hypervisorNodeProps["host_name"] = hypervisor.HostName
		hypervisorNodeProps["ip_address"] = hypervisor.HostIP
		hypervisorNodeViews["network topology"] = "http://localhost:9090/topology/cloudInstace/network/" + hypervisor.Name
		hypervisorNodeViews["ovs topology"] = "http://localhost:9090/topology/cloudInstace/ovs/" + hypervisor.Name
		hyperviosrNode := TopologyNode{
			ID: nodeId, Name: hypervisor.Name,
			DeviceType: "host",
			X:          200,
			Y:          (200 + (5 * i)),
			Color:      "#0000FF",
			Props:      hypervisorNodeProps,
			Views:		hypervisorNodeViews,
		}
		nodeList = append(nodeList, hyperviosrNode)
		hypervisorLink := TopologyLink{
			Name:   "",
			Source: cloudNodeId,
			Target: nodeId,
			Color:  "#0000FF",
			Props:  hypervisorLinkProps,
		}
		linkList = append(linkList, hypervisorLink)
		nodeId++
		linkId++
	}

	for i, network := range networkList {
		networkNodeProps := make(map[string]interface{})
		networkNodeViews := make(map[string]string)
		networkLinkProps := make(map[string]interface{})
		networkNodeProps["id"] = network.ID
		networkNode := TopologyNode{
			ID:         nodeId,
			Name:       network.Name,
			DeviceType: "router",
			X:          400,
			Y:          (200 + (5 * i)),
			Color:      "#0000FF",
			Props:      networkNodeProps,
			Views:		networkNodeViews,
		}
		nodeList = append(nodeList, networkNode)
		networkLink := TopologyLink{
			Name:   "",
			Source: cloudNodeId,
			Target: nodeId,
			Color:  "#0000FF",
			Props:  networkLinkProps,
		}
		linkList = append(linkList, networkLink)
		nodeId++
		linkId++
	}

	groupList := make([]TopologyGroup, 0)

	viewList := make(map[string]string, 2)
	viewList["Cloud Instances"] = "http://localhost:9090/topology/cloudInstanceTopology"
	viewList["Cloud Networks"] = "http://localhost:9090/topology/cloudNetworkTopology"

	cloudTopologyData := TopologyData{
		Title:    topologyTitle,
		Nodes:    nodeList,
		Links:    linkList,
		NodeSets: nodeSetList,
		Groups:   groupList,
		Views:    viewList,
	}

	luddite.WriteResponse(rw, http.StatusOK, cloudTopologyData)
}

func CloudInstanceTopology(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	cloudList := cloudGetCloudList()
	cloudInfo := cloudList[0]

	topologyTitle := cloudInfo.Name + " VNF Topology"

	hypervisorList := cloudGetHypervisorList(cloudInfo)

	var nodeList []TopologyNode
	var linkList []TopologyLink
	var nodeSetList []TopologyNodeSet
	nodeId := 0
	linkId := 0

	cloudProps := make(map[string]interface{})
	cloudProps["authUrl"] = cloudInfo.AuthUrl
	cloudProps["user"] = cloudInfo.User
	cloudProps["password"] = cloudInfo.Password
	cloudProps["tenant"] = cloudInfo.Tenant
	cloudProps["provider"] = cloudInfo.Provider
	cloudNodeId := nodeId
	cloudNode := TopologyNode{
		ID:         cloudNodeId,
		Name:       cloudInfo.Name,
		DeviceType: "cloud",
		Color:      "#0000FF",
		Props:      cloudProps,
	}
	nodeList = append(nodeList, cloudNode)
	nodeId++

	var hypervisorNodeSetIdList []int
	for _, hypervisor := range hypervisorList {
		hypervisorNodeProps := make(map[string]interface{})
		hypervisorLinkProps := make(map[string]interface{})
		hypervisorNodeProps["id"] = hypervisor.ID
		hypervisorNodeProps["host_name"] = hypervisor.HostName
		hypervisorNodeProps["ip_address"] = hypervisor.HostIP
		hypervisorNodeId := nodeId
		hypervisorNode := TopologyNode{
			ID:         hypervisorNodeId,
			Name:       hypervisor.Name,
			DeviceType: "host",
			//X:          200,
			//Y:          (200 + (5 * i)),
			Color: "#0000FF",
			Props: hypervisorNodeProps,
		}
		nodeList = append(nodeList, hypervisorNode)
		//hypervisorNodeSetIdList = append(hypervisorNodeSetIdList, hypervisorNodeId)
		hypervisorLink := TopologyLink{
			Name:   "",
			Source: cloudNodeId,
			Target: hypervisorNodeId,
			Color:  "#0000FF",
			Props:  hypervisorLinkProps,
		}
		linkList = append(linkList, hypervisorLink)
		nodeId++
		linkId++
		instanceList := libvirtGetDomainInstances(hypervisor.HostIP)
		var instanceNodeSetIdList []int
		for _, instance := range instanceList {
			instanceNodeProps := make(map[string]interface{})
			instanceLinkProps := make(map[string]interface{})
			instanceNodeProps["uuid"] = instance.UUID
			instanceNodeProps["name"] = instance.Name
			instanceNode := TopologyNode{
				ID:         nodeId,
				Name:       instance.InstanceName,
				DeviceType: "server",
				//X:          300,
				//Y:          (200 + (5 * j)),
				Color: "#0000FF",
				Props: instanceNodeProps,
			}
			nodeList = append(nodeList, instanceNode)
			//instanceNodeSetIdList = append(instanceNodeSetIdList, nodeId)
			hypervisorLink := TopologyLink{
				Name:   "",
				Source: hypervisorNodeId,
				Target: nodeId,
				Color:  "#0000FF",
				Props:  instanceLinkProps,
			}
			linkList = append(linkList, hypervisorLink)
			nodeId++
			linkId++
		}

		nodeSetProps := make(map[string]interface{})
		if len(instanceNodeSetIdList) > 1 {
			nodeSet := TopologyNodeSet{
				ID:         nodeId,
				Nodes:      instanceNodeSetIdList,
				Name:       "instance-group",
				Root:       hypervisorNodeId,
				DeviceType: "groups",
				//X:          500,
				//Y:          100,
				Color: "#0000FF",
				Props: nodeSetProps,
			}
			nodeSetList = append(nodeSetList, nodeSet)
			nodeId++
		}
	}

	if len(hypervisorNodeSetIdList) > 1 {
		nodeSetProps := make(map[string]interface{})
		nodeSet := TopologyNodeSet{
			ID:         nodeId,
			Nodes:      hypervisorNodeSetIdList,
			Name:       "hypervisor-group",
			Root:       cloudNodeId,
			DeviceType: "groups",
			//X:          200,
			//Y:          200,
			Color: "#0000FF",
			Props: nodeSetProps,
		}
		nodeSetList = append(nodeSetList, nodeSet)
		nodeId++
	}

	groupList := make([]TopologyGroup, 0)

	cloudTopologyData := TopologyData{
		Title:    topologyTitle,
		Nodes:    nodeList,
		Links:    linkList,
		NodeSets: nodeSetList,
		Groups:   groupList,
	}

	luddite.WriteResponse(rw, http.StatusOK, cloudTopologyData)
}

func CloudNetworkTopology(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	cloudList := cloudGetCloudList()
	cloudInfo := cloudList[0]

	topologyTitle := cloudInfo.Name + " VNF Network Topology"

	hypervisorList := cloudGetHypervisorList(cloudInfo)
	networkList := cloudGetNetworkList(cloudInfo)

	var nodeList []TopologyNode
	var linkList []TopologyLink
	var nodeSetList []TopologyNodeSet
	nodeId := 0
	linkId := 0

	for i, network := range networkList {
		networkNodeProps := make(map[string]interface{})
		networkNodeViews := make(map[string]string)
		networkNodeProps["id"] = network.ID
		networkNode := TopologyNode{
			ID:         nodeId,
			Name:       network.Name,
			DeviceType: "router",
			X:          400,
			Y:          (200 + (5 * i)),
			Color:      "#0000FF",
			Props:      networkNodeProps,
			Views:		networkNodeViews,
		}
		nodeList = append(nodeList, networkNode)
		nodeId++
	}

	for _, hypervisor := range hypervisorList {
		instanceList := libvirtGetDomainInstances(hypervisor.HostIP)
		for _, instance := range instanceList {
			instanceNodeProps := make(map[string]interface{})
			instanceNodeProps["uuid"] = instance.UUID
			instanceNodeProps["name"] = instance.Name
			instanceNodeId := nodeId
			instanceNode := TopologyNode{
				ID:         nodeId,
				Name:       instance.InstanceName,
				DeviceType: "server",
				//X:          300,
				//Y:          (200 + (5 * j)),
				Color: "#0000FF",
				Props: instanceNodeProps,
			}
			nodeList = append(nodeList, instanceNode)
			nodeId++
			for _, iface := range instance.Interfaces {
				networkLinkProps := make(map[string]interface{})
				for _, node := range nodeList {
					if node.Name == iface.NetworkName {
						networkLink := TopologyLink{
							Name:   "",
							Source: node.ID,
							Target: instanceNodeId,
							Color:  "#0000FF",
							Props:  networkLinkProps,
						}
						linkList = append(linkList, networkLink)
						linkId++
						break
					}
				}
			}
		}
	}

	groupList := make([]TopologyGroup, 0)

	cloudTopologyData := TopologyData{
		Title:    topologyTitle,
		Nodes:    nodeList,
		Links:    linkList,
		NodeSets: nodeSetList,
		Groups:   groupList,
	}

	luddite.WriteResponse(rw, http.StatusOK, cloudTopologyData)
}

func CloudInstanceNetworkTopology(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	cloudList := cloudGetCloudList()
	cloudInfo := cloudList[0]

	topologyTitle := cloudInfo.Name + " VNF Topology"

	hypervisorList := cloudGetHypervisorList(cloudInfo)

	var nodeList []TopologyNode
	var linkList []TopologyLink
	var nodeSetList []TopologyNodeSet
	nodeId := 0
	linkId := 0

	cloudProps := make(map[string]interface{})
	cloudProps["authUrl"] = cloudInfo.AuthUrl
	cloudProps["user"] = cloudInfo.User
	cloudProps["password"] = cloudInfo.Password
	cloudProps["tenant"] = cloudInfo.Tenant
	cloudProps["provider"] = cloudInfo.Provider
	cloudNodeId := nodeId
	cloudNode := TopologyNode{ID: cloudNodeId, Name: cloudInfo.Name, DeviceType: "cloud", Color: "#0000FF", Props: cloudProps}
	nodeList = append(nodeList, cloudNode)
	nodeId++

	var hypervisorNodeSetIdList []int
	for i, hypervisor := range hypervisorList {
		hypervisorNodeProps := make(map[string]interface{})
		hypervisorLinkProps := make(map[string]interface{})
		hypervisorNodeProps["id"] = hypervisor.ID
		hypervisorNodeProps["host_name"] = hypervisor.HostName
		hypervisorNodeProps["ip_address"] = hypervisor.HostIP
		hypervisorNodeId := nodeId
		hypervisorNode := TopologyNode{
			ID:         hypervisorNodeId,
			Name:       hypervisor.Name,
			DeviceType: "host",
			X:          200,
			Y:          (200 + (5 * i)),
			Color:      "#0000FF",
			Props:      hypervisorNodeProps,
		}
		nodeList = append(nodeList, hypervisorNode)
		//hypervisorNodeSetIdList = append(hypervisorNodeSetIdList, hypervisorNodeId)
		hypervisorLink := TopologyLink{
			Name:   "",
			Source: cloudNodeId,
			Target: hypervisorNodeId,
			Color:  "#0000FF",
			Props:  hypervisorLinkProps,
		}
		linkList = append(linkList, hypervisorLink)
		nodeId++
		linkId++
		instanceList := libvirtGetDomainInstances(hypervisor.HostIP)
		var instanceNodeSetIdList []int
		for j, instance := range instanceList {
			instanceNodeProps := make(map[string]interface{})
			instanceLinkProps := make(map[string]interface{})
			instanceNodeProps["uuid"] = instance.UUID
			instanceNodeProps["name"] = instance.Name
			instanceNodeId := nodeId
			instanceNode := TopologyNode{
				ID:         instanceNodeId,
				Name:       instance.InstanceName,
				DeviceType: "server",
				X:          300,
				Y:          (200 + (5 * j)),
				Color:      "#0000FF",
				Props:      instanceNodeProps,
			}
			nodeList = append(nodeList, instanceNode)
			//instanceNodeSetIdList = append(instanceNodeSetIdList, instanceNodeId)
			hypervisorLink := TopologyLink{
				Name:   "",
				Source: hypervisorNodeId,
				Target: instanceNodeId,
				Color:  "#0000FF",
				Props:  instanceLinkProps,
			}
			linkList = append(linkList, hypervisorLink)
			nodeId++
			linkId++
			var bridgeNodeSetIdList []int
			for k, iface := range instance.Interfaces {
				bridgeNodeProps := make(map[string]interface{})
				bridgeLinkProps := make(map[string]interface{})
				bridgeNodeProps["tap"] = iface.DevName
				bridgeNodeProps["mac_address"] = iface.MacAddress
				bridgeNodeProps["network_name"] = iface.NetworkName
				bridgeLinkProps["rx_bytes"] = iface.RxBytes
				bridgeLinkProps["rx_packets"] = iface.RxPackets
				bridgeLinkProps["rx_errs"] = iface.RxErrs
				bridgeLinkProps["rx_drop"] = iface.RxDrop
				bridgeLinkProps["tx_bytes"] = iface.TxBytes
				bridgeLinkProps["tx_packets"] = iface.TxPackets
				bridgeLinkProps["tx_errs"] = iface.TxErrs
				bridgeLinkProps["tx_drop"] = iface.TxDrop
				bridgeNodeId := nodeId
				bridgeNode := TopologyNode{
					ID:         bridgeNodeId,
					Name:       iface.BridgeName,
					DeviceType: "cloud",
					X:          200,
					Y:          (200 + (5 * k)),
					Color:      "#0000FF",
					Props:      bridgeNodeProps,
				}
				nodeList = append(nodeList, bridgeNode)
				//bridgeNodeSetIdList = append(bridgeNodeSetIdList, bridgeNodeId)
				bridgeLink := TopologyLink{
					Name:   "",
					Source: instanceNodeId,
					Target: bridgeNodeId,
					Color:  "#0000FF",
					Props:  bridgeLinkProps,
				}
				linkList = append(linkList, bridgeLink)
				nodeId++
				linkId++
			}
			if len(bridgeNodeSetIdList) > 0 {
				nodeSetProps := make(map[string]interface{})
				nodeSet := TopologyNodeSet{
					ID:         nodeId,
					Nodes:      bridgeNodeSetIdList,
					Name:       "bridge-group",
					Root:       instanceNodeId,
					DeviceType: "groups",
					X:          (200 + (100 * j)),
					Y:          400,
					Color:      "#0000FF",
					Props:      nodeSetProps,
				}
				nodeSetList = append(nodeSetList, nodeSet)
				nodeId++
			}
		}
		if len(instanceNodeSetIdList) > 0 {
			nodeSetProps := make(map[string]interface{})
			nodeSet := TopologyNodeSet{
				ID:         nodeId,
				Nodes:      instanceNodeSetIdList,
				Name:       "instance-group",
				Root:       cloudNodeId,
				DeviceType: "groups",
				X:          500,
				Y:          200,
				Color:      "#0000FF",
				Props:      nodeSetProps,
			}
			nodeSetList = append(nodeSetList, nodeSet)
			nodeId++
		}
	}
	if len(hypervisorNodeSetIdList) > 1 {
		nodeSetProps := make(map[string]interface{})
		nodeSet := TopologyNodeSet{
			ID:         nodeId,
			Nodes:      hypervisorNodeSetIdList,
			Name:       "hypervisor-group",
			DeviceType: "groups",
			X:          200,
			Y:          200,
			Color:      "#0000FF",
			Props:      nodeSetProps,
		}
		nodeSetList = append(nodeSetList, nodeSet)
		nodeId++
	}

	groupList := make([]TopologyGroup, 0)

	cloudTopologyData := TopologyData{
		Title:    topologyTitle,
		Nodes:    nodeList,
		Links:    linkList,
		NodeSets: nodeSetList,
		Groups:   groupList,
	}

	luddite.WriteResponse(rw, http.StatusOK, cloudTopologyData)
}

func CloudInstanceOvsTopology(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	cloudList := cloudGetCloudList()
	cloudInfo := cloudList[0]

	topologyTitle := cloudInfo.Name + " OVS Topology"

	hypervisorList := cloudGetHypervisorList(cloudInfo)

	var nodeList []TopologyNode
	var linkList []TopologyLink
	var nodeSetList []TopologyNodeSet
	nodeId := 0
	linkId := 0

	for _, hypervisor := range hypervisorList {
		bridgeList := ovsGetBridges(hypervisor.HostIP)
		for j, bridge := range bridgeList {
			bridgeNodeProps := make(map[string]interface{})
			bridgeNodeProps["uuid"] = bridge.UUID
			bridgeNodeProps["name"] = bridge.Name
			bridgeNodeId := nodeId
			bridgeNode := TopologyNode{
				ID:         bridgeNodeId,
				Name:       bridge.Name,
				DeviceType: "cloud",
				X:          (-100 + (-300 * j)),
				Y:          (-600 + (300 * j)),
				Color:      "#00FF00",
				Props:      bridgeNodeProps,
			}
			nodeList = append(nodeList, bridgeNode)
			nodeId++
		}
		bridgeConnections := ovsGetBridgeConnections(hypervisor.HostIP)
		for _, bc := range bridgeConnections {
			var sourceBridgeId int
			var targetBridgeId int
			bridgeLinkProps := make(map[string]interface{})
			bridgeLinkProps["source_interface"] = bc.SourceInterface.Name
			bridgeLinkProps["target_interface"] = bc.TargetInterface.Name
			bridgeLinkProps["source_port"] = bc.SourcePort.Name
			bridgeLinkProps["target_port"] = bc.TargetPort.Name
			for k, v := range bc.SourceInterface.Statistics {
				bridgeLinkProps[k] = v
			}
			for _, node := range nodeList {
				if bc.SourceBridge.Name == node.Name {
					sourceBridgeId = node.ID
					break
				}
			}
			for _, node := range nodeList {
				if bc.TargetBridge.Name == node.Name {
					targetBridgeId = node.ID
					break
				}
			}
			bridgeLink := TopologyLink{
				Name:   "",
				Source: sourceBridgeId,
				Target: targetBridgeId,
				Color:  "#00FF00",
				Props:  bridgeLinkProps,
			}
			linkList = append(linkList, bridgeLink)
			linkId++
		}
		instanceList := libvirtGetDomainInstances(hypervisor.HostIP)
		var instanceNodeSetIdList []int
		for j, instance := range instanceList {
			instanceNodeProps := make(map[string]interface{})
			//instanceLinkProps := make(map[string]interface{})
			instanceNodeProps["uuid"] = instance.UUID
			instanceNodeProps["name"] = instance.Name
			instanceNodeId := nodeId
			instanceNode := TopologyNode{
				ID:         instanceNodeId,
				Name:       instance.InstanceName,
				DeviceType: "server",
				X:          800,
				Y:          (-200 + (5 * j)),
				Color:      "#0000FF",
				Props:      instanceNodeProps,
			}
			nodeList = append(nodeList, instanceNode)
			//instanceNodeSetIdList = append(instanceNodeSetIdList, instanceNodeId)
			nodeId++
			var bridgeNodeSetIdList []int
			for k, iface := range instance.Interfaces {
				bridgeNodeProps := make(map[string]interface{})
				bridgeLinkProps := make(map[string]interface{})
				bridgeNodeProps["tap"] = iface.DevName
				bridgeNodeProps["mac_address"] = iface.MacAddress
				bridgeNodeProps["network_name"] = iface.NetworkName
				bridgeLinkProps["rx_bytes"] = iface.RxBytes
				bridgeLinkProps["rx_packets"] = iface.RxPackets
				bridgeLinkProps["rx_errs"] = iface.RxErrs
				bridgeLinkProps["rx_drop"] = iface.RxDrop
				bridgeLinkProps["tx_bytes"] = iface.TxBytes
				bridgeLinkProps["tx_packets"] = iface.TxPackets
				bridgeLinkProps["tx_errs"] = iface.TxErrs
				bridgeLinkProps["tx_drop"] = iface.TxDrop
				bridgeNodeId := nodeId
				bridgeNode := TopologyNode{
					ID:         bridgeNodeId,
					Name:       iface.BridgeName,
					DeviceType: "cloud",
					X:          700,
					Y:          (-500 + (5 * k)),
					Color:      "#FF00FF",
					Props:      bridgeNodeProps,
				}
				nodeList = append(nodeList, bridgeNode)
				//bridgeNodeSetIdList = append(bridgeNodeSetIdList, bridgeNodeId)
				bridgeLink := TopologyLink{
					Name:   "",
					Source: instanceNodeId,
					Target: bridgeNodeId,
					Color:  "#0000FF",
					Props:  bridgeLinkProps,
				}
				linkList = append(linkList, bridgeLink)
				nodeId++
				linkId++
				bc := ovsGetBridgeConnection(hypervisor.HostIP, iface.MacAddress)
				var targetBridgeId int
				ovsBridgeLinkProps := make(map[string]interface{})
				ovsBridgeLinkProps["target_interface"] = bc.TargetInterface.Name
				ovsBridgeLinkProps["target_port"] = bc.TargetPort.Name
				for k, v := range bc.TargetInterface.Statistics {
					ovsBridgeLinkProps[k] = v
				}
				for _, node := range nodeList {
					if bc.TargetBridge.Name == node.Name {
						targetBridgeId = node.ID
						break
					}
				}
				ovsBridgeLink := TopologyLink{
					Name:   "",
					Source: bridgeNodeId,
					Target: targetBridgeId,
					Color:  "#FF00FF",
					Props:  ovsBridgeLinkProps,
				}
				linkList = append(linkList, ovsBridgeLink)
				linkId++

			}
			if len(bridgeNodeSetIdList) > 0 {
				nodeSetProps := make(map[string]interface{})
				nodeSet := TopologyNodeSet{
					ID:         nodeId,
					Nodes:      bridgeNodeSetIdList,
					Name:       "bridge-group",
					Root:       instanceNodeId,
					DeviceType: "groups",
					X:          (-200 + (100 * j)),
					Y:          400,
					Color:      "#0000FF",
					Props:      nodeSetProps,
				}
				nodeSetList = append(nodeSetList, nodeSet)
				nodeId++
			}
		}
		if len(instanceNodeSetIdList) > 0 {
			nodeSetProps := make(map[string]interface{})
			nodeSet := TopologyNodeSet{
				ID:         nodeId,
				Nodes:      instanceNodeSetIdList,
				Name:       "instance-group",
				Root:       0,
				DeviceType: "groups",
				X:          500,
				Y:          200,
				Color:      "#0000FF",
				Props:      nodeSetProps,
			}
			nodeSetList = append(nodeSetList, nodeSet)
			nodeId++
		}
	}

	groupList := make([]TopologyGroup, 0)

	cloudTopologyData := TopologyData{
		Title:    topologyTitle,
		Nodes:    nodeList,
		Links:    linkList,
		NodeSets: nodeSetList,
		Groups:   groupList,
	}

	luddite.WriteResponse(rw, http.StatusOK, cloudTopologyData)
}

func InitApp(router *httprouter.Router) {
	content := NewContentHandler(cfg.App.Path, cfg.App.Fallback)

	router.NotFound = func(_ context.Context, rw http.ResponseWriter, r *http.Request) {
		delete(rw.Header(), luddite.HeaderContentType)
		content.ServeHTTP(rw, r)
	}
	router.GET("/topology/cloudTopology", CloudTopology)
	router.GET("/topology/cloudInstanceTopology", CloudInstanceTopology)
	router.GET("/topology/cloudNetworkTopology", CloudNetworkTopology)
	router.GET("/topology/cloudInstanceNetworkTopology", CloudInstanceNetworkTopology)
	router.GET("/topology/cloudInstanceOvsTopology", CloudInstanceOvsTopology)
}
