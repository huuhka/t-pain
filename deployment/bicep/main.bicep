param appName string = 't-pain'
param environment string = 'dev'
param location string = 'westeurope'

@secure()
@description('Telegram Bot Token')
param botToken string

param registryResourceId string
param developerGroupObjectId string

param customTableName string = 'PainDesc'

param containerTag string

param buildId string = utcNow()

var namingPrefix = '${appName}-${environment}'
var naming = {
  logAnalytics: '${namingPrefix}-law'
  openAiService: '${namingPrefix}-oai'
  keyVault: '${namingPrefix}-kv'
  userAssignedIdentity: '${namingPrefix}-appid'
  dataCollectionRule: '${namingPrefix}-dcr'
  dataCollectionEndpoint: '${namingPrefix}-dce'
  customTable: customTableName
  speech: '${namingPrefix}-speech'
  containerAppEnv: '${namingPrefix}-appenv'
  containerApp: '${namingPrefix}-app'
}

module lawbase 'lawbase.bicep' = {
  name: 'lawbase-${buildId}'
  params: {
    location: location
    logAnalyticsName: naming.logAnalytics
  }
}

module identity 'identity.bicep' = {
  name: 'identity-${buildId}'
  params: {
    location: location
    userAssignedIdentityName: naming.userAssignedIdentity
    registryResourceId: registryResourceId
    keyVaultName: keyvault.outputs.keyVaultName
  }
}

module keyvault 'keyvault.bicep' = {
  name: 'keyvault-${buildId}'
  params: {
    location: location
    keyVaultName: naming.keyVault
    developerGroupObjectId: developerGroupObjectId
    logAnalyticsResourceId: lawbase.outputs.logAnalyticsId
    botToken: botToken
  }
}

module openAiService 'openai.bicep' = {
  name: 'openAi-${buildId}'
  params: {
    location: 'eastus' // New models not available in EU
    openAiServiceName: naming.openAiService
    keyVaultName: keyvault.outputs.keyVaultName
    logAnalyticsResourceId: lawbase.outputs.logAnalyticsId
  }
}

module speech 'speechService.bicep' = {
  name: 'speech-${buildId}'
  params: {
    keyVaultName: keyvault.outputs.keyVaultName
    speechServiceName: naming.speech
    location: location
  }
}

module data 'data.bicep' = {
  name: 'data-${buildId}'
  params: {
    dataCollectionEndpointName: naming.dataCollectionEndpoint
    dataCollectionRuleName: naming.dataCollectionRule
    keyVaultName: keyvault.outputs.keyVaultName
    location: location
    logAnalyticsName: lawbase.outputs.logAnalyticsName
    painTableShortName: naming.customTable
    managedIdentityObjectId: identity.outputs.userAssignedIdentityObjectId
    developerGroupObjectId: developerGroupObjectId
  }
}

module app 'app.bicep' = {
  name: 'app-${buildId}'
  params: {
    containerAppEnvName: naming.containerAppEnv
    containerAppName: naming.containerApp
    location: location
    logAnalyticsName: lawbase.outputs.logAnalyticsName
    dataCollectionSecretUris: data.outputs.secretUris
    openAiSecretUris: openAiService.outputs.secretUris
    speechSecretUris: speech.outputs.secretUris
    telegramSecretUris: keyvault.outputs.secretUris
    registryResourceId: registryResourceId
    userAssignedIdentityResourceId: identity.outputs.userAssignedIdentityResourceId
    userAssignedIdentityClientId: identity.outputs.userAssignedIdentityClientId
    containerTag: containerTag
  }
}
