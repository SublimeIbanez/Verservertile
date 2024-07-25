namespace cs_server.Models.Server;

public sealed record ServerInfo(string Host, int Port, Guid Uuid = new())
{
    public bool IsEqual(ServerInfo other)
    {
        if (Host == other.Host && Port == other.Port && Uuid == other.Uuid)
        {
            return true;
        }
        return false;
    }
}