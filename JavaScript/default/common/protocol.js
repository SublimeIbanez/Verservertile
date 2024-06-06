/**
 * HTTP Status Codes
 * @type {number}
 */
export const StatusCode = Object.freeze({
    Continue: 100,
    SwitchProto: 101,
    Processing: 102,
    Early: 103,

    Ok: 200,
    Created: 201,
    Accepted: 202,
    NoContent: 204,

    MovedPermanently: 301,
    Found: 302,
    SeeOther: 303,
    TemporaryRedirect: 307,
    PermanentRedirect: 308,

    BadRequest: 400,
    Unauthorized: 401,
    Forbidden: 403,
    NotFound: 404,
    MethodNotAllowed: 405,
    NotAcceptable: 406,
    PreconditionFailed: 412,
    UnsupportedMediaType: 415,

    InternalServerError: 500,
    NotImplemented: 501,
    BadGateway: 502,
    ServiceUnavailable: 503,
    GatewayTimeout: 504,
});

/**
 * HTTP headers
 * @type {string}
 */
export const Header = Object.freeze({
    Accept: "Accept",
    ApplicationJson: "application/json",
    ContentType: "Content-Type",
    TextPlain: "text/plain",
});

/**
 * HTTP Methods
 * @type {string}
 */
export const Method = Object.freeze({
    Delete: "DELETE",
    Get: "GET",
    Head: "HEAD",
    Post: "POST",
    Put: "PUT",
    Trace: "TRACE",
});

/**
 * Request/Response parameters
 * @type {string}
 */
export const Parameter = Object.freeze({
    Aborted: "aborted",
    Close: "close",
    Data: "data",
    End: "end",
    Error: "error",
    Uuid: "uuid",
    Host: "host",
    Status: "status",
    Message: "message",
    Port: "port",
});

/**
 * Request/Response status
 * @type {string}
 */
export const Status = Object.freeze({
    Success: "success",
    Unauthorized: "unauthorized",
    Error: "error",
});
