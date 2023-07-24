param userAssignedIdentityName string
param location string
param registryResourceId string
param keyVaultName string

resource userAssignedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: userAssignedIdentityName
  location: location
}

var registryName = split(registryResourceId, '/')[8]
var registryRg = split(registryResourceId, '/')[4]

module registryPermissions 'registrypermissions.bicep' = {
  name: 'registryPermissions'
  params: {
    registryName: registryName
    identityPrincipalId: userAssignedIdentity.properties.principalId
  }
  scope: resourceGroup(registryRg)
}

module keyVaultPermissions 'keyvaultSecretUser.bicep' = {
  name: 'keyVaultPermissions'
  params: {
    keyVaultName: keyVaultName
    identityPrincipalId: userAssignedIdentity.properties.principalId
  }
}

output userAssignedIdentityResourceId string = userAssignedIdentity.id
output userAssignedIdentityObjectId string = userAssignedIdentity.properties.principalId
output userAssignedIdentityClientId string = userAssignedIdentity.properties.clientId
