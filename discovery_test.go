package main

import (
	"fmt"
)

func discover() {
	cloudLoadCloudList()
	cloudList := cloudGetCloudList()
	for _, cloudInfo := range cloudList {
		if hypervisorList, err := cloudGetHypervisorList(cloudInfo); err == nil {
			fmt.Println(hypervisorList)
			if networkList, err := cloudGetNetworkList(cloudInfo); err == nil {
				fmt.Println(networkList)
			}

			for _, hypervisor := range hypervisorList {
				if c, err := libvirtConnect(hypervisor.HostIP); err == nil {
					if instances, err := c.libvirtGetDomainInstances(); err == nil {
						fmt.Println(instances)
					}
					c.libvirtDisconnect()
				}
				if c, err := ovsConnect(hypervisor.HostIP); err == nil {
					fmt.Println("Bridge List")
					if bridges, err := c.ovsGetBridges(); err == nil {
						for _, bridge := range bridges {
							fmt.Printf("UUID = %s, name = ***%s***\n", bridge.UUID, bridge.Name)
						}
						fmt.Println(bridges)
					}
					if ports, err := c.ovsGetPorts(); err == nil {
						fmt.Println("Port List")
						for _, port := range ports {
							fmt.Printf("UUID = %s, name = %s, interface UUIDS = %v\n", port.UUID, port.Name, port.InterfaceUUIDs)
						}
						fmt.Println(ports)
					}
					if interfaces, err := c.ovsGetInterfaces(); err == nil {
						fmt.Println("Interface List")
						for _, iface := range interfaces {
							if iface.Type == "" {
								fmt.Printf("UUID = %s, name = %s, type = %s\n", iface.UUID, iface.Name, iface.Type, iface.Options["peer"], iface.ExternalIDs["attached-mac"])
							}
						}
						fmt.Println(interfaces)
					}
					fmt.Println("Bridge Connections")
					bridgeConnections := c.ovsGetBridgeConnections()
					for _, bc := range bridgeConnections {
						fmt.Printf("SI = %s, TI = %s, SP = %s, TP = %s, SB = %s, TB = %s\n", bc.SourceInterface.Name, bc.TargetInterface.Name, bc.SourcePort.Name, bc.TargetPort.Name, bc.SourceBridge.Name, bc.TargetBridge.Name)
					}
					c.ovsDisconnect()
				}
			}
		}
	}
}
