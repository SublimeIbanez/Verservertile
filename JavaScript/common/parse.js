import { Result } from "./result.js";
import { Parameter } from "./protocol.js";

export const ParseBody = (request) => {
    let body = "";

    return new Promise((resolve) => {
        request.on(Parameter.Data, (chunk) => {
            body += chunk.toString();
        });

        request.on(Parameter.End, () => {
            try {
                let value = JSON.parse(body);
                return resolve(Result.Ok(value));
            } catch (error) {
                return resolve(Result.Err(`Could not parse data: ${error.message}`));
            }
        });

        request.on(Parameter.Error, (error) => {
            return resolve(Result.Err(error));
        });

        request.on('close', () => {
            //console.log('Request connection closed');
        });

        request.on('aborted', () => {
            //console.log('Request connection aborted');
        });
    })
};
