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
									content: 'Cloud Topology'
								}
							],
							events: {
								'click': '{#loadCloudTopology}'
							}
						},
						{
							tag: 'li',
							content: [
								{
									tag: 'a',
									props: {
										'href': '#'
									},
									content: 'Cloud Instances'
								}
							],
							events: {
								'click': '{#loadCloudInstanceTopology}'
							}
						},
						{
							tag: 'li',
							content: [
								{
									tag: 'a',
									props: {
										'href': '#'
									},
									content: 'Cloud Networks'
								}
							],
							events: {
								'click': '{#loadCloudNetworkTopology}'
							}
						},
						{
							tag: 'li',
							content: [
								{
									tag: 'a',
									props: {
										'href': '#'
									},
									content: 'Cloud Instance Networks'
								}
							],
							events: {
								'click': '{#loadCloudInstanceNetworkTopology}'
							}
						},
						{
							tag: 'li',
							content: [
								{
									tag: 'a',
									props: {
										'href': '#'
									},
									content: 'Cloud Instance OVS Networks'
								}
							],
							events: {
								'click': '{#loadCloudInstanceOvsTopology}'
							}
						},
					]
				},
				/*
				{
					tag: 'h2',
					content: '{#topologyTitle}',
					props: {
						style : {
							'text-align': 'center',
							'margin-left' : '10%'
						}
					}
				},
				*/
			]
		},
		methods: {
			'loadCloudTopology': function (sender, event) {
				event.preventDefault();
				var topo = this.topology();
				topo.clear()
				var topoContainer = this.topologyContainer();
				topoContainer.loadTopology('http://localhost:9090/topology/cloudTopology');
			},
			'loadCloudInstanceTopology': function (sender, event) {
				event.preventDefault();
				var topo = this.topology();
				topo.clear()
				var topoContainer = this.topologyContainer();
				topoContainer.loadTopology('http://localhost:9090/topology/cloudInstanceTopology');
			},
			'loadCloudNetworkTopology': function (sender, event) {
				event.preventDefault();
				var topo = this.topology();
				topo.clear()
				var topoContainer = this.topologyContainer();
				topoContainer.loadTopology('http://localhost:9090/topology/cloudNetworkTopology');
			},
			'loadCloudInstanceNetworkTopology': function (sender, event) {
				event.preventDefault();
				var topo = this.topology();
				topo.clear()
				var topoContainer = this.topologyContainer();
				topoContainer.loadTopology('http://localhost:9090/topology/cloudInstanceNetworkTopology');
			},
			'loadCloudInstanceOvsTopology': function (sender, event) {
				event.preventDefault();
				var topo = this.topology();
				topo.clear()
				var topoContainer = this.topologyContainer();
				topoContainer.loadTopology('http://localhost:9090/topology/cloudInstanceOvsTopology');
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