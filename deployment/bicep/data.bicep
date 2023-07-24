param location string
param logAnalyticsName string
param dataCollectionEndpointName string
param dataCollectionRuleName string
param keyVaultName string

param developerGroupObjectId string
param managedIdentityObjectId string

@description('Short name for the table, used for the stream name and table name. Should not contain the _CL ending. The template will handle that.')
param painTableShortName string = 'PainDescriptions'
var realTableName = '${painTableShortName}_CL'
var dataCollectionStreamName = 'Custom_${painTableShortName}_CL'

var tableSchema = [
  {
    name: 'timestamp'
    type: 'datetime'
  }
  {
    name: 'level'
    type: 'int'
  }
  {
    name: 'locationId'
    type: 'int'
  }
  {
    name: 'sideId'
    type: 'int'
  }
  {
    name: 'description'
    type: 'string'
  }
  {
    name: 'numbness'
    type: 'boolean'
  }
  {
    name: 'numbnessDescription'
    type: 'string'
  }
  {
    name: 'locationName'
    type: 'string'
  }
  {
    name: 'sideName'
    type: 'string'
  }
  {
    name: 'userName'
    type: 'string'
  }
]

resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2022-10-01' existing = {
  name: logAnalyticsName
}

resource customTable 'Microsoft.OperationalInsights/workspaces/tables@2022-10-01' = {
  name: realTableName
  parent: logAnalytics
  properties: {
    plan: 'Analytics'
    retentionInDays: 730
    totalRetentionInDays: 1825
    schema: {
      name: painTableShortName
      columns: tableSchema
    }
  }
}

resource dataCollectionEndpoint 'Microsoft.Insights/dataCollectionEndpoints@2022-06-01' = {
  name: dataCollectionEndpointName
  location: location
  properties: {
    networkAcls: {
      publicNetworkAccess: 'Enabled'
    }
  }
}

resource dataCollectionRule 'Microsoft.Insights/dataCollectionRules@2022-06-01' = {
  name: dataCollectionRuleName
  location: location
  properties: {
    destinations: {
      logAnalytics: [
        {
          workspaceResourceId: logAnalytics.id
          name: guid(logAnalytics.id)
        }
      ]
    }
    dataCollectionEndpointId: dataCollectionEndpoint.id
    dataFlows: [
      {
        streams: [
          dataCollectionStreamName
        ]
        destinations: [
          guid(logAnalytics.id)
        ]
        outputStream: dataCollectionStreamName
        transformKql: 'source | extend TimeGenerated = timestamp'
      }
    ]
    streamDeclarations: {
      '${dataCollectionStreamName}': {
        columns: tableSchema
      }
    }
  }
}

resource dataCollectionRulePublisherGroup 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  name: guid(developerGroupObjectId, dataCollectionEndpoint.id)
  scope: dataCollectionRule
  properties: {
    principalId: developerGroupObjectId
    // Monitoring Metrics Publisher
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '3913510d-42f4-4e42-8a64-420c390055eb')
  }
}

resource dataCollectionRulePublisherApi 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  name: guid(managedIdentityObjectId, dataCollectionEndpoint.id)
  scope: dataCollectionRule
  properties: {
    principalId: managedIdentityObjectId
    // Monitoring Metrics Publisher
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '3913510d-42f4-4e42-8a64-420c390055eb')
  }
}

resource keyVault 'Microsoft.KeyVault/vaults@2023-02-01' existing = {
  name: keyVaultName
}

resource dataCollectionEndpointSecret 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'dataCollectionEndpoint'
  properties: {
    value: dataCollectionEndpoint.properties.logsIngestion.endpoint
  }
}

resource dataCollectionRuleIdSecret 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'dataCollectionRuleId'
  properties: {
    value: dataCollectionRule.properties.immutableId
  }
}

resource dataCollectionStreamNameSecret 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'dataCollectionStreamName'
  properties: {
    value: dataCollectionStreamName
  }
}

output secretUris object = {
  dataCollectionEndpoint: dataCollectionEndpointSecret.properties.secretUri
  dataCollectionRuleId: dataCollectionRuleIdSecret.properties.secretUri
  dataCollectionStreamName: dataCollectionStreamNameSecret.properties.secretUri
}
