# GoBot

A personal Discord bot project built with [Go](https://go.dev/), utilizing the [DiscordGo](https://github.com/bwmarrin/discordgo) library for Discord API integration.

<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about">About</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#command-registration">Command Registration</a></li>
  </ol>
</details>

## About

Gobot is a side project of mine which was solely created to understand the "Go" language better.

### Built With

- [Go](https://go.dev/) - The programming language used
- [DiscordGo](https://github.com/bwmarrin/discordgo) - Go package for Discord API integration
- [GoDotEnv](https://github.com/joho/godotenv) - Environment variable management
- [cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- [Docker](https://www.docker.com/) - Containerization

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (1.24 or higher)
- [Docker](https://www.docker.com/get-started) (optional)
- Discord Bot Token (from [Discord Developer Portal](https://discord.com/developers/applications))

### Installation

1. Clone the repository

```bash
git clone https://github.com/bluejutzu/GoBot.git
```

2. Install dependencies

```bash
go mod download
```

3. Create a `.env` file in the root directory with the following:

```env
BOT_TOKEN="your_bot_token_here"
```

or rename the [`.env.example`](/.env.example) to `.env` and change the `BOT_TOKEN` value.

## Usage

### Running Locally

```bash
go run main.go
```

### Using Docker

```bash
# Build the container
docker build -t gobot .

# Run the container
docker run -d --name gobot --env-file .env gobot
```

### Updating Docker Container

When you make changes to your `.env` file or any other files, you'll need to rebuild and restart the container:

```bash
# Using Docker Compose (recommended):
docker compose down
docker compose up -d --build

# Or using Docker directly:
docker stop gobot
docker rm gobot
docker build -t gobot .
docker run -d --name gobot --env-file .env gobot
```

## Features

- `/ping` - Check if the bot is online
- `/what-is-my-id` - Get your Discord user ID
- _Working on more features_

## Command Registration

Commands in GoBot are registered through a structured system:

1. Command Definition

   - Commands are defined as ApplicationCommand structs in the `commands` slice
   - Each command specifies its `name`, `description`, and any `options`

2. Command Handler Mapping

   - Each command has a corresponding handler function in the `commandHandlers` map
   - Handlers process the command when it's triggered by a user

3. Registration Process
   - Commands are automatically registered with Discord's API during bot startup
   - The bot creates `application_commands` for each defined command

__**Example of how commands are structured:**__

```go
package misc

// Command definition
var ExampleCommand = &discordgo.ApplicationCommand{
    Name: "example",
    Description: "An example command",
}

// Hanlder function
func HandleExampleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
  // ...add code here
}

```

__**And in the [bot/bot.go](/bot/bot.go#L30) file:**__
```go
// command mapping
commands = []*discordgo.ApplicationCommand {
	misc.ExampleCommand,
}
// Handler mapping
commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	"example": HandleExampleCommand,
}

```

## Command Line Interface (CLI)

Gobot has it's own CLI, named `gobot`, that has tools and utilities for this certain project.
> [!WARNING]
> The CLI is in a very experimental stage in which it might not behave as expected.

# [License (MIT)](https://github.com/Bluejutzu/GoBot?tab=MIT-1-ov-file)
