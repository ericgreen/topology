(function(nx) {
    var App = nx.define(nx.ui.Application, {
        methods: {
            start: function() {
                var topologyContainer = new TopologyContainer();
                var topology = topologyContainer.topology();
                var actionBar = new ActionBar();
                this.container(document.getElementById('next-app'));
                actionBar.assignTopology(topology);
                actionBar.assignTopologyContainer(topologyContainer);
                topologyContainer.assignActionBar(actionBar);
                var viewBar = new ViewBar();
                actionBar.attach(this);
                viewBar.attach(this);
                topology.attach(this);
                actionBar.loadCloudTopology();
            }
        }
    });
    new App().start();
})(nx);
