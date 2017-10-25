package main

var htmlPage = `<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
	<title>ligo workbook</title>
	<style type="text/css">
	body {
		width: 100%;
		height: 100%;
		top: 0px;
		left: 0px;
		margin: 0px;
		background-color: #222222;
	}
	.nilReturn {
		color : #555555;
	}
        .error {
                color : #F30000;
        }
	.workspace {
		padding-top : 10px;
		width : 70%;
		height : 90%;
		border-radius: 3px;
		background-color: #888888;
		overflow-y : auto;
	}

	#interp {
		resize: both;
		width : 98%;
		box-shadow : none;
		border : none;
		min-height: 35px;
		font-size: 30px;
		color : #FFFFFF;
		padding: 1%;
		font-family: monospace;
		margin-top : 10px;
		text-align: left;
		background-color: #111111;
	}
	.interpDone {
		border-radius: 3px;
		min-height: 4%;
		font-size: 30px;
		color : #777777;
		padding-left: 10px;
		font-family: monospace;
		margin : 5px;
		text-align: left;
		background-color: #333333;	
	}
	.result {
		border-radius: 3px;
		min-height: 4%;
		font-size: 30px;
		color : #333333;
		padding-left: 10px;
		font-family: monospace;
		margin : 5px;
		text-align: left;
		background-color: #AAAAAA;
	}
</style>
</head>
<body>
	<center>
		<div class="workspace">
			<div id="interp" onkeydown="readInput(this,event);" onblur="this.focus();" contenteditable></div>
		</div>
	</center>
	<script>
		function insertBefore(el, referenceNode) {
	    	referenceNode.parentNode.insertBefore(el, referenceNode);
		}
		function url(s) {
			var l = window.location;
			return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + l.pathname + s;
		}
		class VMWS {
			constructor() {
				this.ws = new WebSocket(url("ws"));
				this.ws.onmessage = this.receive;
			}

			write(subexp) {
				subexp = subexp.trim();
				if (subexp != "") {
					var hist = document.createElement("div");
					hist.className += " interpDone";
					hist.innerHTML = subexp;
					insertBefore(hist, interp);
					this.ws.send(subexp);
				}
			}

			receive(message) {
				var histRes = document.createElement("div");
				histRes.className += " result";
				histRes.innerHTML = message.data;
				insertBefore(histRes, interp);
				console.log(message.data);
			}
		}
		var vm = new VMWS();
		var interp = document.getElementById("interp");
		window.onload = function() {
			interp.focus();
		};
		function readInput(el, e) {
			if (e.keyCode == 13  && e.ctrlKey) {
				vm.write(el.innerText);
				el.innerText = "";
			}
		}
	</script>
</body>
</html>
`
