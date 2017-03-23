package main

import (
	"github.com/SpirentOrion/httprouter"
	"github.com/SpirentOrion/luddite"
	"golang.org/x/net/context"
	"net/http"
)

func Discover(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	service.Logger().Infof("Discovery started")

	cloudLoadCloudInfo()
	cloudList := cloudGetCloudList()
	for _, cloudInfo := range cloudList {
		hypervisorList := cloudGetHypervisorList(&cloudInfo)
		for _, hypervisor := range hypervisorList {
			if hypervisor.State == "up" {
				libvirtLoadInfo(cloudInfo, hypervisor.HostIP)
				ovsLoadInfo(cloudInfo, hypervisor.HostIP)
			}
		}
	}

	service.Logger().Infof("Discovery completed")

	apiError := APIError{http.StatusOK, "OK"}
	luddite.WriteResponse(rw, apiError.ErrorCode, apiError)
}

func DiscoverCloud(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	cloudName := httprouter.ContextParams(ctx).ByName("cloud_name")
	cloudInfo := cloudGetCloudInfo(cloudName)
	if cloudInfo == nil {
		apiError := APIError{http.StatusNotFound, "Cloud " + cloudName + " Not discovered"}
		luddite.WriteResponse(rw, apiError.ErrorCode, apiError)
		return
	}
	service.Logger().Infof("Discovery started for cloud %s", cloudName)

	hypervisorList := cloudGetHypervisorList(cloudInfo)
	for _, hypervisor := range hypervisorList {
		if hypervisor.State == "up" {
			libvirtLoadInfo(*cloudInfo, hypervisor.HostIP)
			ovsLoadInfo(*cloudInfo, hypervisor.HostIP)
		}
	}
	service.Logger().Infof("Discovery completed for cloud %s", cloudName)

	apiError := APIError{http.StatusOK, "OK"}
	luddite.WriteResponse(rw, apiError.ErrorCode, apiError)

}

func GetHypervisors(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	cloudName := httprouter.ContextParams(ctx).ByName("cloud_name")
	cloudInfo := cloudGetCloudInfo(cloudName)
	if cloudInfo == nil {
		apiError := APIError{http.StatusNotFound, "Cloud " + cloudName + " Not discovered"}
		luddite.WriteResponse(rw, apiError.ErrorCode, apiError)
		return
	}
	hypervisorList := cloudGetHypervisorList(cloudInfo)

	hypervisors := CloudHypervisors{hypervisorList}
	luddite.WriteResponse(rw, http.StatusOK, hypervisors)
}

func GetInstancesForHypervisor(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	cloudName := httprouter.ContextParams(ctx).ByName("cloud_name")
	cloudInfo := cloudGetCloudInfo(cloudName)
	if cloudInfo == nil {
		apiError := APIError{http.StatusNotFound, "Cloud " + cloudName + " Not discovered"}
		luddite.WriteResponse(rw, apiError.ErrorCode, apiError)
		return
	}

	hostName := httprouter.ContextParams(ctx).ByName("host_name")
	hypervisorInfo := cloudGetHypervisorInfoByHostName(cloudName, hostName)
	if hypervisorInfo == nil {
		apiError := APIError{http.StatusNotFound, "Hypervisor " + hostName + " for cloud "+ cloudName + " Not discovered"}
		luddite.WriteResponse(rw, apiError.ErrorCode, apiError)
		return
	}

	intanceList := cloudGetIntanceListForHypervisor(cloudInfo, hostName)

	instances := CloudInstances{intanceList}
	luddite.WriteResponse(rw, http.StatusOK, instances)
}

func InitDiscovery(router *httprouter.Router) {
	router.POST("/topology/discovery/discover", Discover)
	router.POST("/topology/discovery/discover/:cloud_name", DiscoverCloud)
	router.GET("/topology/discovery/hypervisors/:cloud_name", GetHypervisors)
	router.GET("/topology/discovery/instances/:cloud_name/:host_name", GetInstancesForHypervisor)
}

