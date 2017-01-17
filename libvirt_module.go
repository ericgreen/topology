package main

import (
	"fmt"
	log "github.com/SpirentOrion/logrus"
	"github.com/beevik/etree"
	libvirt "github.com/rgbkrk/libvirt-go"
	"strconv"
)

type LibvirtDomainInstance struct {
	UUID           string
	Name           string
	InstanceName   string
	HypervisorName string
	Interfaces     []LibvirtDomainInterface
}

type LibvirtDomainInterface struct {
	DevName     string
	BridgeName  string
	MacAddress  string
	NetworkName string
	RxBytes     int64
	RxPackets   int64
	RxErrs      int64
	RxDrop      int64
	TxBytes     int64
	TxPackets   int64
	TxErrs      int64
	TxDrop      int64
}

type LibvirtPhysicalInterface struct {
	Name       string
	MacAddress string
}

type LibvirtConnection struct {
	cloudInfo  CloudInfo
	ipAddress  string
	connection libvirt.VirConnection
}

var libvirtDomainInstances map[string][]LibvirtDomainInstance = make(map[string][]LibvirtDomainInstance)
var libvirtPhysicalInterfaces map[string][]LibvirtPhysicalInterface = make(map[string][]LibvirtPhysicalInterface)

func libvirtConnect(cloudInfo CloudInfo, ipAddress string) (*LibvirtConnection, error) {
	var err error
	connectString := fmt.Sprintf("qemu+tcp://%s/system", ipAddress)

	logFields := log.Fields{
		"Name":          cloudInfo.Name,
		"IpAddress":     ipAddress,
		"ConnectString": connectString,
	}

	service.Logger().WithFields(logFields).Info("Establishing libvirt connection")

	c, err := libvirt.NewVirConnectionReadOnly(connectString)
	if err != nil {
		logFields = log.Fields{
			"Name":          cloudInfo.Name,
			"IpAddress":     ipAddress,
			"ConnectString": connectString,
			"Error":         err.Error(),
		}
		service.Logger().WithFields(logFields).Error("libvirt open error: ", err.Error())
		return nil, err
	}
	libvirtConnection := LibvirtConnection{
		cloudInfo:  cloudInfo,
		ipAddress:  ipAddress,
		connection: c,
	}

	service.Logger().WithFields(logFields).Info("Successfullly established libvirt connection")

	return &libvirtConnection, nil
}

func (c *LibvirtConnection) libvirtDisconnect() {
	logFields := log.Fields{
		"Name":      c.cloudInfo.Name,
		"IpAddress": c.ipAddress,
	}

	service.Logger().WithFields(logFields).Info("Disconnecting libvirt connection")

	c.connection.CloseConnection()

	service.Logger().WithFields(logFields).Info("Successfully disconnected libvirt connection")
}

func (c *LibvirtConnection) libvirtLoadDomainInstances() error {
	logFields := log.Fields{
		"Name":      c.cloudInfo.Name,
		"IpAddress": c.ipAddress,
	}

	service.Logger().WithFields(logFields).Info("Loading libvirt domain instances")

	domains, err := c.connection.ListAllDomains(libvirt.VIR_CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		logFields = log.Fields{
			"Name":      c.cloudInfo.Name,
			"IpAddress": c.ipAddress,
			"Error":     err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error listing all libvirt domains")
		return err
	}
	instances := make([]LibvirtDomainInstance, len(domains))
	for i, domain := range domains {
		instance, err := c.libvirtLoadDomainInfo(domain)
		if err != nil {
			logFields = log.Fields{
				"Name":      c.cloudInfo.Name,
				"IpAddress": c.ipAddress,
				"Error":     err.Error(),
			}
			service.Logger().WithFields(logFields).Error("Error getting libvirt domain instance information")
			return err
		}
		instances[i] = *instance
	}
	libvirtDomainInstances[c.ipAddress] = instances

	service.Logger().WithFields(logFields).Info("Succesfully loaded libvirt domain instances")

	return nil
}

func (c *LibvirtConnection) libvirtLoadDomainInfo(domain libvirt.VirDomain) (*LibvirtDomainInstance, error) {
	uuid, err := domain.GetUUIDString()
	if err != nil {
		logFields := log.Fields{
			"Name":      c.cloudInfo.Name,
			"IpAddress": c.ipAddress,
			"Error":     err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error getting libvirt domain instance UUID")
		return nil, err
	}
	name, err := domain.GetName()
	if err != nil {
		logFields := log.Fields{
			"Name":       c.cloudInfo.Name,
			"IpAddress":  c.ipAddress,
			"DomainUUID": uuid,
			"Error":      err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error getting libvirt domain instance name")
		return nil, err
	}

	interfaces, err := c.libvirtLoadDomainInterfaceInfo(domain)
	if err != nil {
		return nil, err
	}

	instance := LibvirtDomainInstance{
		UUID:           uuid,
		Name:           name,
		InstanceName:   resolveInstanceName(&c.cloudInfo, uuid),
		HypervisorName: resolveHypervisorName(&c.cloudInfo, c.ipAddress),
		Interfaces:     interfaces,
	}

	return &instance, nil
}

func (c *LibvirtConnection) libvirtLoadDomainInterfaceInfo(domain libvirt.VirDomain) ([]LibvirtDomainInterface, error) {
	xmlDesc, _ := domain.GetXMLDesc(0)
	doc := etree.NewDocument()

	if err := doc.ReadFromString(xmlDesc); err != nil {
		logFields := log.Fields{
			"Name":      c.cloudInfo.Name,
			"IpAddress": c.ipAddress,
			"Error":     err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error reading from file getting libvirt interface info")
		return nil, err
	}

	interfaces := make([]LibvirtDomainInterface, len(doc.FindElements(".//devices/interface/target")))
	for i, e := range doc.FindElements(".//devices/interface/target") {
		devName := e.SelectAttr("dev").Value
		se := doc.FindElement(".//devices/interface[" + strconv.Itoa(i+1) + "]/source")
		bridgeName := se.SelectAttr("bridge").Value
		me := doc.FindElement(".//devices/interface[" + strconv.Itoa(i+1) + "]/mac")
		macAddr := me.SelectAttr("address").Value
		stats, err := domain.InterfaceStats(devName)
		if err != nil {
			logFields := log.Fields{
				"Name":      c.cloudInfo.Name,
				"IpAddress": c.ipAddress,
				"DevName":   devName,
				"Error":     err.Error(),
			}
			service.Logger().WithFields(logFields).Error("Error reading libvirt interface stats for device")
			return nil, err
		}
		interfaces[i] = LibvirtDomainInterface{
			DevName:     devName,
			BridgeName:  bridgeName,
			MacAddress:  macAddr,
			NetworkName: resolveNetworkName(&c.cloudInfo, macAddr),
			RxBytes:     stats.RxBytes,
			RxPackets:   stats.RxPackets,
			RxErrs:      stats.RxErrs,
			RxDrop:      stats.RxDrop,
			TxBytes:     stats.TxBytes,
			TxPackets:   stats.TxPackets,
			TxErrs:      stats.TxErrs,
			TxDrop:      stats.TxDrop,
		}
	}
	return interfaces, nil
}

func (c *LibvirtConnection) libvirtLoadPhysicalInterfaces() error {
	logFields := log.Fields{
		"Name":      c.cloudInfo.Name,
		"IpAddress": c.ipAddress,
	}

	service.Logger().WithFields(logFields).Info("Loading libvirt physical interfaces")

	ifaces, err := c.connection.ListAllInterfaces(0)
	if err != nil {
		logFields := log.Fields{
			"Name":      c.cloudInfo.Name,
			"IpAddress": c.ipAddress,
			"Error":     err.Error(),
		}
		service.Logger().WithFields(logFields).Error("Error getting libvirt physical interfaces")
		return err
	}

	physicalInterfaces := make([]LibvirtPhysicalInterface, len(ifaces))
	for i, iface := range ifaces {
		name, err := iface.GetName()
		if err != nil {
			logFields = log.Fields{
				"Name":      c.cloudInfo.Name,
				"IpAddress": c.ipAddress,
				"Error":     err.Error(),
			}
			service.Logger().WithFields(logFields).Error("Error getting libvirt physical interface name")
			return err
		}
		macAddress, err := iface.GetMACString()
		if err != nil {
			logFields = log.Fields{
				"Name":          c.cloudInfo.Name,
				"IpAddress":     c.ipAddress,
				"InterfaceName": name,
				"Error":         err.Error(),
			}
			service.Logger().WithFields(logFields).Error("Error getting libvirt physical interface mac address")
			return err
		}

		pIface := LibvirtPhysicalInterface{
			Name:       name,
			MacAddress: macAddress,
		}

		physicalInterfaces[i] = pIface
	}
	libvirtPhysicalInterfaces[c.ipAddress] = physicalInterfaces

	service.Logger().WithFields(logFields).Info("Succesfully loaded libvirt physical interfaces")

	return nil
}

func libvirtGetDomainInstances(ipAddress string) []LibvirtDomainInstance {
	if dList, ok := libvirtDomainInstances[ipAddress]; ok == true {
		return dList
	}
	var dList []LibvirtDomainInstance
	return dList
}

func libvirtGetDomainInstance(ipAddress string, instanceName string) *LibvirtDomainInstance {
	if dList, ok := libvirtDomainInstances[ipAddress]; ok == true {
		for _, di := range dList {
			if di.InstanceName == instanceName {
				return &di
			}
		}
	}
	return nil
}

func libvirtGetPhysicalInterfaces(ipAddress string) []LibvirtPhysicalInterface {
	if iList, ok := libvirtPhysicalInterfaces[ipAddress]; ok == true {
		return iList
	}
	var iList []LibvirtPhysicalInterface
	return iList
}
