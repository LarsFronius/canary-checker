## <img src='https://raw.githubusercontent.com/flanksource/flanksource-ui/main/src/icons/http.svg' style='height: 32px'/> HTTP

??? example
     ```yaml
     apiVersion: canaries.flanksource.com/v1
     kind: Canary
     metadata:
       name: http-pass
     spec:
       interval: 30
       http:
         - endpoint: http://status.savanttools.com/?code=200
           thresholdMillis: 3000
           responseCodes: [201, 200, 301]
           responseContent: ""
           maxSSLExpiry: 7
         - endpoint: http://status.savanttools.com/?code=404
           thresholdMillis: 3000
           responseCodes: [404]
           responseContent: ""
           maxSSLExpiry: 7
         - endpoint: http://status.savanttools.com/?code=500
           thresholdMillis: 3000
           responseCodes: [500]
           responseContent: ""
           maxSSLExpiry: 7
     
     ```

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| authentication | Credentials for authentication headers | *[Authentication](#authentication) |  |
| body | Request Body Contents | string |  |
| description | Description for the check | string |  |
| display |  | [Template](#template) |  |
| **endpoint** | HTTP endpoint to check.  Mutually exclusive with Namespace | string | Yes |
| headers | Header fields to be used in the query | \[\][kommons.EnvVar](https://pkg.go.dev/github.com/flanksource/kommons#EnvVar) |  |
| icon | Icon for overwriting default icon on the dashboard | string |  |
| maxSSLExpiry | Maximum number of days until the SSL Certificate expires. | int |  |
| method | Method to use - defaults to GET | string |  |
| name | Name of the check | string |  |
| namespace | Namespace to crawl for TLS endpoints.  Mutually exclusive with Endpoint | string |  |
| ntlm | NTLM when set to true will do authentication using NTLM v1 protocol | bool |  |
| ntlmv2 | NTLM when set to true will do authentication using NTLM v2 protocol | bool |  |
| responseCodes | Expected response codes for the HTTP Request. | \[\]int |  |
| responseContent | Exact response content expected to be returned by the endpoint. | string |  |
| responseJSONContent | Path and value to of expect JSON response by the endpoint | [JSONCheck](#jsoncheck) |  |
| test |  | [Template](#template) |  |
| thresholdMillis | Maximum duration in milliseconds for the HTTP request. It will fail the check if it takes longer. | int |  |