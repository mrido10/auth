## Before Running
before run this code, please generate config.yaml file into config folder

```yaml
mySql:
  user: admin
  password: admin
  dbName: yourDB
server:
  host: localhost
  sqlPort: 3306
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
```