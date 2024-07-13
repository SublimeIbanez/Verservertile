namespace cs_server.Models.Server;
using cs_server.Utils;

public class Node(ServerId uuid, ServerInfo nodeInfo, ServerInfo leaderInfo)
{
    public ServerId Uuid { get; } = uuid;
    public ServerInfo NodeInfo { get; private set; } = nodeInfo;
    public ServerInfo LeaderInfo { get; private set; } = leaderInfo;
    public bool Leader { get; private set; }

    public void Init()
    {

    }
}