import { argv } from "node:process";
import { Mode } from "./common/utils.js";
import { HandleServer } from "./Server/node.js";

const Arguments = Object.freeze({
    ModeArg: 0,
    LocalArg: 1,
    RemoteArg: 2,
    None: 3,
});

let remote = "";
let local = "localhost:8080";
let mode = Mode.None;
let currentMode = Arguments.None;

const MatchArgs = (arg) => {
    switch (arg.toLowerCase()) {
        case "-m":
        case "--mode":
            return Arguments.ModeArg;
        case "-l":
        case "--local":
            return Arguments.LocalArg;
        case "-r":
        case "--remote":
            return Arguments.RemoteArg;
        default:
            return Arguments.None;
    }
};

argv.slice(0).forEach(arg => {
    switch (currentMode) {
        case Arguments.ModeArg: {
            switch (arg.toLowerCase().trim()) {
                case "client":
                    mode = Mode.Client;
                    break;
                case "node":
                    mode = Mode.ServerNode;
                    break;
            }
            break;
        }

        case Arguments.LocalArg: {
            // Parse and error check the argument
            let pos = arg.indexOf(":");
            let values = arg.split(":");
            if (pos >= 0 && values.length > 1) {
                try {
                    parseInt(values[1]);
                } catch (error) {
                    console.log(`Could not parse port - defaulting to 8080: ${error}`);
                    values[1] = "8080";
                }
            } else {
                values[1] = "8080";
            }
            local = `${values[0]}:${values[1]}`
            break;
        }

        case Arguments.RemoteArg: {
            let pos = arg.indexOf(":");
            if (pos === -1) {
                console.log("Must pass correct remote\n  Expected: hostname:port");
                process.exit(1);
            }
            let values = arg.split(":");
            if (pos > 0 && values.length > 1) {
                try {
                    parseInt(values[1]);
                } catch (error) {
                    console.log(`Could not parse port - defaulting to 8080: ${error}`);
                    values[1] = "8080";
                }
            } else {
                values[1] = "8080";
            }
            remote = `${values[0]}:${values[1]}`
            break;
        }
    }
    currentMode = MatchArgs(arg)
});

switch (mode) {
    case Mode.Client: {

        break;
    }

    case Mode.ServerNode: {
        HandleServer(remote, local);
        break;
    }

    default:
        console.log("Must pass a mode using --mode or -m. Expected: client | node");
        process.exit(1);
}
