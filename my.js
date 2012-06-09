$(document).ready(function() {
		$("#tree1").jstree({
		plugins : [ "themes", "json_data"],
		core : {animation: 0},
		"json_data": {
			"ajax": {
				"url": function(n) {
					console.log("url",n);
					return "/tree";
				},
				"data": function(n) {
			        // the result is fed to the AJAX request 'data' option
			        var reply = {
						"id": "fish"
			            // "id": n.attr ? n.attr("id").replace("node-", "") : 1,
			        };
		
					console.log("data", reply, "n", n);
					return reply;
			    },
			}
		},
	});
});
