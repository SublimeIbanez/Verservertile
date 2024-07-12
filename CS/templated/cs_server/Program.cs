using cs_server.Utils;

WebApplicationBuilder builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
    options.AddPolicy(Policy.AllowSpecificOrigin,
        builder => builder.WithOrigins(HOST.Get)
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
