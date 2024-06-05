export class UpdateRequest {
    constructor(leaderId, leaderAddress, nodes, services) {
        this.LeaderId = leaderId;
        this.LeaderAddress = leaderAddress;
        this.Nodes = nodes;
        this.services = services;
    }
}

export class RegistrationRequest {
    constructor(address, services) {
        this.Address = address;
        this.Services = services;
    }
}

export class ServiceListResponse {
    constructor(services) {
        this.Services = services;
    }
}
