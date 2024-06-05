import { Parameter } from "./protocol.js";

/**
 * @param {import("http").IncomingMessage | import("http").OutgoingMessage} message 
 * @returns {Promise<object | Error>}
 */
export const ParseBody = (message) => {
    return new Promise((resolve, reject) => {
        let body = "";

        message.on(Parameter.Data, (chunk) => {
            body += chunk.toString();
        });

        message.on(Parameter.End, () => {
            try {
                let value = JSON.parse(body);
                resolve(value);
            } catch (error) {
                reject(error);
            }
        });

        message.on(Parameter.Error, (error) => {
            reject(error);
        });

        message.on(Parameter.Close, () => {
            //console.log('Request connection closed');
        });

        message.on(Parameter.Reject, () => {
            reject(new Error("Request was aborted"));
        });
    })
};
