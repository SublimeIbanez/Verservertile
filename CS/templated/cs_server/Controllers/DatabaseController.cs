using Microsoft.AspNetCore.Mvc;

namespace cs_server.Controllers;

[Route(Utils.Route.DatabaseItem)]
[ApiController]
public class DatabaseItemController() : ControllerBase
{
    [HttpGet("{item_id}")]
    [Produces(Utils.Header.ApplicationJson)]
    public IEnumerable<string> Get(string item_id)
    {
        return [""];
    }

    [HttpPost]
    [Produces(Utils.Header.ApplicationJson)]
    public string Post()
    {
        return "";
    }
}