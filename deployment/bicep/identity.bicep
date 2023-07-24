param userAssignedIdentityName string
param location string
param registryResourceId string

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
    principalId: userAssignedIdentity.properties.principalId
  }
  scope: resourceGroup(registryRg)
}

output userAssignedIdentityResourceId string = userAssignedIdentity.id
output userAssignedIdentityObjectId string = userAssignedIdentity.properties.principalId
output userAssignedIdentityClientId string = userAssignedIdentity.properties.clientId
