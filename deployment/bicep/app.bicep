param containerAppName string
param containerAppEnvName string

@description('name of the container image to use. Should not include tag, registry name or :')
param containerImage string = 't-painbot'
param containerTag string
param userAssignedIdentityResourceId string
param userAssignedIdentityClientId string

param location string

param speechSecretUris object
param openAiSecretUris object
param dataCollectionSecretUris object
param telegramSecretUris object

param logAnalyticsName string

param registryResourceId string

var registryName = split(registryResourceId, '/')[8]
var registryRg = split(registryResourceId, '/')[4]

resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2022-10-01' existing = {
  name: logAnalyticsName
}

resource environment 'Microsoft.App/managedEnvironments@2023-04-01-preview' = {
  name: containerAppEnvName
  location: location
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: logAnalytics.properties.customerId
        sharedKey: logAnalytics.listKeys().primarySharedKey
      }
    }
  }
}

resource registry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' existing = {
  name: registryName
  scope: resourceGroup(registryRg)
}

resource app 'Microsoft.App/containerApps@2023-04-01-preview' = {
  name: containerAppName
  location: location
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${userAssignedIdentityResourceId}': {}
    }
  }
  properties: {
    environmentId: environment.id
    configuration: {
      activeRevisionsMode: 'Single'
      secrets: [
        {
          name: 'bot-token'
          keyVaultUrl: telegramSecretUris.botToken
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'speech-region'
          keyVaultUrl: speechSecretUris.speechRegion
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'speech-key'
          keyVaultUrl: speechSecretUris.speechKey
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'data-collection-endpoint'
          keyVaultUrl: dataCollectionSecretUris.dataCollectionEndpoint
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'data-collection-rule-id'
          keyVaultUrl: dataCollectionSecretUris.dataCollectionRuleId
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'data-collection-stream-name'
          keyVaultUrl: dataCollectionSecretUris.dataCollectionStreamName
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'open-ai-key'
          keyVaultUrl: openAiSecretUris.openAiKey
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'open-ai-deployment'
          keyVaultUrl: openAiSecretUris.openAiDeployment
          identity: userAssignedIdentityResourceId
        }
        {
          name: 'open-ai-endpoint'
          keyVaultUrl: openAiSecretUris.openAiEndpoint
          identity: userAssignedIdentityResourceId
        }
      ]
      registries: [
        {
          server: registry.properties.loginServer
          identity: userAssignedIdentityResourceId
        }
      ]
    }
    template: {
      scale: {
        maxReplicas: 1
        minReplicas: 1
      }
      containers: [
        {
          name: containerAppName
          image: '${registry.properties.loginServer}/${containerImage}:${containerTag}'
          env: [
            {
              name: 'BOT_TOKEN'
              secretRef: 'bot-token'
            }
            {
              name: 'DATA_COLLECTION_ENDPOINT'
              secretRef: 'data-collection-endpoint'
            }
            {
              name: 'DATA_COLLECTION_RULE_ID'
              secretRef: 'data-collection-rule-id'
            }
            {
              name: 'DATA_COLLECTION_STREAM_NAME'
              secretRef: 'data-collection-stream-name'
            }
            {
              name: 'OPENAI_DEPLOYMENT'
              secretRef: 'open-ai-deployment'
            }
            {
              name: 'OPENAI_ENDPOINT'
              secretRef: 'open-ai-endpoint'
            }
            {
              name: 'OPENAI_KEY'
              secretRef: 'open-ai-key'
            }
            {
              name: 'SPEECH_KEY'
              secretRef: 'speech-key'
            }
            {
              name: 'SPEECH_REGION'
              secretRef: 'speech-region'
            }
            {
              name: 'AZURE_CLIENT_ID'
              value: userAssignedIdentityClientId
            }
          ]
          resources: {
            cpu: 1
            memory: '2Gi'
          }
        }
      ]
    }
  }
}
