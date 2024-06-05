import { Result } from "./result.js";
import { Parameter } from "./utils.js";

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
            } catch {
                return resolve(Result.Err("Could not parse data"));
            }
        })

        request.on(Parameter.Error, (error) => {
            return resolve(Result.Err(error));
        })
    })
};