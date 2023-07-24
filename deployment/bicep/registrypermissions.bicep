param registryName string
param identityPrincipalId string

resource registry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' existing = {
  name: registryName
}
//AcrPull
resource roleAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  name: guid(identityPrincipalId, registryName, 'AcrPull')
  scope: registry
  properties: {
    principalId: identityPrincipalId
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '7f951dda-4ed3-4680-a7ca-43fe172d538d')
  }
}
