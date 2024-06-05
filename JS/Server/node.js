import { Parameter, Path, StatusCode, Header, Mode, Method } from "../common/utils.js";
import { createServer, request } from "node:http";
import { HandleNodeRegistration } from "./node_communication.js";
import { v4 as uuidv4 } from "uuid";
import readline from "node:readline";
import { ParseBody } from "../common/parse.js";

export class NodeInfo {
    constructor(uuid, host, port) {
        this.Uuid = uuid;
        this.Host = host;
        this.Port = parseInt(port);
    }
}

export class Server {
    constructor(localHost, localPort, remoteHost, remotePort, leader) {
        this.Uuid = uuidv4().replace(/-/g, "");
        this.Port = parseInt(localPort);
        this.Host = localHost;
        this.Leader = new NodeInfo(leader ? this.Uuid : "", remoteHost, remotePort);
        this.Type = leader ? Mode.LeaderNode : Mode.ServerNode;
        this.InputBuffer = [];
        this.OutputBuffer = [];
        this.NodeList = [];
        this.Services = [];
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
        const server = createServer((request, response) => this.HandleConnection(request, response));

        this.HandleInput();
        this.HandleOutput();

        server.listen(this.Port, this.Host, () => {
            console.log(`Server running at http://${this.Host}:${this.Port}`);
        });

        // Register with leader
        if (this.Type === Mode.ServerNode) {
            const options = {
                hostname: this.Leader.Host,
                port: this.Leader.Port,
                path: Path.Registration,
                method: Method.Post,
                headers: {
                    [Header.ContentType]: Header.ApplicationJson,
                },
            };

            const registerRequest = request(options, (registerResponse) => {
                ParseBody(registerResponse).then((result) => {
                    result.match({
                        Ok: body => { 
                            this.Leader.Uuid = body[Parameter.Uuid];
                            console.log(body);
                        },
                        Err: error => { 
                            console.log(error);
                        }
                    });
                });
            });

            registerRequest.on(Parameter.Error, (e) => {
                console.error(`Problem with request: ${e.message}`);
            });

            const registrationData = JSON.stringify({
                [Parameter.Uuid]: this.Uuid,
                [Parameter.Host]: this.Host,
                [Parameter.Port]: this.Port,
            });

            console.log(registrationData);

            registerRequest.write(registrationData);
            registerRequest.end();
        }
    }

    HandleConnection = (request, response) => {
        const url = new URL(request.url, `http://${request.headers.host}`);

        switch (url.pathname) {
            case Path.Registration:
                HandleNodeRegistration(request, response, this);
                break;

            default: {
                response.statusCode = StatusCode.NotFound;
                response.setHeader(Header.ContentType, Header.TextPlain);
                response.end("Not Found");
            }
        }
    }
}

export const HandleServer = (remote, local) => {
    const [localHost, localPort] = local.split(":");
    let server;
    if (remote.length === 0) {
        server = new Server(localHost, localPort, localHost, localPort, true);
    } else {
        const [remoteHost, remotePort] = remote.split(":");
        server = new Server(localHost, localPort, remoteHost, remotePort, false);
    }

    server.Start();
}

