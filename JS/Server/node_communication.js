import { StatusCode, Header, Method, Parameter, Status } from "../common/utils.js";
import { ParseBody } from "../common/parse.js";
import { NodeInfo } from "./node.js";

export const HandleNodeRegistration = (request, response, node) => {
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

            console.log("bl0p");
            break;
        }

        default: {
            response.statusCode = StatusCode.NotFound;
            response.setHeader(Header.ContentType, Header.TextPlain);
            response.end("Not Found");
        }
    }
};