import { Parameter, Path, StatusCode, Header, Mode, Method } from "../common/utils.js";
import { createServer, request } from "node:http";
import { HandleNodeRegistration } from "./node_communicaiton.js";
import { v4 as uuidv4 } from "uuid";
import readline from "node:readline";

export class Leader {
    constructor(host, port) {
        this.LeaderHost = host;
        this.LeaderPort = parseInt(port);
    }
}

export class Server {
    constructor(localHost, localPort, remoteHost, remotePort, leader) {
        this.Uuid = uuidv4().replace("-", "");
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
        const server = createServer((request, response) => HandleConnection(request, response));

        this.HandleInput();
        this.HandleOutput();

        server.listen(this.Port, this.Host, () => {
            console.log(`Server running at http://${this.Host}:${this.Port}`);
        });

        if (this.Type === Mode.ServerNode) {
            const options = {
                hostname: this.Leader.remoteHost,
                port: this.Leader.remotePort,
                path: Path.Registration,
                method: Method.Post,
                headers: {
                    [Header.ContentType]: Header.ApplicationJson,
                },
            };

            const registerRequest = request(options, (registerResponse) => {
                let data = "";

                registerResponse.on(Parameter.Data, (chunk) => {
                    data += chunk;
                });

                registerResponse.on(Parameter.End, () => {
                    console.log(`Response from leader: ${data}`);
                });
            });

            registerRequest.on(Parameter.Error, (e) => {
                console.error(`Problem with request: ${e.message}`);
            });

            const registrationData = JSON.stringify({
                Uuid: this.Uuid,
                Host: this.Host,
                Port: this.Port,
            });

            registerRequest.write(registrationData);
            registerRequest.end();
        }
    }

    HandleConnection = (request, response) => {
        const url = new URL(request.url, `http://${this.Host}:${this.Port}`);

        switch (url.pathname) {
            case Path.Registration:
                HandleNodeRegistration(request, response);

            default: {
                response.statusCode = StatusCode.BadRequest;
                response.setHeader(Header.ContentType, Header.TextPlain);
                response.end("Invalid request");
            }
        }
    }
}

export const HandleServer = (remote, local) => {
    const [localHost, localPort] = local.split(":");
    let server;
    if (remote === local) {
        server = new Server(localHost, localPort, localHost, localPort, true);
    } else {
        const [remoteHost, remotePort] = remote.split(":");
        server = new Server(localHost, localPort, remoteHost, remotePort, false);
    }

    server.Start();
}

