using System.Net;
using System.Net.Sockets;

namespace NodeSocket
{
    class NodeSocket
    {
        ProtocolType Protocol { get; }
        string Address { get; }
        UInt16 Port { get; }
        Socket Socket { get;}


        public NodeSocket(ProtocolType protocol, string address, UInt16 port) {
            if (protocol != ProtocolType.Tcp || protocol != ProtocolType.Udp) {
                throw new ArgumentException("Only TCP or UDP protocols can be used");
            }
            Protocol = protocol;
            Address = address;
            Port = port;
            Socket = new Socket(protocol == ProtocolType.Tcp ? SocketType.Stream : SocketType.Dgram, protocol);
        }
    }
}