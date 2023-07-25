# T-Pain Bot

T-Pain Bot is a telegram bot that helps the user track their daily pain levels using either audio or text messages. It utilizes Azure Cognitive Services, OpenAI and Azure Log Analytics and at this point is designed for only a very limited set of users.

# Requirements

The following environment variables are expected to be populated before running the bot:

- **BOT_TOKEN**: The token of the telegram bot
- **DATA_COLLECTION_ENDPOINT**: Log ingestion address of the data collection endpoint
- **DATA_COLLECTION_RULE_ID**: the immutable rule id of the azure monitor data collection rule
- **DATA_COLLECTION_STREAM_NAME**: the name of the "Data Source" in Azure Monitor Data collection rule. Should be something like "Custom-YourDataTableName_CL"
- **OPENAI_DEPLOYMENT_NAME**: the name of the Azure OpenAI deployment
- **OPENAI_KEY**: the api key of the Azure OpenAI deployment
- **OPENAI_ENDPOINT**: the endpoint of the Azure OpenAI service
- **SPEECH_KEY**: the api key of the Azure Cognitive Services Speech service
- **SPEECH_REGION**: the region of the Azure Cognitive Services Speech service

You might also need to install the Speech Service SDK for Go. It's a bit of a mess:
https://learn.microsoft.com/en-us/azure/ai-services/speech-service/quickstarts/setup-platform?pivots=programming-language-go&tabs=windows,ubuntu,dotnetcli,dotnet,jre,maven,browser,mac,pypi#platform-requirements

Or just build the dockerfile, that should work with non-mac environments.

# Usage

The bot is currently only usable by the users listed in models/users.go. It's also only tested in a private chat.

The user message should contain description of their current pains: their location, levels from 0-10 and optionally
further description regarding radiation, numbness etc.

The bot will then generate an object based on the data given and log it into Azure Log Analytics.

The user has access to a Azure workbook that allows them to use premade charts of their data and create
their own queries based on Kusto Query Language.

## IaC Deployment

```powershell
Dev:
New-AzResourceGroupDeployment -ResourceGroupName "t-pain-dev" -TemplateFile ./deployment/bicep/main.bicep -TemplateParameterFile ./deployment/bicep/params-dev.json

Prod:
New-AzResourceGroupDeployment -ResourceGroupName "t-pain-prod" -TemplateFile ./deployment/bicep/main.bicep -TemplateParameterFile ./deployment/bicep/params-prod.json
```

# Notes

- The current implementation with Azure Log Analytics is not perfect. The obvious tradeoff is that the data cannot
  be edited or deleted by the user (or really, the admin either). However, as the use case is for
  such a limited userbase, we can probably live with that.
- In later versions, we can easily import the current data to another database, if needed.

# Still missing

- First implementation of the visualization on top of the data (e.g. Azure Workbooks)
- Better testing (or any, really)
- Health check support in the container
- /about or other commands support
- Maybe some kind of approval flow from the user to avoid saving incorrect data?
