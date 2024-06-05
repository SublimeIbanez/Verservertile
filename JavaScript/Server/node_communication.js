import { StatusCode, Header, Method, Parameter, Status } from "../common/protocol.js";
import { ParseBody } from "../common/parse.js";
import { NodeInfo } from "./node.js";
import { Path } from "../common/utils.js";

export const HandleNodeRegistration = (url, request, response, node) => {
    switch (request.method) {
        case Method.Post: {
            ParseBody(request).then((result) => {
                result.match({
                    Ok: (body) => {
                        let exists = false;
                        for (let n of node.NodeList) {
                            if (n.Uuid === body[Parameter.Uuid] || (n.Host === body[Parameter.Host] && n.Port === body[Parameter.Port])) {
                                exists = true;
                                break;
                            }
                        }
                        if (!exists) {
                            node.NodeList.push(new NodeInfo(body[Parameter.Uuid], body[Parameter.Host], body[Parameter.Port]));
                            response.writeHead(StatusCode.Created, {
                                [Header.ContentType]: Header.ApplicationJson,
                            });
                            response.end(JSON.stringify({
                                [Parameter.Status]: Status.Success,
                                [Parameter.Message]: "Node added successfully",
                                [Parameter.Data]: {
                                    [Parameter.Uuid]: node.Leader.Uuid
                                }
                            }));
                        } else {
                            response.writeHead(StatusCode.Unauthorized, {
                                [Header.ContentType]: Header.ApplicationJson,
                            });
                            response.end(JSON.stringify({
                                [Parameter.Status]: Status.Unauthorized,
                                [Parameter.Message]: "Node already exists",
                            }));
                        }
                        console.log(node.NodeList);
                    },
                    Err: (error) => {
                        console.log(error);
                        response.writeHead(StatusCode.InternalServerError, {
                            [Header.ContentType]: Header.ApplicationJson,
                        });
                        response.end(JSON.stringify({
                            [Parameter.Status]: Status.Error,
                            [Parameter.Message]: error,
                        }));
                    },
                });
            });
            break;
        }

        case Method.Delete: {
            const uuid = url.pathname.replace(`${Path.Registration}/`, "");

            node.NodeList = node.NodeList.filter((n) => n.Uuid !== uuid);
            console.log("Node List:", node.NodeList);

            response.writeHead(StatusCode.Ok, {
                [Header.ContentType]: Header.ApplicationJson,
            });
            response.end(JSON.stringify({
                [Parameter.Status]: Status.Success,
                [Parameter.Message]: "Node Removed",
            }));
            break;
        }

        default: {
            response.statusCode = StatusCode.NotFound;
            response.setHeader(Header.ContentType, Header.TextPlain);
            response.end("Not Found");
        }
    }
};