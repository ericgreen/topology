(function (nx) {
	nx.define('ActionBar', nx.ui.Component, {
		properties: {
			'topology': null,
			'topologyContainer': null,
			'topologyTitle': {},
			'exportedData': ''
		},
		view: {
			content: [
				{
					tag: 'header',
					props: {
						'class': 'w3-container w3-teal',
						style : {
							'text-align': 'center',
						}
					},
					content: [
						{
							tag: 'h1',
							content: 'NFVI Toplogy Viewer',
						},
						{
							tag: 'h5',
							content: '{#topologyTitle}'
						}
					]
				},
				{
					tag: 'ul',
					props: {
						'class': 'w3-navbar w3-blue-grey',
					},
					content: [
						{
							tag: 'li',
							content: [
								{
									tag: 'a',
									props: {
										'href': '#'
									},
									content: 'Clouds'
								}
							],
							events: {
								'click': '{#loadCloudsTopology}'
							}
						},
					]
				},
			]
		},
		methods: {
			'loadCloudsTopology': function (sender, event) {
				event.preventDefault();
				var topo = this.topology();
				topo.clear()
				var topoContainer = this.topologyContainer();
				$urlStack.length = 0;
				var url = document.baseURI + '/topology/cloudsTopology';
				$urlStack.push(url);
				topoContainer.loadTopology(url);
			},
			'assignTopology': function (topo) {
				this.topology(topo);
			},
			'assignTopologyContainer': function (topoContainer) {
				this.topologyContainer(topoContainer);
			}
		}
	});
})(nx);