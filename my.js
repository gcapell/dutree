$(document).ready(function() {
		$("#tree1").jstree({
		plugins : [ "themes", "json_data"],
		core : {animation: 0},
		"themes" : {
			"icons" : false
		},
		"json_data": {
			"ajax": {
				"url": "/tree",
				"data": function(n) {
					console.log("n", n);
			        // the result is fed to the AJAX request 'data' option
					if (n.attr) {
						console.log("got attr");
						if (n.attr("id")) {
							console.log("got id");
							var retval = {"id": n.attr("id").replace("node-", "")};
							console.log("retval", retval);
							return retval;
						} else {
							console.log("no id");
						}
					} else {
						console.log("no attr");
					}
					console.log("-1");
					return {"id": -1};
			    },
			}
		},
	});
});
