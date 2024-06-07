import { Node } from "../Server/node";
import { StatusCode } from "../common/protocol";
import { Path } from "../common/utils";

describe("Protocol", () => {
    let node;

    beforeEach(() => {
        node = new Node("localhost", "8000", "localhost", "8000", true);
    });

    // // ################## Node Communication ##################
    // /** Node registration:POST | deregistration:DELETE */
    // NodeRegistration: "/node/registration",
    //------------
    // Delete: "DELETE",
    // Get: "GET",
    // Head: "HEAD",
    // Post: "POST",
    // Put: "PUT",
    // Trace: "TRACE",
    it("Should provide correct responses to node communication", async () => {
        node.Start();

        // add a new "node" to the list of nodes
        const putResponse = await request(node.Server).put(Path.NodeRegistration);
        expect(putResponse.status).toBe(404);

        {
            const deleteResponse = await request(node.Server).delete(Path.NodeRegistration);
            expect(deleteResponse.status).toBe(StatusCode.NotFound);
            const getResponse = await request(node.Server).get(Path.NodeRegistration);
            expect(getResponse.status).toBe(StatusCode.NotFound);
            const postResponse = await request(node.Server).post(Path.NodeRegistration);
            expect(postResponse.status).toBe(StatusCode.NotFound);
            const traceResponse = await request(node.Server).trace(Path.NodeRegistration);
            expect(traceResponse.status).toBe(StatusCode.NotFound);
        }

        await node.Shutdown();
    });

    // // ################## Authentication ##################
    // /** User log-  in:Get | out:Trace | registration:PUT | edit:POST */
    // AuthAccess: "/auth/access",
    //------------
    // Delete: "DELETE",
    // Get: "GET",
    // Head: "HEAD",
    // Post: "POST",
    // Put: "PUT",
    // Trace: "TRACE",
    it("Should provide correct responses to authentication", async () => {
        node.Start();

        // add a new "node" to the list of nodes
        const putResponse = await request(node.Server).put(Path.NodeRegistration);
        expect(putResponse.status).toBe(404);

        const deleteResponse = await request(node.Server).delete(Path.AuthAccess);
        expect(deleteResponse.status).toBe(StatusCode.NotFound);
        const getResponse = await request(node.Server).get(Path.AuthAccess);
        expect(getResponse.status).toBe(StatusCode.NotFound);
        const postResponse = await request(node.Server).post(Path.AuthAccess);
        expect(postResponse.status).toBe(StatusCode.NotFound);
        const traceResponse = await request(node.Server).trace(Path.AuthAccess);
        expect(traceResponse.status).toBe(StatusCode.NotFound);

        await node.Shutdown();
    });

    // // ################## DB Management ##################
    // /** Database add:PUT | remove:DELETE | edit:POST | query:GET category */
    // DatabaseCategory: "/database/category",
    // /** Database add:PUT | remove:DELETE | edit:POST | query:GET item */
    // DatabaseItem: "/database/item",
    //------------
    // Delete: "DELETE",
    // Get: "GET",
    // Head: "HEAD",
    // Post: "POST",
    // Put: "PUT",
    // Trace: "TRACE",
    it("Should provide correct responses to authentication", async () => {
        node.Start();

        // Add a new category
        const putCatResponse = await request(node.Server).put(Path.DatabaseCategory);
        expect(putCatResponse.status).toBe(404);

        const deleteCatResponse = await request(node.Server).delete(Path.DatabaseCategory);
        expect(deleteCatResponse.status).toBe(StatusCode.NotFound);
        const getCatResponse = await request(node.Server).get(Path.DatabaseCategory);
        expect(getCatResponse.status).toBe(StatusCode.NotFound);
        const postCatResponse = await request(node.Server).post(Path.DatabaseCategory);
        expect(postCatResponse.status).toBe(StatusCode.NotFound);
        const traceCatResponse = await request(node.Server).trace(Path.DatabaseCategory);
        expect(traceCatResponse.status).toBe(StatusCode.NotFound);

        // Put a new item
        const putItemResponse = await request(node.Server).put(Path.DatabaseCategory);
        expect(putItemResponse.status).toBe(404);

        const deleteItemResponse = await request(node.Server).delete(Path.DatabaseCategory);
        expect(deleteItemResponse.status).toBe(StatusCode.NotFound);
        const getItemResponse = await request(node.Server).get(Path.DatabaseCategory);
        expect(getItemResponse.status).toBe(StatusCode.NotFound);
        const postItemResponse = await request(node.Server).post(Path.DatabaseCategory);
        expect(postItemResponse.status).toBe(StatusCode.NotFound);
        const traceItemResponse = await request(node.Server).trace(Path.DatabaseCategory);
        expect(traceItemResponse.status).toBe(StatusCode.NotFound);
        await node.Shutdown();
    });
});