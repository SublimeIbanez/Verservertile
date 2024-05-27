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
let local = "localhost:8000";
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
            let values = arg.split(":");
            // For now just use the passed argument
            local = arg;
            break;
        }

        case Arguments.RemoteArg: {
            // Parse and error check the argument
            let values = arg.split(":");
            // For now just use the passed argument
            remote = arg;
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
