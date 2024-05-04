FROM mcr.microsoft.com/dotnet/sdk:8.0
WORKDIR /CS_Server
COPY ["CS_Server.csproj", "/CS_Server"]

RUN dotnet restore "CS_Server.csproj"

COPY . .

CMD ["sh"]