package main

import (
	"encoding/json"
	log "github.com/SpirentOrion/logrus"
	"github.com/socketplane/libovsdb"
)

type OvsBridge struct {
	UUID      string
	Name      string
	PortUUIDs []string
}

type OvsPort struct {
	UUID           string
	Name           string
	InterfaceUUIDs []string
}

type OvsInterface struct {
	UUID            string
	Name            string
	MacAddressInUse string
	Type            string
	Options         map[string]string
	ExternalIDs     map[string]string
	Statistics      map[string]float64
}

type OvsBridgeConnection struct {
	HostIP          string
	SourceInterface OvsInterface
	TargetInterface OvsInterface
	SourcePort      OvsPort
	TargetPort      OvsPort
	SourceBridge    OvsBridge
	TargetBridge    OvsBridge
}

type OvsConnection struct {
	cloudInfo  CloudInfo
	ipAddress  string
	connection *libovsdb.OvsdbClient
}

var ovsBridges map[string][]OvsBridge = make(map[string][]OvsBridge)
var ovsPorts map[string][]OvsPort = make(map[string][]OvsPort)
var ovsInterfaces map[string][]OvsInterface = make(map[string][]OvsInterface)

func ovsConnect(cloudInfo CloudInfo, ipAddress string, port int) (*OvsConnection, error) {
	logFields := log.Fields{
		"Name":      cloudInfo.Name,
		"IpAddress": ipAddress,
		"Port":      port,
	}

	service.Logger().WithFields(logFields).Info("Establishing ovsdb connection")

	c, err := libovsdb.Connect(ipAddress, port)
	if err != nil {
		logFields = log.Fields{
			"Name":      cloudInfo.Name,
			"IpAddress": ipAddress,
			"Port":      port,
			"Error":     err.Error(),
		}
		service.Logger().WithFields(logFields).Info("ovsdb connection failed")
		return nil, err
	}
	ovsConnection := OvsConnection{
		cloudInfo:  cloudInfo,
		ipAddress:  ipAddress,
		connection: c,
	}

	service.Logger().WithFields(logFields).Info("Succesfully established ovsdb connection")

	return &ovsConnection, nil
}

func (c *OvsConnection) ovsDisconnect() {
	logFields := log.Fields{
		"Name":      c.cloudInfo.Name,
		"IpAddress": c.ipAddress,
	}

	service.Logger().WithFields(logFields).Info("Disconnecting ovsdb connection")

	c.connection.Disconnect()

	service.Logger().WithFields(logFields).Info("Successfully disconnected ovsdb connection")
}

func (c *OvsConnection) ovsLoadBridges() error {
	logFields := log.Fields{
		"Name":      c.cloudInfo.Name,
		"IpAddress": c.ipAddress,
	}

	service.Logger().WithFields(logFields).Info("Loading OVS bridge list")

	condition := libovsdb.NewCondition("name", "!=", "")
	selectOp := libovsdb.Operation{
		Op:      "select",
		Table:   "Bridge",
		Where:   []interface{}{condition},
		Columns: []string{},
	}
	operations := []libovsdb.Operation{selectOp}
	reply, err := c.connection.Transact("Open_vSwitch", operations...)

	if err == nil && len(reply) > 0 && len(reply[0].Rows) > 0 {
		bridges := make([]OvsBridge, len(reply[0].Rows))
		for i, row := range reply[0].Rows {
			sl := row["ports"].([]interface{})
			bsliced, _ := json.Marshal(sl)
			var oSet libovsdb.OvsSet
			json.Unmarshal(bsliced, &oSet)
			portUUIDs := make([]string, len(oSet.GoSet))
			for k, v := range oSet.GoSet {
				portUUIDs[k] = v.(libovsdb.UUID).GoUUID
			}
			bridges[i] = OvsBridge{
				UUID:      row["_uuid"].([]interface{})[1].(string),
				Name:      row["name"].(string),
				PortUUIDs: portUUIDs,
			}
		}
		ovsBridges[c.ipAddress] = bridges
	} else {
		logFields = log.Fields{
			"Name":      c.cloudInfo.Name,
			"IpAddress": c.ipAddress,
			"Error":     err.Error(),
		}
		service.Logger().WithFields(logFields).Info("Error loading OVS bridge list")
		return err
	}

	service.Logger().WithFields(logFields).Info("Successfully loaded OVS bridge list")

	return err
}

func (c *OvsConnection) ovsLoadPorts() error {
	logFields := log.Fields{
		"Name":      c.cloudInfo.Name,
		"IpAddress": c.ipAddress,
	}

	service.Logger().WithFields(logFields).Info("Loading OVS port list")

	condition := libovsdb.NewCondition("name", "!=", "")
	selectOp := libovsdb.Operation{
		Op:      "select",
		Table:   "Port",
		Where:   []interface{}{condition},
		Columns: []string{},
	}
	operations := []libovsdb.Operation{selectOp}
	reply, err := c.connection.Transact("Open_vSwitch", operations...)

	if err == nil && len(reply) > 0 && len(reply[0].Rows) > 0 {
		ports := make([]OvsPort, len(reply[0].Rows))
		for i, row := range reply[0].Rows {
			interfaceUUIDs := make([]string, 1)
			uuids := row["interfaces"].([]interface{})
			interfaceUUIDs[0] = uuids[1].(string)
			ports[i] = OvsPort{
				UUID:           row["_uuid"].([]interface{})[1].(string),
				Name:           row["name"].(string),
				InterfaceUUIDs: interfaceUUIDs,
			}
		}
		ovsPorts[c.ipAddress] = ports
	} else {
		logFields = log.Fields{
			"Name":      c.cloudInfo.Name,
			"IpAddress": c.ipAddress,
			"Error":     err.Error(),
		}
		service.Logger().WithFields(logFields).Info("Error loading OVS port list")
		return err
	}

	service.Logger().WithFields(logFields).Info("Successfully loaded OVS port list")

	return nil
}

func (c *OvsConnection) ovsLoadInterfaces() error {
	logFields := log.Fields{
		"Name":      c.cloudInfo.Name,
		"IpAddress": c.ipAddress,
	}

	service.Logger().WithFields(logFields).Info("Loading OVS interface list")

	condition := libovsdb.NewCondition("name", "!=", "")
	selectOp := libovsdb.Operation{
		Op:    "select",
		Table: "Interface",
		Where: []interface{}{condition},
		//Columns: []string{"name", "type", "statistics"},
		Columns: []string{},
	}
	operations := []libovsdb.Operation{selectOp}
	reply, err := c.connection.Transact("Open_vSwitch", operations...)

	if err == nil && len(reply) > 0 && len(reply[0].Rows) > 0 {
		interfaces := make([]OvsInterface, len(reply[0].Rows))
		for i, row := range reply[0].Rows {
			sl := row["options"].([]interface{})
			bsliced, _ := json.Marshal(sl)
			var oMap libovsdb.OvsMap
			json.Unmarshal(bsliced, &oMap)
			options := make(map[string]string, len(oMap.GoMap))
			for k, v := range oMap.GoMap {
				options[k.(string)] = v.(string)
			}
			sl = row["external_ids"].([]interface{})
			bsliced, _ = json.Marshal(sl)
			json.Unmarshal(bsliced, &oMap)
			externalIDs := make(map[string]string, len(oMap.GoMap))
			for k, v := range oMap.GoMap {
				externalIDs[k.(string)] = v.(string)
			}
			sl = row["statistics"].([]interface{})
			bsliced, _ = json.Marshal(sl)
			json.Unmarshal(bsliced, &oMap)
			statistics := make(map[string]float64, len(oMap.GoMap))
			for k, v := range oMap.GoMap {
				statistics[k.(string)] = v.(float64)
			}
			var macInUse string
			switch row["mac_in_use"].(type) {
			case string:
				macInUse = row["mac_in_use"].(string)
			default:
				macInUse = ""

			}
			interfaces[i] = OvsInterface{
				UUID:            row["_uuid"].([]interface{})[1].(string),
				Name:            row["name"].(string),
				MacAddressInUse: macInUse,
				Type:            row["type"].(string),
				Options:         options,
				ExternalIDs:     externalIDs,
				Statistics:      statistics,
			}
		}
		ovsInterfaces[c.ipAddress] = interfaces
	} else {
		service.Logger().Error("Operation Failed due to an error:", err.Error())
		return err
	}

	service.Logger().WithFields(logFields).Info("Successfully loaded OVS interface list")

	return nil
}

func ovsGetBridges(ipAddress string) []OvsBridge {
	if bList, ok := ovsBridges[ipAddress]; ok == true {
		return bList
	}
	var bList []OvsBridge
	return bList
}

func ovsGetPorts(ipAddress string) []OvsPort {
	if pList, ok := ovsPorts[ipAddress]; ok == true {
		return pList
	}
	var pList []OvsPort
	return pList
}

func ovsGetInterfaces(ipAddress string) []OvsInterface {
	if iList, ok := ovsInterfaces[ipAddress]; ok == true {
		return iList
	}
	var iList []OvsInterface
	return iList
}

func ovsGetDPDKInterfaces(ipAddress string) []OvsInterface {
	if iList, ok := ovsInterfaces[ipAddress]; ok == true {
		var pIList []OvsInterface
		for _, iface := range iList {
			if iface.Type == "dpdk" {
				pIList = append(pIList, iface)
			}
		}
		return pIList
	}
	var pIList []OvsInterface
	return pIList
}

func ovsGetBridgeConnections(ipAddress string) []OvsBridgeConnection {
	var bridgeConnections []OvsBridgeConnection
	interfaceList := ovsGetInterfaces(ipAddress)
	portList := ovsGetPorts(ipAddress)
	bridgeList := ovsGetBridges(ipAddress)
	for _, iface := range interfaceList {
		if iface.Type == "patch" {
			bridgeConnection := OvsBridgeConnection{
				HostIP:          ipAddress,
				SourceInterface: iface,
			}
			bridgeConnections = append(bridgeConnections, bridgeConnection)
		}
	}
	for i, bridgeConnection := range bridgeConnections {
		for _, iface := range interfaceList {
			if iface.Name == bridgeConnection.SourceInterface.Options["peer"] {
				bridgeConnections[i].TargetInterface = iface
				break
			}
		}
	}
	for i, bridgeConnection := range bridgeConnections {
		for _, port := range portList {
			if bridgeConnection.SourceInterface.Name == port.Name {
				bridgeConnections[i].SourcePort = port
				break
			}
		}
		for _, port := range portList {
			if bridgeConnection.TargetInterface.Name == port.Name {
				bridgeConnections[i].TargetPort = port
				break
			}
		}
	}
	for i, bridgeConnection := range bridgeConnections {
		for _, bridge := range bridgeList {
			for _, portUUID := range bridge.PortUUIDs {
				if bridgeConnection.SourcePort.UUID == portUUID {
					bridgeConnections[i].SourceBridge = bridge
					break
				}
			}
		}
		for _, bridge := range bridgeList {
			for _, portUUID := range bridge.PortUUIDs {
				if bridgeConnection.TargetPort.UUID == portUUID {
					bridgeConnections[i].TargetBridge = bridge
					break
				}
			}
		}
	}
	return bridgeConnections
}

func ovsGetBridgeConnection(ipAddress string, macAddress string) *OvsBridgeConnection {
	var bridgeConnection OvsBridgeConnection
	interfaceList := ovsGetInterfaces(ipAddress)
	portList := ovsGetPorts(ipAddress)
	bridgeList := ovsGetBridges(ipAddress)
	found := false
	for _, iface := range interfaceList {
		if iface.ExternalIDs["attached-mac"] == macAddress {
			bridgeConnection.HostIP = ipAddress
			bridgeConnection.TargetInterface = iface
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	found = false
	for _, port := range portList {
		if bridgeConnection.TargetInterface.Name == port.Name {
			bridgeConnection.TargetPort = port
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	found = false
	for _, bridge := range bridgeList {
		for _, portUUID := range bridge.PortUUIDs {
			if bridgeConnection.TargetPort.UUID == portUUID {
				bridgeConnection.TargetBridge = bridge
				found = true
				break
			}
		}
	}

	if !found {
		return nil
	}

	return &bridgeConnection
}

func ovsGetPhysicalPortConnection(ipAddress string, name string, macAddress string) *OvsBridgeConnection {
	var bridgeConnection OvsBridgeConnection
	interfaceList := ovsGetInterfaces(ipAddress)
	portList := ovsGetPorts(ipAddress)
	bridgeList := ovsGetBridges(ipAddress)
	found := false
	for _, iface := range interfaceList {
		if (iface.Type == "" || iface.Type == "dpdk" ) && iface.Name == name && iface.MacAddressInUse == macAddress {
			bridgeConnection.HostIP = ipAddress
			bridgeConnection.SourceInterface = iface
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	found = false
	for _, port := range portList {
		if bridgeConnection.SourceInterface.Name == port.Name {
			bridgeConnection.SourcePort = port
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	found = false
	for _, bridge := range bridgeList {
		for _, portUUID := range bridge.PortUUIDs {
			if bridgeConnection.SourcePort.UUID == portUUID {
				bridgeConnection.SourceBridge = bridge
				found = true
				break
			}
		}
	}

	if !found {
		return nil
	}

	return &bridgeConnection
}
