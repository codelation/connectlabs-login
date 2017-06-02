Travis CI Status [![Build Status](https://travis-ci.org/codelation/connectlabs-login.svg?branch=master)](https://travis-ci.org/codelation/connectlabs-login)

### Notes


#### Regions

| URL                              | Description                               |
| -------------------------------- | ----------------------------------------- |
| `/auth/{provider}/callback`      | Callback endpoint for OAuth 2.0           |
| `/auth/{provider}/login`         | URL for login; automatic or redirect      |
| `/auth/{provider}/logout`        | URL for logging user out of provider      |
| `/auth/login.html`               | Login page for router to redirect to      |
| `/ap/auth.html`                  | Internal endpoint for access point        |
| `/api/users/{mac}`               | Get a User by one of it's MAC addresses   |


### Template

The template used on the login page uses the following view model:

```json
{
  "TwitterProvider": true,
  "TwitterUser": {
    "RawData": null,
    "Provider": "",
    "Email": "",
    "Name": "",
    "FirstName": "",
    "LastName": "",
    "NickName": "",
    "Description": "",
    "UserID": "",
    "AvatarURL": "",
    "Location": "",
    "AccessToken": "",
    "AccessTokenSecret": "",
    "RefreshToken": "",
    "ExpiresAt": "0001-01-01T00:00:00Z"
  },
  "FacebookProvider": true,
  "FacebookUser": {
    "RawData": null,
    "Provider": "",
    "Email": "",
    "Name": "",
    "FirstName": "",
    "LastName": "",
    "NickName": "",
    "Description": "",
    "UserID": "",
    "AvatarURL": "",
    "Location": "",
    "AccessToken": "",
    "AccessTokenSecret": "",
    "RefreshToken": "",
    "ExpiresAt": "0001-01-01T00:00:00Z"
  },
  "GPlusProvider": true,
  "GPlusUser": {
    "RawData": null,
    "Provider": "",
    "Email": "",
    "Name": "",
    "FirstName": "",
    "LastName": "",
    "NickName": "",
    "Description": "",
    "UserID": "",
    "AvatarURL": "",
    "Location": "",
    "AccessToken": "",
    "AccessTokenSecret": "",
    "RefreshToken": "",
    "ExpiresAt": "0001-01-01T00:00:00Z"
  },
  "SSID": "",
  "Title": "",
  "SubTitle": "",
  "Email": "",
  "Message": ""
}
```
