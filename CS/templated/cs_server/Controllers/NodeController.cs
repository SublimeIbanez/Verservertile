using System.Net.Mime;
using Microsoft.AspNetCore.Mvc;
using cs_server.Models.Server;

namespace cs_server.Controllers;

[Route(Utils.Route.NodeRegistration)]
[ApiController]
public class NodeController() : ControllerBase
{
    [HttpPost]
    [Produces(MediaTypeNames.Application.Json)]
    [Consumes(MediaTypeNames.Application.Json)]
    public IActionResult Post([FromBody] INodeRegistrationRequest request)
    {
        return Ok();
    }

    [HttpDelete("{id}")]
    [Produces(MediaTypeNames.Application.Json)]
    public IActionResult Delete()
    {
        return Ok();
    }
}
