/**
 * Entity Mode Types
 * @type {number}
 */
export const Mode = Object.freeze({
    None: 0,
    Client: 1,
    ServerNode: 2,
    LeaderNode: 3,
});

/**
 * @param {boolean} value 
 * @param {string} message 
 */
export const Assert = (value, message) => {
    if (!value) {
        console.error(message);
    }
}

export const Type = Object.freeze({
    String: "string",
    Number: "number",
    Boolean: "boolean",
});

/**
 * REST API paths
 * @type {Object<string, string>}
 */
export const Path = Object.freeze({
    // ################## Node Communication ##################
    /** Node registration:POST | deregistration:DELETE */
    NodeRegistration: "/node/registration",

    // ################## Authentication ##################
    /** User log-  in:PUT | out:POST */
    AuthAccess: "/auth/access",
    /** User registration:PUT */
    AuthRegister: "/auth/register",

    // ################## DB Management ##################
    /** Database add:PUT | remove:DELETE | edit:POST | query:GET category */
    DatabaseCategory: "/database/category",
    /** Database add:PUT | remove:DELETE | edit:POST | query:GET item */
    DatabaseItem: "/database/item",
});
