param speechServiceName string
param keyVaultName string
param location string

resource account 'Microsoft.CognitiveServices/accounts@2023-05-01' = {
  name: speechServiceName
  location: location
  kind: 'SpeechServices'
  sku: {
    name: 'S0'
  }
  properties: {}
}

resource keyVault 'Microsoft.KeyVault/vaults@2022-07-01' existing = {
  name: keyVaultName
}

resource speechKey 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'speechKey'
  properties: {
    value: account.listKeys().key1
  }
}

resource speechRegion 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'speechRegion'
  properties: {
    value: account.location
  }
}

resource speechEndpoint 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: keyVault
  name: 'speechEndpoint'
  properties: {
    value: account.properties.endpoint
  }
}

output secretUris object = {
  speechKey: speechKey.properties.secretUri
  speechRegion: speechRegion.properties.secretUri
  speechEndpoint: speechEndpoint.properties.secretUri
}
