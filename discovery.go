package main

func discover() error {
	if err := cloudLoadCloudInfo(); err != nil {
		return err
	}
	cloudList := cloudGetCloudList()
	for _, cloudInfo := range cloudList {
		if err := cloudLoadHypervisors(cloudInfo); err != nil {
			//return err
			continue
		}
		if err := cloudLoadInstances(cloudInfo); err != nil {
			//return err
			continue
		}
		if err := cloudLoadNetworks(cloudInfo); err != nil {
			//return err
			continue
		}
		if err := cloudLoadNetworkPorts(cloudInfo); err != nil {
			//return err
			continue
		}
		hypervisorList := cloudGetHypervisorList(&cloudInfo)
		for _, hypervisor := range hypervisorList {
			lc, err := libvirtConnect(cloudInfo, hypervisor.HostIP)
			if err != nil {
				//return err
				continue
			}
			if err = lc.libvirtLoadDomainInstances(); err != nil {
				//return err
			}
			if err = lc.libvirtLoadPhysicalInterfaces(); err != nil {
				//return err
			}
			lc.libvirtDisconnect()

			oc, err := ovsConnect(cloudInfo, hypervisor.HostIP, 6640)
			if err != nil {
				oc, err = ovsConnect(cloudInfo, hypervisor.HostIP, 6641)
				if err != nil {
					//return err
					continue
				}
			}
			if err = oc.ovsLoadBridges(); err != nil {
				//return err
				continue
			}
			if err = oc.ovsLoadInterfaces(); err != nil {
				//return err
			}
			if err = oc.ovsLoadPorts(); err != nil {
				//return err
			}
			oc.ovsDisconnect()
		}
	}
	return nil
}
