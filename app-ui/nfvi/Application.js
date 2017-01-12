(function(nx) {
    var App = nx.define(nx.ui.Application, {
        methods: {
            start: function() {
                var topologyContainer = new TopologyContainer();
                var topology = topologyContainer.topology();
                var actionBar = new ActionBar();
                var viewBar = new ViewBar();
                this.container(document.getElementById('next-app'));
                actionBar.assignTopology(topology);
                actionBar.assignTopologyContainer(topologyContainer);
                viewBar.assignTopology(topology);
                viewBar.assignTopologyContainer(topologyContainer);
                topologyContainer.assignActionBar(actionBar);
                topologyContainer.assignViewBar(viewBar);
                actionBar.attach(this);
                viewBar.attach(this);
                topology.attach(this);
                topologyContainer.loadTopology(document.baseURI + '/topology/cloudsTopology');
            }
        }
    });
    new App().start();
})(nx);
