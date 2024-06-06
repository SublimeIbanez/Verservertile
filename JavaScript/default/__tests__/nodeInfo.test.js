import { NodeInfo } from "../Server/node";

describe('NodeInfo', () => {
    it("Create a NodeInfo instance with requisite properties", () => {
        const uuid = 'test-uuid';
        const host = 'localhost';
        const port = 8000;
        const nodeInfo = new NodeInfo(uuid, host, port);

        expect(nodeInfo.Uuid).toBe(uuid);
        expect(nodeInfo.Host).toBe(host);
        expect(nodeInfo.Port).toBe(port);
    });

    it("port should be parsed as an integer", () => {
        const nodeInfo = new NodeInfo("test-uuid", "localhost", "8000");
        expect(nodeInfo.Port).toBe(8000);
    });
});