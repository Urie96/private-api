const http = require("http");

http
	.get("http://localhost:8030/openai?Prompt=hello", (response) => {
		let todo = "";

		// called when a data chunk is received.
		response.on("data", (chunk) => {
			console.log(chunk.toString());
		});

		// called when the complete response is received.
		response.on("end", () => {
			console.log("end");
		});
	})
	.on("error", (error) => {
		console.log("Error: " + error.message);
	});
