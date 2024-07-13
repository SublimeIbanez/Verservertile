using System.Numerics;
using CommandLine;
using cs_server.Utils;

string localHost = "localhost";
int localPort = 8000;
string remoteHost = "localhost";
int remotePort = 8000;

Parser.Default.ParseArguments<Options>(args)
    .WithParsed<Options>(opts =>
    {
        if (opts.Local != String.Empty)
        {
            string[] localComponents = opts.Local.Split(":");

            if (localComponents.Length != 2)
            {
                int separator = opts.Local.IndexOf(':');
                if (separator != -1) {
                    localPort = int.TryParse(localComponents[0], out int port) ? port : localPort;
                }
            }
        }
    })
    .WithNotParsed<Options>(errs =>
    {

    });





WebApplicationBuilder builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
    options.AddPolicy(Policy.AllowSpecificOrigin,
        builder => builder.WithOrigins(cs_server.Utils.Host.Get)
            .AllowAnyMethod()
            .AllowAnyHeader());
});

WebApplication webApp = builder.Build();

// Configure the HTTP request pipeline.
if (webApp.Environment.IsDevelopment())
{
    webApp.UseSwagger();
}

webApp.UseHttpsRedirection();
webApp.UseStaticFiles();

webApp.UseRouting();

webApp.Run();
