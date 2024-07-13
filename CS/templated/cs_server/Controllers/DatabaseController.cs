using System.Net.Mime;
using Microsoft.AspNetCore.Mvc;

namespace cs_server.Controllers;

[Route(Utils.Route.DatabaseItem)]
[ApiController]
public class DatabaseItemController() : ControllerBase
{
    [HttpGet("{item_id}")]
    [Produces(MediaTypeNames.Application.Json)]
    public IEnumerable<string> Get(string item_id)
    {
        return [""];
    }

    [HttpPost]
    [Produces(MediaTypeNames.Application.Json)]
    public string Post()
    {
        return "";
    }
}