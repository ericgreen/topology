package main

type TopologyData struct {
	Title    string            `json:"title,required"`
	Nodes    []TopologyNode    `json:"nodes,required"`
	Links    []TopologyLink    `json:"links,required"`
	NodeSets []TopologyNodeSet `json:"nodeSet,required"`
	Groups   []TopologyGroup   `json:"groups,required"`
	Views    map[string]string `json:"views,required"`
}

type TopologyNode struct {
	ID         int                    `json:"id,required"`
	Name       string                 `json:"name,required"`
	DeviceType string                 `json:"device_type,required"`
	X          int                    `json:"x,required"`
	Y          int                    `json:"y,required"`
	Color      string                 `json:"color,required"`
	Props      map[string]interface{} `json:"props,required"`
	Views      map[string]string      `json:"views,required"`
}

type TopologyNodeSet struct {
	ID         int                    `json:"id,required"`
	Nodes      []int                  `json:"nodes,required"`
	Name       string                 `json:"name,required"`
	Root       int                    `json:"root,required"`
	DeviceType string                 `json:"device_type,required"`
	X          int                    `json:"x,required"`
	Y          int                    `json:"y,required"`
	Color      string                 `json:"color,required"`
	Props      map[string]interface{} `json:"props,required"`
}

type TopologyLink struct {
	Name   string                 `json:"name,required"`
	Source int                    `json:"source,required"`
	Target int                    `json:"target,required"`
	Color  string                 `json:"color,required"`
	Width  int                    `json:"width,required"`
	Props  map[string]interface{} `json:"props,required"`
}

type TopologyGroup struct {
	NodeIDs []int  `json:"node_ids,required"`
	Shape   string `json:"shape,required"`
	Label   string `json:"label,required"`
	Color   string `json:"color,required"`
}
