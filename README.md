# LemBaas Go SDK

Go reference: [https://pkg.go.dev/github.com/KaiserWerk/lembaas-go-sdk](https://pkg.go.dev/github.com/KaiserWerk/lembaas-go-sdk)

## App authentication

Use `ClientID` and `ClientSecret` to authenticate your app and get an access token. The access token in turn can be used to make authenticated requests.

The following example shows how to authenticate an app and get an access token using the `AppClient`. The same access token can then be used to make authenticated requests with the other clients.

```go
using lembaas "github.com/KaiserWerk/lembaas-go-sdk"

ctx := context.Background()
baseURL := "https://api.lembaas.net"
apiVersion := 1
clientID := "your-client-id"
clientSecret := "your-client-secret"

appClient := lembaas.NewAppClient(baseURL, apiVersion)
appTokenResponse, _ := appClient.Authenticate(ctx, clientID, clientSecret)

fmt.Printf("Access Token: %s\n", appTokenResponse.Token)
```

## Clients

Use the different client implementations with the respective API version to make requests to the `LemBaas API`.

### AppClient

```go
appResponse, _ := appClient.GetAppInfo(ctx, appTokenResponse.Token)
fmt.Printf("App Name: %s\n", appResponse.Name)
```

### AppConfigClient

```go
configClient := lembaas.NewConfigClient(baseURL, apiVersion, token)
configResponse, _ := configClient.ListCustomConfigValues(ctx)
```

### RoleClient

```go
roleClient := lembaas.NewRoleClient(baseURL, apiVersion, token)
roleResponse, _ := roleClient.ListRoles(ctx)
```

### UserClient

```go
userClient := lembaas.NewUserClient(baseURL, apiVersion, token)

newUserRequest := lembaas.RegisterUserRequest{
    Email: "newuser@example.com",
    Password: "password123",
}
newuserResponse, _ := userClient.RegisterUser(ctx, newUserRequest)

userListResponse, _ := userClient.ListUsers(ctx)
```

## Utility functions

```go
var userAuthResponse *AppUserAuthResponse = userClient.LoginUser(ctx, &lembaas.AppUserAuthRequest{...})

// Tells whether the user authentication response indicates a valid login
validLogin := lembaas.IsValidLogin(userAuthResponse)
    
// Tells whether the user authentication response indicates that TOTP is required
isTOTPRequired := lembaas.IsTOTPRequired(userAuthResponse)
```