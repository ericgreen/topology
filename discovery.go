package main

func discover() error {
	if err := cloudLoadCloudInfo(); err != nil {
		return err
	}
	cloudList := cloudGetCloudList()
	for _, cloudInfo := range cloudList {
		if err := cloudLoadHypervisors(cloudInfo); err != nil {
			return err
		}
		if err := cloudLoadInstances(cloudInfo); err != nil {
			return err
		}
		if err := cloudLoadNetworks(cloudInfo); err != nil {
			return err
		}
		if err := cloudLoadNetworkPorts(cloudInfo); err != nil {
			return err
		}
		hypervisorList := cloudGetHypervisorList(cloudInfo)
		for _, hypervisor := range hypervisorList {
			lc, err := libvirtConnect(cloudInfo, hypervisor.HostIP)
			if err != nil {
				return err
			}
			if err = lc.libvirtLoadDomainInstances(); err != nil {
				return err
			}
			lc.libvirtDisconnect()

			oc, err := ovsConnect(cloudInfo, hypervisor.HostIP)
			if err != nil {
				return err
			}
			if err = oc.ovsLoadBridges(); err != nil {
				return err
			}
			if err = oc.ovsLoadInterfaces(); err != nil {
				return err
			}
			if err = oc.ovsLoadPorts(); err != nil {
				return err
			}
			oc.ovsDisconnect()
		}
	}
	return nil
}
