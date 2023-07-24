param logAnalyticsName string
param location string

resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2022-10-01' = {
  name: logAnalyticsName
  location: location
  properties: {
    sku: {
      name: 'PerGB2018'
    }
    retentionInDays: 60
  }
}

output logAnalyticsName string = logAnalytics.name
output logAnalyticsId string = logAnalytics.id
output logAnalyticsCustomerId string = logAnalytics.properties.customerId
