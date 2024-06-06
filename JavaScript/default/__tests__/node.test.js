import request from "supertest";
import { Node } from "../Server/node";
import { Mode } from "../common/utils";

describe("Node", () => {
    let node;

    beforeEach(() => {
        node = new Node("localhost", "8000", "localhost", "8000", true);
    });

    it("should create a Node object with correct properties", async () => {
        expect(node.Info).toHaveProperty("Uuid");
        expect(node.Info.Host).toBe("localhost");
        expect(node.Info.Port).toBe(8000);
        expect(node.Leader.Host).toBe("localhost");
        expect(node.Leader.Port).toBe(8000);
        expect(node.Type).toBe(Mode.LeaderNode);
        expect(node.InputBuffer).toEqual([]);
        expect(node.OutputBuffer).toEqual([]);
        expect(node.NodeList).toEqual([]);
        expect(node.Services).toEqual([]);

        await node.Shutdown();
    });

    it("Start Leader and connect follower, verify follower exists in NodeList, then disconnect the follower and verify follower removed", async () => {
        node.Start();
        const followerNode = new Node("localhost", "8001", "localhost", "8000", false);
        followerNode.Start();

        // Wait 3 seconds to give the nodes time to communicate, unsure of how to do this better
        await new Promise(resolve => setTimeout(resolve, 2000));

        // Verify the node exists in the Leader's NodeList
        expect(node.NodeList.length).toBe(1);
        // Verify the Leader's uuid has been sent back and added appropriately
        expect(followerNode.Leader.Uuid).toBe(node.Info.Uuid);

        followerNode.AddInput("exit");

        // Wait 3 seconds to give the nodes time to communicate, unsure of how to do this better
        await new Promise(resolve => setTimeout(resolve, 2000));
        expect(node.NodeList.length).toBe(0);

        await node.Shutdown();
    });

    it("should start the server and respond to requests", async () => {
        node.Start();
        const response = await request(node.Server).get("/node/registration");
        expect(response.status).toBe(404);

        await node.Shutdown();
    });
});