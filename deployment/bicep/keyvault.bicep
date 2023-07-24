param location string
param keyVaultName string

@secure()
@description('Telegram bot token')
param botToken string

@description('The objectId of the group to give full permissions to this keyvault to')
param developerGroupObjectId string
param logAnalyticsResourceId string

resource keyVault 'Microsoft.KeyVault/vaults@2022-07-01' = {
  name: keyVaultName
  location: location
  properties: {
    enabledForDeployment: false
    enabledForTemplateDeployment: true
    enabledForDiskEncryption: false
    enableSoftDelete: true
    enablePurgeProtection: true
    tenantId: subscription().tenantId
    sku: {
      name: 'standard'
      family: 'A'
    }
    enableRbacAuthorization: true
  }
}

resource developerGroupFullControl 'Microsoft.Authorization/roleAssignments@2020-08-01-preview' = {
  name: guid('${keyVault.id}-${developerGroupObjectId}-KeyVaultAdministrator')
  scope: keyVault
  properties: {
    // Key Vault Administrator
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '00482a5a-887f-4fb3-b363-3b7fe8e74483')
    principalId: developerGroupObjectId
  }
}

resource diagnosticsSettings 'Microsoft.Insights/diagnosticSettings@2021-05-01-preview' = {
  name: 'AuditLogsToLogAnalytics'
  scope: keyVault
  properties: {
    workspaceId: logAnalyticsResourceId
    logs: [
      {
        category: 'AuditEvent'
        enabled: true
      }
    ]
  }
}

resource botTokenSecret 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'botToken'
  properties: {
    value: botToken
  }
}

output keyVaultName string = keyVault.name
output keyVaultId string = keyVault.id

output secretUris object = {
  botToken: botTokenSecret.properties.secretUri
}
