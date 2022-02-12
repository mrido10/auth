## Introduction
This application is authentication service that can currently be used in money manager applications 
in other repositories at [`github.com/mrido10`](https://github.com/mrido10). It is possible that it can be used for other 
applications by adding some permissions into database, so that other applications can use this application.

## Before Running
before run this code, generate config.yaml file into config folder

```yaml
server:
  host: localhost
  servicePort: 3003
  protocol: http
jwt:
  key: yourKey
hash:
  secret: yourSecret
sendEmail:
  SMTP_HOST: smtp.gmail.com
  SMTP_PORT: 587
  SENDER_NAME: Your Name
  AUTH_EMAIL: yourEmail@mail.com
  AUTH_PASSWORD: yourEmailPassword
encryptDecrypt:
  key: yourKey
postgres:
  address: host=localhost port=5432 user=postgres password=password dbname=dbname sslmode=disable
```

## API Doc
### Register
Registration can only use an active email, because after registration there will be an account 
activation process, so that registered accounts can log in
````http request
POST /register
````
````json
{
    "email": "activeEmail@gmail.com",
    "password": "password123",
    "rePassword": "password123",
    "name": "Your Name",
    "gender": "male",
    "accessID": 1
}
````

### Activation
````http request
PUT /activate?d=uXlgc3vvUxhtXJG6ae3uAoG0pmFRdiuCl-hl1LCjCeR4gyqcW5vOLXTt2psAmv0HQNqxp-PWUL20mPspI9tlBbeQIlJ2eMtCqubR8j_cBKI5ltQ%3D
````
| Parameter | Type | Description |
| :--- | :--- | :--- |
| `d` | `string` | **Required**. AES encrypted data |

### Resend Activation
````http request
POST /resendActivation
````
````json
{
    "email": "activeEmail@gmail.com"
}
````

### Login
````http request
POST /login
````
````json
{
    "email": "activeEmail@gmail.com",
    "password": "password123"
}
````
Response
````json
{
    "message": "Success",
    "status": true,
    "data": {
        "authorization": "<JWT token>"
    }
}
````