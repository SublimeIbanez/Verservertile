using System.Net.Mime;
using Microsoft.AspNetCore.Mvc;

namespace cs_server.Program.Controllers;

[Route(Utils.Route.NodeRegistration)]
[ApiController]
public class AuthController
{
    [HttpPost]
    [Produces(MediaTypeNames.Application.Json)]
    [Consumes(MediaTypeNames.Application.Json)]
    public string Post()
    {
        return "";
    }
}