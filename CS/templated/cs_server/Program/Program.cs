using CommandLine;
using cs_server.Program.Models.Server;

var (localHost, remoteHost) = ("localhost", "localhost");
var (localPort, remotePort) = (8000, 8000);

Parser.Default.ParseArguments<cs_server.Program.Utils.Options>(args)
    .WithParsed<cs_server.Program.Utils.Options>(opts =>
    {
        if (string.IsNullOrEmpty(opts.Local) && string.IsNullOrEmpty(opts.Remote))
        {
            Console.WriteLine("No argument provided, using defualts");
            return;
        }

        // Parse local opts
        string[] localComponents = opts.Local.Split(":");
        int localSeparator = opts.Local.IndexOf(':');
        if (localComponents.Length != 0)
        {
            switch (localSeparator)
            {
                case -1:
                    {
                        localHost = string.IsNullOrEmpty(localComponents[0]) ?
                            localHost : localComponents[0];
                        localPort = int.TryParse(localComponents[0], out int port) ?
                            port : localPort;
                        break;
                    }
                case 0:
                    {
                        localPort = int.TryParse(localComponents[1], out int port) ? port
                            : localPort;
                        break;

                    }
                default:
                    {
                        localHost = localComponents[0];
                        break;
                    }
            }
        }

        // Parse remote opts
        string[] remoteComponents = opts.Remote.Split(":");
        int remoteSeparator = opts.Remote.IndexOf(':');
        if (remoteComponents.Length != 0)
        {

            switch (remoteSeparator)
            {
                case -1:
                    {
                        remoteHost = string.IsNullOrEmpty(remoteComponents[0]) ?
                            remoteHost : remoteComponents[0];
                        remotePort = int.TryParse(remoteComponents[0], out int port) ? port
                            : remotePort;
                        break;
                    }
                case 0:
                    {
                        remotePort = int.TryParse(remoteComponents[1], out int port) ? port
                            : remotePort;
                        break;

                    }
                default:
                    {
                        remoteHost = remoteComponents[0];
                        break;
                    }
            }
        }
    })
    .WithNotParsed<cs_server.Program.Utils.Options>(errs =>
    {
        Console.WriteLine(errs);
    });

Node node = new(new ServerInfo(localHost, localPort), new ServerInfo(remoteHost, remotePort));
node.Init();
