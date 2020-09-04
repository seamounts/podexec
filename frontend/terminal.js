function getQueryVariable(variable) {
	let query = window.location.search.substring(1);
	let vars = query.split("&");
	for (let i=0;i<vars.length;i++) {
			let pair = vars[i].split("=");
			if(pair[0] == variable){return pair[1];}
	}
	return(false);
}

function connect(){
	clusterID=getQueryVariable("cluster")

	if (clusterID == false) {
		alert("无法获取到集群，请联系管理员")
		return
	}
	console.log(clusterID)
	url = "ws://"+document.location.host+ "/webshell/" + clusterID
	console.log(url);
	
	// function findDimensions() {
	// 	height = window.innerHeight
	// 	width = window.innerWidth

	// 	document.getElementById('terminal').style.height = height;
	// 	document.getElementById('terminal').style.width = width;
	// 	console.log(width, height)
	// }

	// findDimensions();

	document.getElementById('terminal').style.height = window.innerHeight

	
	let term = new Terminal({
		"cursorBlink":true,
	});
	if (window["WebSocket"]) {
		
		term.open(document.getElementById("terminal"));
		term.write("connecting to cluster "+ clusterID + "...")
		term.fit();

		conn = new WebSocket(url);
		conn.onopen = function(e) {
			// term.write("\r");
			// msg = {operation: "stdin", data: "export TERM=xterm && clear \r"}
			// conn.send(JSON.stringify(msg))
			// term.clear()
		};
		conn.onmessage = function(event) {
			msg = JSON.parse(event.data)
			if (msg.operation === "stdout") {
				term.write(msg.data)
			} else {
				console.log("invalid msg operation: "+msg)
			}
		};
		conn.onclose = function(event) {
			if (event.wasClean) {
				console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
			} else {
				console.log('[close] Connection died');
				term.writeln("")
			}
			term.write('Connection Reset By Peer! Try Refresh.');
		};
		conn.onerror = function(error) {
			console.log('[error] Connection error');
			term.write("error: "+error.message);
			term.destroy();
		};


		term.on('data', function (data) {
			console.log('data xterm=>',data)
			msg = {operation: "stdin", data: data}
			conn.send(JSON.stringify(msg))
		});

		term.on('resize', function (size) {
			console.log('resize', [size.cols, size.rows]);
			msg = {operation: "resize", cols: size.cols, rows: rows}
			conn.send(JSON.stringify(msg))
		});
		
	} else {
		var item = document.getElementById("terminal");
		item.innerHTML = "<h2>Your browser does not support WebSockets.</h2>";
	}
}