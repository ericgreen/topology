package main

func discover() error {
	cloudLoadCloudInfo()
	cloudList := cloudGetCloudList()
	for _, cloudInfo := range cloudList {
		hypervisorList := cloudGetHypervisorList(&cloudInfo)
		for _, hypervisor := range hypervisorList {
			libvirtLoadInfo(cloudInfo, hypervisor.HostIP)
			ovsLoadInfo(cloudInfo, hypervisor.HostIP)
		}
	}
	return nil
}
