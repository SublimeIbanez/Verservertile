namespace cs_server.Models.Server;

public class Node(ServerInfo nodeInfo, ServerInfo leaderInfo)
{
    public ServerInfo NodeInfo { get; private set; } = nodeInfo;
    public ServerInfo LeaderInfo { get; private set; } = leaderInfo;
    public bool Leader { get; private set; } = nodeInfo.IsEqual(leaderInfo);

    public void Init()
    {
        Console.WriteLine($"Initializing {(Leader ? "Leader" : "Follower")} node");
        WebApplicationBuilder builder = WebApplication.CreateBuilder();
        builder.Services.AddCors(options =>
        {
            options.AddPolicy(Utils.Policy.AllowSpecificOrigin,
                builder => builder.WithOrigins(Utils.Host.Get)
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
    }
}