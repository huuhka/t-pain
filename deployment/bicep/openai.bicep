param location string
param openAiServiceName string
param keyVaultName string

param logAnalyticsResourceId string

resource openAiService 'Microsoft.CognitiveServices/accounts@2023-05-01' = {
  name: openAiServiceName
  location: location
  sku: {
    name: 'S0'
  }
  kind: 'OpenAI'
  properties: {
    publicNetworkAccess: 'Enabled'
    customSubDomainName: openAiServiceName
    disableLocalAuth: true
  }
}

resource openAi_diagnosticsSettings 'Microsoft.Insights/diagnosticSettings@2021-05-01-preview' = {
  name: 'appServiceAuditToLogAnalytics'
  scope: openAiService
  properties: {
    workspaceId: logAnalyticsResourceId
    logs: [
      {
        category: 'Audit'
        enabled: true
      }
    ]
  }
}

resource openAI_deployment 'Microsoft.CognitiveServices/accounts/deployments@2023-05-01' = {
  parent: openAiService
  name: 'gpt35'
  sku: {
    name: 'S0'
  }
  properties: {
    model: {
      name: 'gpt-35-turbo'
      version: '0613'
      format: 'OpenAi'
    }
  }
}

resource keyVault 'Microsoft.KeyVault/vaults@2023-02-01' existing = {
  name: keyVaultName
}

resource openAiEndpoint 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'openAiEndpoint'
  properties: {
    value: openAiService.properties.endpoint
  }
}

resource openAiKey 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'openAiKey'
  properties: {
    value: openAiService.listKeys().key1
  }
}

resource openAiDeployment 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'openAiDeployment'
  properties: {
    value: openAiService.properties.endpoint
  }
}

output secretUris object = {
  openAiEndpoint: openAiEndpoint.properties.secretUri
  openAiKey: openAiKey.properties.secretUri
  openAiDeployment: openAiDeployment.properties.secretUri
}
