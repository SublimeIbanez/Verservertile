import request from "supertest";
import { Node } from "../Server/node";

describe("Node", () => {
  let node;

  beforeEach(() => {
    node = new Node("localhost", "8001", "localhost", "8000", true);
  });

  it("should create a Node object with correct properties", () => {
    expect(node.Info).toHaveProperty("Uuid");
    expect(node.Info.Host).toBe("localhost");
    expect(node.Info.Port).toBe(8001);
    expect(node.Leader.Host).toBe("localhost");
    expect(node.Leader.Port).toBe(8000);
    expect(node.Type).toBe("LeaderNode");
    expect(node.InputBuffer).toEqual([]);
    expect(node.OutputBuffer).toEqual([]);
    expect(node.NodeList).toEqual([]);
    expect(node.Services).toEqual([]);
  });

  it("should handle input correctly", () => {
    node.AddInput("test input");
    expect(node.InputBuffer).toContain("test input");
  });

  it("should handle output correctly", () => {
    node.AddOutput("test output");
    expect(node.OutputBuffer).toContain("test output");
  });

  it("should start the server and respond to requests", async () => {
    node.Start();
    const response = await request(node.Server).get("/node/registration");
    expect(response.status).toBe(404); // Adjust based on your expected behavior

    // Close the server after the test
    node.Server.close();
  });
});