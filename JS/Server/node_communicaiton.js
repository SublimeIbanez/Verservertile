import { Method } from "../common/utils.js";

export const HandleNodeRegistration = (request, response) => {
    console.log(request);
    switch (request.method) {
        case Method.Post: {

        }

        case Method.Delete: {

        }

        default: {
            response.statusCode = StatusCode.BadRequest;
            response.setHeader(ResponseHeader.ContentType, ResponseHeader.TextPlain);
            response.end("Invalid request");
        }
    }
};