(function (nx) {
	nx.define('ViewBar', nx.ui.Component, {
		properties: {
			'exportedData': ''
		},
		view: {
			content: [
				{
                    name: 'list',
                    tag: 'nav',
                    props: {
                        'class': 'w3-sidenav w3-gray'
                    },
                    style: {
                        'width': '5%',
                    },
                    content: [
                        {
                            content: [
                                {
                                    tag: 'a',
                                    content: 'Cloud Instances'
                                }
                            ]
                        },
                        {
                            content: [
                                {
                                    tag: 'a',
                                    content: 'Cloud Instances'
                                }
                            ]
                        }
                    ]
                }
            ]
		}
	});
})(nx);