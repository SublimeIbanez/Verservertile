namespace cs_server.Program.Models.Server;

public interface INodeRegistrationRequest
{
    public string Uuid { get; set; }
    public string Host { get; set; }
    public int Port { get; set; }
}