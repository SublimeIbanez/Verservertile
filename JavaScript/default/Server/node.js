import { Assert, Path, Mode, Type } from "../common/utils.js";
import { StatusCode, Header, Method, Parameter, Status } from "../common/protocol.js";
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
        Assert(typeof (host) === Type.String, "host must be a string.");
        Assert(typeof (port) === Type.Number || typeof (port) === Type.String, "port must be a number or a string");

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
     * Interval to manage input/output parsing
     * @type {number} */
    #updateInterval = 100;

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
        this.Leader = new NodeInfo(leader ? this.Info.Uuid : "", remoteHost, remotePort);

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
                    case "leader": {
                        this.AddOutput(this.Leader);
                    }
                }
            }
        }, this.#updateInterval);
    }

    HandleOutput = async () => {
        this.OutputTask = setInterval(() => {
            while (this.OutputBuffer.length > 0) {
                console.log(this.OutputBuffer.shift());
            }
        }, this.#updateInterval);
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

            const registerRequest = request(options, async (registerResponse) => {
                try {
                    const body = await ParseBody(registerResponse);
                    if (body[Parameter.Status] !== Status.Success) {
                        console.error("Unable to connect to leader");
                        this.Shutdown();
                    }
                    
                    if (body[Parameter.Data][Parameter.Uuid] === undefined) {
                        console.error("Invalid response from Leader:", body);
                        this.Shutdown();
                    }
                    this.Leader.Uuid = body[Parameter.Data][Parameter.Uuid];
                    this.AddOutput(`ATTENTION: Connected to ${this.Leader.Uuid}@${this.Leader.Host}:${this.Leader.Port}`);
                } catch (error) {
                    console.error(`Unable to connect to leader: ${error}`);
                    this.Shutdown();
                }
            });

            registerRequest.on(Parameter.Error, (e) => {
                this.AddOutput(`Could not connect to leader: ${e.message}\n-- Shutting down`);
                this.Shutdown();
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
                        console.log(Parameter.Status);
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

