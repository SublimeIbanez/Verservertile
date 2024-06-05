import { Path, Mode } from "../common/utils.js";
import { StatusCode, Header, Method, Parameter } from "../common/protocol.js";
import { createServer, request } from "node:http";
import { HandleNodeRegistration } from "./node_communication.js";
import { v4 as uuidv4 } from "uuid";
import readline from "node:readline";
import { ParseBody } from "../common/parse.js";
import { clearInterval } from "node:timers";

export class NodeInfo {
    constructor(uuid, host, port) {
        this.Uuid = uuid;
        this.Host = host;
        this.Port = parseInt(port);
    }
}

export class Node {
    constructor(localHost, localPort, remoteHost, remotePort, leader) {
        this.Info = new NodeInfo(uuidv4().replace(/-/g, ""), localHost, localPort);
        this.Leader = new NodeInfo(leader ? this.Uuid : "", remoteHost, remotePort);
        this.Type = leader ? Mode.LeaderNode : Mode.ServerNode;
        this.InputBuffer = [];
        this.ConsoleReader = null;
        this.InputTask = null;
        this.OutputBuffer = [];
        this.OutputTask = null;
        this.NodeList = [];
        this.Services = [];
        this.Server = null;
    }

    AddInput = (input) => {
        this.InputBuffer.push(input);
    }

    AddOutput = (output) => {
        this.OutputBuffer.push(output);
    }

    HandleInput = async () => {
        this.ConsoleReader = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: true,
        });

        this.ConsoleReader.on("line", (line) => {
            console.log(`test: ${line}`);
            this.InputBuffer.push(line);
        })
    }

    ParseInput = async () => {
        this.InputTask = setInterval(() => {
            while (this.InputBuffer.length > 0) {
                let input = this.InputBuffer.shift();
                switch (input.toLowerCase()) {
                    case "exit": {
                        this.Shutdown();
                        break;
                    }
                }
            }
        })
    }

    HandleOutput = async () => {
        this.OutputTask = setInterval(() => {
            while (this.OutputBuffer.length > 0) {
                console.log(this.OutputBuffer.shift());
            }
        }, 100)
    }

    Start = () => {
        this.Server = createServer((request, response) => this.HandleConnection(request, response));

        this.HandleInput();
        this.ParseInput();
        this.HandleOutput();

        this.Server.listen(this.Info.Port, this.Info.Host, () => {
            this.OutputBuffer.push(`Server running at http://${this.Info.Host}:${this.Info.Port}`);
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
                this.OutputBuffer.push(`Problem with request: ${e.message}`);
            });

            const registrationData = JSON.stringify({
                [Parameter.Uuid]: this.Info.Uuid,
                [Parameter.Host]: this.Info.Host,
                [Parameter.Port]: this.Info.Port,
            });

            registerRequest.write(registrationData);
            registerRequest.end();
        }
    }

    HandleConnection = (request, response) => {
        const url = new URL(request.url, `http://${request.headers.host}`);

        if (url.pathname.includes(Path.Registration)) {
            HandleNodeRegistration(url, request, response, this);
            return;
        }

        response.statusCode = StatusCode.NotFound;
        response.setHeader(Header.ContentType, Header.TextPlain);
        response.end("Not Found");
    }

    Shutdown = async () => {
        // De-register
        await (new Promise((resolve) => {
            if (this.Type === Mode.ServerNode) {
                const options = {
                    hostname: this.Leader.Host,
                    port: this.Leader.Port,
                    path: `${Path.Registration}/${this.Info.Uuid}`,
                    method: Method.Delete, // /api/something?id=somenumber
                    headers: {
                        [Header.ContentType]: Header.ApplicationJson,
                    },
                };

                const deregRequest = request(options, (deregResponse) => {
                    ParseBody(deregResponse).then((result) => {
                        result.match({
                            Ok: body => console.log(body),
                            Err: error => console.log(error),
                        });
                        resolve();
                    });
                });

                deregRequest.on(Parameter.Error, (e) => {
                    this.OutputBuffer.push(`Problem with request: ${e.message}`);
                    resolve();
                });

                deregRequest.end();
            } else {
                // Set up the next leader
                resolve();
            }
        }));

        if (this.Server) {
            await (new Promise((resolve) => {
                // Close Server
                this.Server.close(() => {
                    this.OutputBuffer.push("Closing server");
                    resolve();
                });
            }));
        }

        // Close the reader
        if (this.ConsoleReader) {
            this.ConsoleReader.close();
        }

        // Close the tasks
        if (this.InputTask) {
            clearInterval(this.InputTask);
        }

        if (this.OutputTask) {
            clearInterval(this.OutputTask);
        }

        console.log("Node shutdown");
        process.exit(0);
    }
}

export const HandleServer = (remote, local) => {
    const [localHost, localPort] = local.split(":");
    let node;
    if (remote.length === 0) {
        node = new Node(localHost, localPort, localHost, localPort, true);
    } else {
        const [remoteHost, remotePort] = remote.split(":");
        node = new Node(localHost, localPort, remoteHost, remotePort, false);
    }

    node.Start();
}

