namespace cs_server.Utils;

/// <summary>
/// Class-based source of truth for policy string variables
/// </summary>
public static class Policy
{
    public const string AllowSpecificOrigin = "AllowSpecificOrigin";
}

/// <summary>
/// Class-based source of truth for host address
/// </summary>
public static class HOST
{
    public const string Get = "http://localhost:5000";
}

public static class Header
{
    public const string ApplicationJson = "application/json";
}

public static class Route
{
    public const string DatabaseItem = "database/item";
}