(function(nx, global) {
    nx.define('TopologyContainerNodeTooltipContent', nx.ui.Component, {
        properties: {
            node: {
                set: function (value) {
                    var model = value.model();
                    this.view('list').set('items', new nx.data.Dictionary(model.getData().props));
                    this.view('views').set('items', new nx.data.Dictionary(model.getData().views));
                    this.title(value.label());
                }
            },
            topology: {},
            title: {}
        },
        view: {
            content: [
                {
                    name: 'header',
                    props: {
                        'class': 'n-topology-tooltip-header'
                    },
                    content: [
                        {
                            tag: 'span',
                            props: {
                                'class': 'n-topology-tooltip-header-text'
                            },
                            name: 'title',
                            content: '{#title}'
                        }
                    ]
                },
                {
                    name: 'content',
                    props: {
                        'class': 'n-topology-tooltip-content n-list'
                    },
                    content: [
                        {
                            name: 'list',
                            tag: 'ul',
                            props: {
                                'class': 'n-list-wrap',
                                template: {
                                    tag: 'li',
                                    props: {
                                        'class': 'n-list-item-i',
                                        role: 'listitem'
                                    },
                                    events: {
                                        'click': '{#`openView`}'
                                    },
                                    content: [
                                        {
                                            tag: 'label',
                                            content: '{key}: ',
                                        },
                                        {
                                            tag: 'span',
                                            content: '{value}',
                                        }
                                    ]

                                }
                            }
                        },
                        {
                            name: 'header',
                            props: {
                                'class': 'n-topology-tooltip-header'
                            },
                            content: [
                                {
                                    tag: 'span',
                                    props: {
                                        'class': 'n-topology-tooltip-header-text'
                                    },
                                    name: 'title',
                                    content: 'views'
                                }
                            ]
                        },
                        {
                            name: 'views',
                            tag: 'ul',
                            props: {
                                'class': 'n-list-wrap views',
                                'id': 'viewlist',
                                template: {
                                    tag: 'li',
                                    props: {
                                        'class': 'n-list-item-i',
                                        role: 'listitem'
                                    },
                                    content:
                                        {
                                            tag: 'a',
                                            props: {
                                                'href': '#',
                                                'addr': '{value}'
                                            },
                                            events: {
                                                'click': '{#openView}'
                                            },                                            
                                            content:
                                                {
                                                    tag: 'span',
                                                    content: '{key}',
                                                }
                                        }
                                }
                            }
                        }
                    ]
                }
            ]
        },
        methods: {
            init: function (args) {
                this.inherited(args);
                this.sets(args);
            },
            openView: function(sender, event) {
                event.preventDefault();
                alert($(event.srcElement.parentElement.parentElement.innerHTML).attr("addr"));
            },
        }
    });
    nx.define('TopologyContainerLinkTooltipContent', nx.graphic.Topology.LinkTooltipContent, {
        properties: {
            link: {
                set: function (value) {
                    var model = value.model();
                    var items = new nx.data.Dictionary(model.getData().props);
                    items.setItem("source", model.getData().source);
                    items.setItem("target", model.getData().target);
                    this.view('list').set('items', items);
                }
            },
            topology: {},
            tooltipmanager: {}
        },
        view: {
            content: {
                props: {
                    'class': 'n-topology-tooltip-content n-list'
                },
                content: [
                    {
                        name: 'list',
                        tag: 'ul',
                        props: {
                            'class': 'n-list-wrap',
                            template: {
                                tag: 'li',
                                props: {
                                    'class': 'n-list-item-i',
                                    role: 'listitem'
                                },
                                content: [
                                    {
                                        tag: 'label',
                                        content: '{key}: '
                                    },
                                    {
                                        tag: 'span',
                                        content: '{value}'
                                    }
                                ]

                            }
                        }
                    }
                ]
            }
        }
    });
    nx.define('TopologyContainer', nx.ui.Component, {
        view: {
            props: {
                'class': "demo-topology-via-api"
            },
            content: [
                {
                    name: 'topology',
                    type: 'nx.graphic.Topology',
                    style: {
                        'margin-left': '25%',
                    },
                    props: {
                        adaptive: true,
                        dataProcessor: 'force',
                        showIcon: true,
                        theme: 'green',
                        identityKey: 'id',
                        data: '{#topologyData}',
                        nodeConfig: {
                            iconType: 'model.device_type',
                            label: 'model.name',
                            color: 'model.color'
                        },
                        nodeSetConfig: {
                            iconType: 'model.device_type',
                            label: 'model.name',
                            color: 'model.color'
                        },
                        linkConfig: {
                            linkType: 'parallel',
                            label: 'model.name',
                            color: 'model.color',
                            width: 'model.width'
                        },
                        tooltipManagerConfig: {
                            nodeTooltipContentClass: 'TopologyContainerNodeTooltipContent',
                            linkTooltipContentClass: 'TopologyContainerLinkTooltipContent'
                        },
                    },
                    events: {
                        'topologyGenerated': '{#updateTopology}'
                    }
                }
            ]
        },
        properties: {
            topologyData: {},
			topology: {
				get: function () {
					return this.view('topology');
				}
			},
            actionBar: null,
        },
        methods: {
            init: function(options) {
                this.inherited(options);
            },
			assignActionBar: function (actionBar) {
				this.actionBar(actionBar);
			},
            loadTopology: function(topologyUrl) {
                $.ajax({
                    url: topologyUrl,
                    success: function(data) {
                        this.topologyData(data);
                    }.bind(this)
                });
            },
            updateTopology: function(sender, event) {
                this.actionBar().topologyTitle(this.topologyData().title);
                this.addGroups(sender, event);
            },
            addGroups: function(sender, event) {
                if (this.topologyData().groups.length > 0) {
                    var groupsLayer = sender.getLayer('groups');
                    for (i = 0; i < this.topologyData().groups.length; i++) {
                        group = this.topologyData().groups[i];
                        var nodes = [];
                        for (j = 0; j < group.node_ids.length; j++) {
                            nodeId = group.node_ids[j];
                            nodes[j] = sender.getNode(nodeId);
                        }
                        var group = groupsLayer.addGroup({
                            nodes: nodes,
                            label: group.label,
                            shapeType: group.shape,
                            color: group.color
                        });
                    }
                }
            }
        }
    });
})(nx, nx.global);
