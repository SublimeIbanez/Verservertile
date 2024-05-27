import { Mode } from "../common/utils.js";
import { createServer } from "node:http";
import readline from "node:readline";
import { parentPort, workerData } from "worker_threads";

export class Leader {
    constructor(host, port) {
        this.LeaderHost = host;
        this.LeaderPort = parseInt(port);
    }
}


export class Server {
    constructor(localHost, localPort, remoteHost, remotePort, leader) {
        this.Port = parseInt(localPort);
        this.Host = localHost;
        this.Leader = new Leader(remoteHost, remotePort);
        this.Type = leader ? Mode.ServerNode : Mode.LeaderNode;
        this.InputBuffer = [];
        this.OutputBuffer = [];
    }

    AddInput = (input) => {
        this.InputBuffer.push(input);
    }

    AddOutput = (output) => {
        this.OutputBuffer.push(output);
    }

    HandleInput = async () => {
        const read = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: true,
        });

        read.on("line", (line) => {
            console.log(`test: ${line}`);
            this.InputBuffer.push(line);
        })
    }

    HandleOutput = async () => {
        setInterval(() => {
            while (this.OutputBuffer.length > 0) {
                console.log(this.OutputBuffer.shift());
            }
        }, 100)
    }

    Start = () => {
        const server = createServer((request, response) => {
            response.statusCode = 200;
            response.setHeader("Content-Type", "text/plain");
            response.end("Hello World");
        });

        this.HandleInput();
        this.HandleOutput();

        server.listen(this.Port, this.Host, () => {
            console.log(`Server running at http://${this.Host}:${this.Port}`)
        });
    }
}

export const HandleServer = (remote, local) => {
    let server;
    const [localHost, localPort] = local.split(":");
    if (remote === "") {
        server = new Server(localHost, localPort, localHost, localPort, true);
    } else {
        const [remoteHost, remotePort] = remote.split(":");
        server = new Server(localHost, localPort, remoteHost, remotePort, false);
    }

    server.Start();
}

