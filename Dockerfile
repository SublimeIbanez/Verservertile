FROM mcr.microsoft.com/dotnet/sdk:8.0
WORKDIR /cs_server
COPY ["C#.csproj", "/cs_server"]

RUN dotnet restore "C#.csproj"

COPY . .

CMD ["sh"]