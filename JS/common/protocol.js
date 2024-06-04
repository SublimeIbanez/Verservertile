

export const Direction = {
    NilDirection: 0,
    Request: 1,
    Response: 2,
};

export const Directive = {
    NilDirective: 0,
    Register: 1,
    Shutdown: 2,
    UpdateNodesList: 3,
    ServicesRequest: 4,
    ServiceOperation: 5,
    ServiceChoice: 6,
}

export class BaseMessage {
    constructor(ok, direction, directive, entityType, uuid, message, data) {
        this.Ok = ok; // True/False
        this.Direction = direction; // Response/Request
        this.Directive = directive; // Intent
        this.entityType = entityType; // Client/Leader/Node
        this.Uuid = uuid; // Unique identifier
        this.Message = message; // Any intended message
        this.Data = data; // Object of values
    }
}

export const NewRequest = (ok, directive, entityType, uuid, message, data) => {
    return new JSON.stringify(BaseMessage(ok, Direction.Request, directive, entityType, uuid, message, data));
}

export const NewResponse = (ok, directive, entityType, uuid, message, data) => {
    return new JSON.stringify(BaseMessage(ok, Direction.Response, directive, entityType, uuid, message, data));
}