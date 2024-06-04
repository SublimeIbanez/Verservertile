
export const StatusCode = {
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
};

export const Method = {
    Delete: "DELETE",
    Get: "GET",
    Head: "HEAD",
    Post: "POST",
    Put: "PUT",
    Trace: "TRACE",
};

export const Parameter = {
    Data: "data",
    End: "end",
    Error: "error",
}

export const Header = {
    Accept: "Accept",
    ApplicationJson: "application/json",
    ContentType: "Content-Type",
    TextPlain: "text/plain",
};

export const Mode = {
    None: 0,
    Client: 1,
    ServerNode: 2,
    LeaderNode: 3,
};

export const Path = {
    Registration: "node/registration",
}

export class Result {
    constructor(okVal = null, errVal = null) {
        this.Ok = okVal;
        this.Err = errVal;
    }

    Ok(value) {
        return new Result(value, null);
    }

    Err(value) {
        return new Result(null, value);
    }

    isOk() {
        return this.Ok !== null;
    }

    isErr() {
        return this.Err !== null;
    }

    match({ Ok, Err }) {
        return isOk() ? Ok(this.Ok) : Err(this.Err);
    }

    unwrap() {
        return this.Ok;
    }
};
