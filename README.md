# Verservertile

The entire goal behind this project is to provide a standardized testing frontend and backend development center.

## API Calls

- All will be with "application/json"

### Node Communication:

- "/node/registration"
  - Node registration:POST | deregistration:DELETE

#### Example:

```json
// Node Registration
{
    // Information for the leader
    "hostname": "localhost",
    "port": 8000,
    "path": "/node/registration",
    "method": "POST",
    "headers": {
        "Content-Type": "application/json",
    },
    "data": {
        "uuid": "globallyuniqueid", // Globally unique id
        "host": "localhost",
        "port": 8001,
    }
}

// Node Deregistration
{
    "hostname": "localhost",
    "port": 8000,
    "path": "/node/registration/globallyuniqueid",
    "method": "DELETE",
    "headers": {
        "Content-Type": "application/json",
    },
}
```

### Authentication:

- "/auth/access"
  - User sign -in:POST | -out:POST | register:POST

#### Example

```json
{
  // Information for the leader
  "hostname": "localhost",
  "port": 8000,
  "path": "/auth/access",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json"
  },
  "data": {
    "action": "sign_in", // "sign_in" | "sign_out" | "register"
    "username": "uniqueusername", // Globally unique id
    "password": "password"
  }
}
```

### DB Management -- UNDER CONSTRUCTION

- Catagories: "/database/category"
  - Database add:POST | remove:POST | edit:POST | query:GET category
- Items: "/database/item"
  - Database add:POST | remove:POST | edit:POST | query:GET item

#### Example

```json
{
  // Add/Edit categories
  "hostname": "localhost",
  "port": 8000,
  "path": "/database/category/session_id",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json"
  },
  "data": {
    "action": "add", // "add" | "edit" | "delete"
    // Fill in information here
  }
}
{
  // Add/Edit categories
  "hostname": "localhost",
  "port": 8000,
  "path": "/database/category/item_id",
  "method": "GET",
  "headers": {
    "Content-Type": "application/json"
  },
}

```
