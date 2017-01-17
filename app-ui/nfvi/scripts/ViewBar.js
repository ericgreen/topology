(function (nx) {
	nx.define('ViewBar', nx.ui.Component, {
		properties: {
			'topology': null,
			'topologyContainer': null,
			'viewItems': {
				set: function (value) {
					var items = new nx.data.Dictionary(value);
					if ($urlStack.length > 1) {
						items.setItem("Back", $urlStack[$urlStack.length -2]);
                    }
					this.view('views').set('items', items);
				}
			},
			'exportedData': ''
		},
		view: {
			content: [
				{
                    name: 'views',
                    tag: 'nav',
                    style: {
                        'width': '5%',
                    },
                    props: {
                        'class': 'w3-sidenav w3-gray',
						template: {
							tag: 'a',
							props: {
								'addr': '{value}'
							},
							content: {
								tag: 'span',
								content: '{key}',
							},
							events: {
								'click': '{#openView}'
							}
						},
						items: '{viewItems}'
					}
				}
            ]
		},
		methods: {
			'assignTopology': function (topo) {
				this.topology(topo);
			},
			'assignTopologyContainer': function (topoContainer) {
				this.topologyContainer(topoContainer);
			},
            'openView': function(sender, event) {
                event.preventDefault();
				var topo = this.topology();
				topo.clear()
				var topoContainer = this.topologyContainer();
				url = ($(event.srcElement.parentElement.outerHTML).attr("addr"));
				text = $(event.srcElement.parentElement.innerHTML).text();
				if (url == undefined) {
                    url = $(event.srcElement.outerHTML).attr("addr");
					text = $(event.srcElement.outerHTML).text();
                }
				if (text == 'Back') {
                    $urlStack.pop();
                }
				topoContainer.loadTopology(url);
            },
		}
	});
})(nx);