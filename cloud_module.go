package main

import (
	"bufio"
	"encoding/json"
	log "github.com/SpirentOrion/logrus"
	"os/exec"
)

type CloudInfo struct {
	Name     string
	AuthUrl  string
	User     string
	Password string
	Tenant   string
	Provider string
}

type CloudHypervisors struct {
	Hypervisors []CloudHypervisorInfo `json:"hypervisors,required"`
}

type CloudHypervisorInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HostIP   string `json:"host_ip"`
	HostName string `json:"host_name"`
}

type CloudInstances struct {
	Instances []CloudInstanceInfo `json:"instances,required"`
}

type CloudInstanceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CloudNetworks struct {
	Networks []CloudNetworkInfo `json:"networks,required"`
}

type CloudNetworkInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CloudNetworkPorts struct {
	NetworkPorts []CloudNetworkPortInfo `json:"network_ports,required"`
}

type CloudNetworkPortInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MacAddress  string `json:"mac_address,required"`
	NetworkName string `json:"network_name,required"`
}

var clouds []CloudInfo

var cloudHypervisors map[string][]CloudHypervisorInfo = make(map[string][]CloudHypervisorInfo)
var cloudInstances map[string][]CloudInstanceInfo = make(map[string][]CloudInstanceInfo)
var cloudNetworks map[string][]CloudNetworkInfo = make(map[string][]CloudNetworkInfo)
var cloudNetworkPorts map[string][]CloudNetworkPortInfo = make(map[string][]CloudNetworkPortInfo)

var cloudHypervisorNames map[string]map[string]string = make(map[string]map[string]string)
var cloudInstanceNames map[string]map[string]string = make(map[string]map[string]string)
var cloudNetworkNames map[string]map[string]string = make(map[string]map[string]string)

func cloudLoadCloudInfo() error {
	for i, _ := range cfg.CloudProviders.Providers {
		cloudInfo := CloudInfo{
			Name:     cfg.CloudProviders.Providers[i].Name,
			AuthUrl:  cfg.CloudProviders.Providers[i].AuthUrl,
			User:     cfg.CloudProviders.Providers[i].User,
			Password: cfg.CloudProviders.Providers[i].Password,
			Tenant:   cfg.CloudProviders.Providers[i].Tenant,
			Provider: cfg.CloudProviders.Providers[i].Provider,
		}
		clouds = append(clouds, cloudInfo)

		logFields := log.Fields{
			"Name":     cloudInfo.Name,
			"AuthUrl":  cloudInfo.AuthUrl,
			"User":     cloudInfo.User,
			"Password": cloudInfo.Password,
			"Tenant":   cloudInfo.Tenant,
			"Provider": cloudInfo.Provider,
		}
		service.Logger().WithFields(logFields).Info("Loaded Cloud Parameters")
	}
	return nil
}

func cloudLoadHypervisors(cloudInfo CloudInfo) error {
	cmd := exec.Command("./glimpse",
		"-auth-url", cloudInfo.AuthUrl,
		"-user", cloudInfo.User,
		"-pass", cloudInfo.Password,
		"-tenant", cloudInfo.Tenant,
		"-provider", cloudInfo.Provider,
		"list", "hypervisors")

	logFields := log.Fields{
		"Name": cloudInfo.Name,
	}

	service.Logger().WithFields(logFields).Info("Loading hypervisor list")

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error creating StdoutPipe for glimpse to list hypervisors")
		return err
	}
	// read the data from stdout
	buf := bufio.NewReader(cmdReader)

	err = cmd.Start()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error starting glimpse to list hypervisors")
		return err
	}

	output, _ := buf.ReadString('\n')

	err = cmd.Wait()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error returned from glimpse to list hypervisors")
		return err
	}

	var hypervisors CloudHypervisors
	json.Unmarshal([]byte(output), &hypervisors)

	cloudHypervisors[cloudInfo.Name] = hypervisors.Hypervisors

	hypervisorNames := make(map[string]string, len(hypervisors.Hypervisors))
	for _, hypervisor := range hypervisors.Hypervisors {
		hypervisorNames[hypervisor.HostIP] = hypervisor.Name
	}

	cloudHypervisorNames[cloudInfo.Name] = hypervisorNames

	service.Logger().WithFields(logFields).Info("Sucessfully loaded hypervisor list")

	return nil
}

func cloudLoadInstances(cloudInfo CloudInfo) error {
	cmd := exec.Command("./glimpse",
		"-auth-url", cloudInfo.AuthUrl,
		"-user", cloudInfo.User,
		"-pass", cloudInfo.Password,
		"-tenant", cloudInfo.Tenant,
		"-provider", cloudInfo.Provider,
		"list", "instances")

	logFields := log.Fields{
		"Name": cloudInfo.Name,
	}

	service.Logger().WithFields(logFields).Info("Loading instance list")

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error creating StdoutPipe for glimpse to list instances")
		return err
	}
	// read the data from stdout
	buf := bufio.NewReader(cmdReader)

	err = cmd.Start()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error starting glimpse to list instances")
		return err
	}

	output, _ := buf.ReadString('\n')

	err = cmd.Wait()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error returned from glimpse to list instances")
		return err
	}

	var instances CloudInstances
	json.Unmarshal([]byte(output), &instances)

	_, ok := cloudInstances[cloudInfo.Name]
	if ok {
		cloudInstances[cloudInfo.Name] = append(cloudInstances[cloudInfo.Name], instances.Instances...)
	} else {
		cloudInstances[cloudInfo.Name] = instances.Instances
	}

	instanceNames := make(map[string]string, len(instances.Instances))
	for _, instance := range instances.Instances {
		instanceNames[instance.ID] = instance.Name
	}

	_, ok = cloudInstanceNames[cloudInfo.Name]
	if ok {
		for k, v := range instanceNames {
			cloudInstanceNames[cloudInfo.Name][k] = v
		}
	} else {
		cloudInstanceNames[cloudInfo.Name] = instanceNames
	}

	service.Logger().WithFields(logFields).Info("Sucessfully loaded instance list")

	return nil
}

func cloudLoadNetworks(cloudInfo CloudInfo) error {
	cmd := exec.Command("./glimpse",
		"-auth-url", cloudInfo.AuthUrl,
		"-user", cloudInfo.User,
		"-pass", cloudInfo.Password,
		"-tenant", cloudInfo.Tenant,
		"-provider", cloudInfo.Provider,
		"list", "networks")

	logFields := log.Fields{
		"Name": cloudInfo.Name,
	}

	service.Logger().WithFields(logFields).Info("Loading network list")

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error creating StdoutPipe for glimpse to list networks")
		return err
	}
	// read the data from stdout
	buf := bufio.NewReader(cmdReader)

	if err = cmd.Start(); err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error starting glimpse to list networks")
		return err
	}

	output, _ := buf.ReadString('\n')
	if err = cmd.Wait(); err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error returned from glimpse to list networks")
		return err
	}

	var networks CloudNetworks
	json.Unmarshal([]byte(output), &networks)

	cloudNetworks[cloudInfo.Name] = networks.Networks

	service.Logger().WithFields(logFields).Info("Sucessfully loaded network list")

	return nil
}

func cloudLoadNetworkPorts(cloudInfo CloudInfo) error {
	cmd := exec.Command("./glimpse",
		"-auth-url", cloudInfo.AuthUrl,
		"-user", cloudInfo.User,
		"-pass", cloudInfo.Password,
		"-tenant", cloudInfo.Tenant,
		"-provider", cloudInfo.Provider,
		"list", "network-ports")

	logFields := log.Fields{
		"Name": cloudInfo.Name,
	}

	service.Logger().WithFields(logFields).Info("Loading network port list")

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error creating StdoutPipe for glimpse to list network-ports")
		return err
	}
	// read the data from stdout
	buf := bufio.NewReader(cmdReader)

	if err = cmd.Start(); err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error starting glimpse to list network-ports")
		return err
	}

	output, _ := buf.ReadString('\n')
	if err = cmd.Wait(); err != nil {
		logFields = log.Fields{
			"Name":  cloudInfo.Name,
			"Error": err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error returned from glimpse to list network-ports")
		return err
	}

	var networkPorts CloudNetworkPorts
	json.Unmarshal([]byte(output), &networkPorts)

	cloudNetworkPorts[cloudInfo.Name] = networkPorts.NetworkPorts

	networkNames := make(map[string]string, len(networkPorts.NetworkPorts))
	for _, networkPort := range networkPorts.NetworkPorts {
		networkNames[networkPort.MacAddress] = networkPort.NetworkName
	}

	cloudNetworkNames[cloudInfo.Name] = networkNames

	service.Logger().WithFields(logFields).Info("Sucessfully loaded network port list")

	return nil
}

func cloudGetCloudList() []CloudInfo {
	var cloudList []CloudInfo
	for _, cloudInfo := range clouds {
		found := false
		for i, _ := range cloudList {
			if cloudList[i].Name == cloudInfo.Name {
				found = true
				break
			}
		}
		if !found {
			cloudList = append(cloudList, cloudInfo)
		}
	}
	return cloudList
}

func cloudGetCloudInfo(cloudName string) *CloudInfo {
	for _, cloudInfo := range clouds {
		if cloudInfo.Name == cloudName {
			return &cloudInfo
		}
	}
	return nil
}

func cloudGetHypervisorList(cloudInfo *CloudInfo) []CloudHypervisorInfo {
	if hList, ok := cloudHypervisors[cloudInfo.Name]; ok == true {
		return hList
	}
	var hList []CloudHypervisorInfo
	return hList
}

func cloudGetHypervisorInfo(cloudName string, hypervisorName string) *CloudHypervisorInfo {
	cloudInfo := cloudGetCloudInfo(cloudName)
	if cloudInfo != nil {
		if hList, ok := cloudHypervisors[cloudInfo.Name]; ok == true {
			for _, hypervisorInfo := range hList {
				if hypervisorInfo.Name == hypervisorName {
					return &hypervisorInfo
				}
			}
		}
	}
	return nil
}

func cloudGetIntanceList(cloudInfo *CloudInfo) []CloudInstanceInfo {
	if iList, ok := cloudInstances[cloudInfo.Name]; ok == true {
		return iList
	}
	var iList []CloudInstanceInfo
	return iList
}

func cloudGetInstanceInfo(cloudName string, instanceName string) *CloudInstanceInfo {
	cloudInfo := cloudGetCloudInfo(cloudName)
	if cloudInfo != nil {
		if hList, ok := cloudInstances[cloudInfo.Name]; ok == true {
			for _, instanceInfo := range hList {
				if instanceInfo.Name == instanceName {
					return &instanceInfo
				}
			}
		}
	}
	return nil
}

func cloudGetNetworkList(cloudInfo *CloudInfo) []CloudNetworkInfo {
	if nList, ok := cloudNetworks[cloudInfo.Name]; ok == true {
		return nList
	}
	var nList []CloudNetworkInfo
	return nList
}

func cloudGetNetworkInfo(cloudName string, networkName string) *CloudNetworkInfo {
	cloudInfo := cloudGetCloudInfo(cloudName)
	if cloudInfo != nil {
		if nList, ok := cloudNetworks[cloudInfo.Name]; ok == true {
			for _, networkInfo := range nList {
				if networkInfo.Name == networkName {
					return &networkInfo
				}
			}
		}
	}
	return nil
}

func cloudGetNetworkPortList(cloudInfo *CloudInfo) []CloudNetworkPortInfo {
	if npList, ok := cloudNetworkPorts[cloudInfo.Name]; ok == true {
		return npList
	}
	var npList []CloudNetworkPortInfo
	return npList
}

func resolveHypervisorName(cloudInfo *CloudInfo, ipAddress string) string {
	if hypervisorNames, ok := cloudHypervisorNames[cloudInfo.Name]; ok {
		if hypervisorName, ok := hypervisorNames[ipAddress]; ok {
			return hypervisorName
		}
	}
	return "unknown"
}

func resolveInstanceName(cloudInfo *CloudInfo, instanceID string) string {
	if instanceNames, ok := cloudInstanceNames[cloudInfo.Name]; ok {
		if instanceName, ok := instanceNames[instanceID]; ok {
			return instanceName
		}
	}
	return "unknown"
}

func resolveNetworkName(cloudInfo *CloudInfo, macAddress string) string {
	if networkNames, ok := cloudNetworkNames[cloudInfo.Name]; ok {
		if networkName, ok := networkNames[macAddress]; ok {
			return networkName
		}
	}
	return "unknown"
}
