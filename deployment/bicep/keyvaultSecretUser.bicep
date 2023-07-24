// This template has been split from the main API template to allow use of principalID in the assignment name. 
// This avoids having to delete the previous assignment from the vault if API web app or function has been recreated
// Also supports key vaults in separate resource groups, though use there always creates a semi-invisible dependency

param keyVaultName string
param identityPrincipalId string

resource keyVault 'Microsoft.KeyVault/vaults@2022-07-01' existing = {
  name: keyVaultName
}

resource apiSecretsUser 'Microsoft.Authorization/roleAssignments@2020-08-01-preview' = {
  name: guid('${keyVault.id}-${identityPrincipalId}-KeyVaultSecretsUser')
  scope: keyVault
  properties: {
    // Key Vault Secrets User
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', '4633458b-17de-408a-b874-0445c86b69e6')
    principalId: identityPrincipalId
  }
}
