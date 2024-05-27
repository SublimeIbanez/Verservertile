import { createServer } from "node:http";
import { argv } from "node:process";

const remote = "";
const local = "localhost:8000";
const mode = Mode.node;
const currentMode = Arguments.none;

argv.slice(0).forEach(arg => {
    switch (currentMode) {
        case Arguments.modeArg: {
            if (arg.toLowerCase() === "client") {
                mode = Mode.client;
            }
            break;
        }
        case Arguments.localArg: {
            let values = arg.split(":");
            break;
        }
        case Arguments.remoteArg: {

            break;
        }
    }
    switch (Arguments(arg.toLowerCase())) {
        case Arguments.modeArg:
            currentMode = Arguments.modeArg;
            break;
        case Arguments.localArg:
            currentMode = Arguments.localArg;
            break;
        case Arguments.remoteArg:
            currentMode = Arguments.remoteArg;
            break;
        default:
            currentMode = Arguments.none;
            break;
    }
});

const Arguments = {
    modeArg: "-m" | "--mode",
    localArg: "-l" | "--local",
    remoteArg: "-r" | "--remote",
    none,
};

const Mode = client | node;

const server = createServer((req, res) => {
    res.statusCode = 200;
    res.setHeader("Content-Type", "text/plain");
    res.end("Hello World");
});

server.listen(port, hostname, () => {
    console.log(`Server running at http://${local}:${port}/`)
});