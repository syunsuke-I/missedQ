## What is missedQ

* This application detects conversations in Slack channels that have no replies.
* I made this app to make sure we don't **miss any unanswered** questions in my programming school's question channel.
* This application is intended to be regularly executed by some means, such as a cron job

## Setting Up the Environment

1. Create an `env` directory:

```zsh
mkdir env
```
2. Navigate to the env directory:
```zsh
cd env
```

3. Create a .env file:
```zsh
touch .env
```
4. Inside the .env file, add your Slack API token:
```.env
TOKEN=your_slack_token
```

##  Checking Your Configuration

* Please check the following settings in your configuration(setting.json):

  * monitoredChannel: Specify the Slack channel you want to monitor.
  * sendTo: Specify the Slack channel or user where messages will be sent.
 
```json
{
  "slackApiURL": "https://slack.com/api/conversations.history",
  "monitoredChannel": "Specify the Slack channel you want to monitor here",
  "sendTo": "Specify the Slack channel or user where messages will be sent here",
  "postMessage" : "https://slack.com/api/chat.postMessage"
}

```

## How the Notification System Works

```mermaid
sequenceDiagram
    participant Main as Main Application
    participant EnvLoad as Env_load()
    participant LoadConfig as LoadConfig()
    participant GetMessages as getMessagesFromSlack()
    participant FilterMessages as filterMessages()
    participant PostMessage as postMessageToSlack()

    Main ->> EnvLoad: Call Env_load()
    EnvLoad -->> Main: Token Loaded

    Main ->> GetMessages: Call getMessagesFromSlack()
    GetMessages ->> LoadConfig: Call LoadConfig()
    LoadConfig -->> GetMessages: Config Loaded
    GetMessages ->> HTTPRequest: Send HTTP Request
    HTTPRequest -->> GetMessages: Response Received
    GetMessages ->> JSONDecode: Decode JSON Response
    JSONDecode -->> GetMessages: Messages Parsed
    GetMessages -->> Main: Messages Retrieved

    Main ->> FilterMessages: Call filterMessages()
    FilterMessages -->> Main: Messages Filtered

    Main ->> PostMessage: Call postMessageToSlack()
    PostMessage ->> LoadConfig: Call LoadConfig()
    LoadConfig -->> PostMessage: Config Loaded
    PostMessage ->> ComposeMessage: Compose Message
    ComposeMessage -->> PostMessage: Message Composed
    PostMessage ->> HTTPRequest: Send HTTP Request
    HTTPRequest -->> PostMessage: Response Received
    PostMessage -->> Main: Message Sent
