# Verservertile

The entire goal behind this project is to provide a standardized testing frontend and backend development center.

## Start
- To run the program, provide the relevant compile/run command along with the requisite arguments:
  - `-l <host:port>` or `--local <host:port>` (e.g. `-l localhost:8000`)
  - `-r <host:port>` or `--remote <host:port>` (e.g. `-r localhost:8000`)
  - NOTE: Default will use `localhost:8000`

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
  - DB: categoryID - categoryName
- Items: "/database/item"
  - Database add:POST | remove:POST | edit:POST | query:GET item
  - DB: itemID - categoryName - itemName - quantity - price

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
  // Query category
  "hostname": "localhost",
  "port": 8000,
  "path": "/database/category/item_id",
  "method": "GET",
  "headers": {
    "Content-Type": "application/json"
  },
}

{
  // Add/Edit items
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
  // Add/Edit items
  "hostname": "localhost",
  "port": 8000,
  "path": "/database/category/item_id",
  "method": "GET",
  "headers": {
    "Content-Type": "application/json"
  },
}
```
