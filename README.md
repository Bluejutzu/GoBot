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
  </ol>
</details>

## About

Gobot is a side project of mine which was solely created to understand the "Go" language better.

### Built With

* [Go](https://go.dev/) - The programming language used
* [DiscordGo](https://github.com/bwmarrin/discordgo) - Go package for Discord API integration
* [GoDotEnv](https://github.com/joho/godotenv) - Environment variable management
* [Docker](https://www.docker.com/) - Containerization

## Getting Started

### Prerequisites

* [Go](https://go.dev/dl/) (1.24 or higher)
* [Docker](https://www.docker.com/get-started) (optional)
* Discord Bot Token (from [Discord Developer Portal](https://discord.com/developers/applications))

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
- *Working on more features*

### [License](https://github.com/Bluejutzu/GoBot?tab=MIT-1-ov-file)
