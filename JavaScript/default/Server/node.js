import { Path, Mode } from "../common/utils.js";
import { StatusCode, Header, Method, Parameter } from "../common/protocol.js";
import { createServer, IncomingMessage, OutgoingMessage, request, Server } from "node:http";
import { HandleNodeRegistration } from "./node_communication.js";
import { v4 as uuidv4 } from "uuid";
import readline from "node:readline";
import { ParseBody } from "../common/parse.js";
import { clearInterval } from "node:timers";

export class NodeInfo {
    /** 
     * @param {string} uuid 
     * @param {string} host 
     * @param {number} port 
     */
    constructor(uuid, host, port) {
        let error = [];
        if (typeof(uuid) !== "string") {
            error.push("uuid must be a string.");
        }
        if (typeof(host) !== "string") {
            error.push("host must be a string.");
        }
        if (typeof(port) !== "number") {
            error.push("port must be a number.");
        }
        if (error.length !== 0) {
            let e = error.join(" ");
            throw new TypeError(e);
        }
        /** @type {string} */
        this.Uuid = uuid;

        /** @type {string} */
        this.Host = host;

        /** @type {number} */
        this.Port = parseInt(port, 10);
        if (this.Port < 2_000 || this.Port > 40_000) {
            throw new Error("expected port range between 2_000 and 40_000");
        }
    }
}

export class Node {
    /**
     * @param {string} localHost 
     * @param {string} localPort 
     * @param {string} remoteHost 
     * @param {string} remotePort 
     * @param {boolean} leader 
     */
    constructor(localHost, localPort, remoteHost, remotePort, leader) {
        /** @type {NodeInfo} */
        this.Info = new NodeInfo(uuidv4().replace(/-/g, ""), localHost, localPort);

        /** @type {NodeInfo} */
        this.Leader = new NodeInfo(leader ? this.Uuid : "", remoteHost, remotePort);

        /** @type {Mode} */
        this.Type = leader ? Mode.LeaderNode : Mode.ServerNode;

        /** @type {string[]} */
        this.InputBuffer = [];

        /** @type {readline.Interface} */
        this.ConsoleReader = null;

        /** @type {NodeJS.Timeout} */
        this.InputTask = null;

        /** @type {string[]} */
        this.OutputBuffer = [];

        /** @type {NodeJS.Timeout} */
        this.OutputTask = null;

        /** @type {NodeInfo[]} */
        this.NodeList = [];

        /** @type {string[]} */
        this.Services = [];

        /** @type {Server} */
        this.Server = null;
    }

    /** @param {string} input */
    AddInput = (input) => {
        this.InputBuffer.push(input);
    }

    /** @param {string} output */
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
        }, 100);
    }

    HandleOutput = async () => {
        this.OutputTask = setInterval(() => {
            while (this.OutputBuffer.length > 0) {
                console.log(this.OutputBuffer.shift());
            }
        }, 100);
    }

    Start = () => {
        this.Server = createServer((request, response) => this.HandleConnection(request, response));

        this.HandleInput();
        this.ParseInput();
        this.HandleOutput();

        this.Server.listen(this.Info.Port, this.Info.Host, () => {
            this.AddOutput(`Server running at http://${this.Info.Host}:${this.Info.Port}`);
        });

        // Register with leader
        if (this.Type === Mode.ServerNode) {
            const options = {
                hostname: this.Leader.Host,
                port: this.Leader.Port,
                path: Path.NodeRegistration,
                method: Method.Post,
                headers: {
                    [Header.ContentType]: Header.ApplicationJson,
                },
            };
            console.log(options);

            const registerRequest = request(options, (registerResponse) => {
                try {
                    const body = ParseBody(registerResponse);
                    this.Leader.Uuid = body[Parameter.Uuid];
                    console.log(body);
                } catch (error) {
                    console.log(error);
                }
            });

            registerRequest.on(Parameter.Error, (e) => {
                this.AddOutput(`Problem with request: ${e.message}`);
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

    /**
     * @param {IncomingMessage} request 
     * @param {OutgoingMessage} response 
     */
    HandleConnection = (request, response) => {
        const url = new URL(request.url, `http://${request.headers.host}`);

        if (url.pathname.includes(Path.NodeRegistration)) {
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
                    path: `${Path.NodeRegistration}/${this.Info.Uuid}`,
                    method: Method.Delete, 
                    headers: {
                        [Header.ContentType]: Header.ApplicationJson,
                    },
                };

                const deregRequest = request(options, (deregResponse) => {
                    try {
                        const body = ParseBody(deregResponse);
                        console.log(body);
                        resolve();
                    } catch (error) {
                        console.log(error);
                    }
                });

                deregRequest.on(Parameter.Error, (e) => {
                    this.AddOutput(`Problem with request: ${e.message}`);
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
                    this.AddOutput("Closing server");
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

/**
 * @param {string} remote 
 * @param {string} local 
 */
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

