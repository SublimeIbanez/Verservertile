namespace cs_server.Utils;
using CommandLine;

public class Options
{
    [Option('l', "local", Required = false, HelpText = "Host and port information for the local machine")]
    public required string Local { get; set; } = String.Empty;


    [Option('r', "remote", Required = false, HelpText = "Host and port information for the remote machine")]
    public required string Remote { get; set; } = String.Empty;
}