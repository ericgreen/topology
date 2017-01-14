(function (nx) {
	nx.define('ViewBar', nx.ui.Component, {
		properties: {
			'topology': null,
			'topologyContainer': null,
			'viewItems': {
				set: function (value) {
					this.view('views').set('items', new nx.data.Dictionary(value));
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
                //url = $(event.srcElement.parentElement.innerHTML).attr('addr')
				topoContainer.loadTopology(url);
            },
		}
	});
})(nx);