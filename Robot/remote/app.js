// to handle incoming HTTP requests,
// use the standard node.js http.Server API
var server = http.createServer(function(req, res) {
    if (req.method != "GET") {
        res.writeHead(405, {"Content-Type": "text/plain"})
        return res.end("Method not allowed")
    }

    var pathname = req.url.split("?")[0]

    if (pathname != "/hello") {
        res.writeHead(404, {"Content-Type": "text/plain"})
        return res.end("Not found.")
    }

    res.writeHead(200, {"Content-Type": "text/plain"})
    res.end("Hello, world!")
})

// establish a WebSocket (or compatible) connection,
// in this case using engine.io
var ws = new eio.Socket({host: "myserver.com"})

// bind the browserver HTTP server to the WebSocket
// and wait for connections from the browserver proxy!
server.listen(ws)

// to make outgoing HTTP requests w/o cross-domain issues,
// use http.get or http.request
http.get("http://www.google.com/index.html", function(res) {
    console.log("Google answered back!")
})