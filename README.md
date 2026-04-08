# GitPlus

GitPlus is a terminal application that simplifies GitHub project management.  
Connect your GitHub account, list repositories, create pull requests, and leverage AI to generate commit messages and PR descriptions—all from the command line.

## Features

- OAuth2 authentication with GitHub
- Browse your repositories in a sleek TUI
- Create pull requests with a few keystrokes
- AI‑assisted commit message and PR description generation (OpenAI)
- Lightweight and cross‑platform

## Installation

```bash
git clone https://github.com/iciwhite/gitplus.git
cd gitplus
make build
```

The binary will be placed in dist/gitplus.

Configuration

Copy .env.example to .env and fill in your credentials:

```
GITHUB_CLIENT_ID=your_github_oauth_app_client_id
GITHUB_CLIENT_SECRET=your_github_oauth_app_client_secret
OPENAI_API_KEY=your_openai_api_key
GITPLUS_PORT=8080
```

Usage

```bash
./dist/gitplus
```

On first run, you will be prompted to authorize the app with GitHub.
After authentication, you can:

· Type list and press Enter to refresh the repository list.
· Type pr to start a pull request workflow.
· Type ai followed by a diff or description to get AI suggestions.
