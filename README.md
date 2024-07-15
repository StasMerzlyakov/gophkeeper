# gophkeeper - Yandex Practicum graduation project

## iteration1 - authirization

Client <-> Server:
- gRPC interaction  (based on tls, CA authority -> Server cert, start client with CA certificate)
- registration
- authorization (OTP)
- authentfication (JWT)
- TUI

Refs:
[AES-256 key encryption](https://github.com/andrewromanenco/gcrypt)

### Registration

prview в VSCode - Alt-D.

```plantuml
@startuml
Client->Client: Promt Authentification password
Client->GoKeeper: Registration Request (Email: string, Password: string)
GoKeeper-> EmailServer: Send OTP QR
EmailServer-->GoKeeper:
GoKeeper-->Client: Prompt OTP password
Client->GoKeeper: OTP pass
GoKeeper-->Client: Registration complete, JWT
@enduml
```


### Authorization
```plantuml
@startuml
Client->GoKeeper: Authorization Request(Email, Password)
GoKeeper-->Client: Prompt OTP password
Client->GoKeeper: OTP pass
GoKeeper-->Client: JWT
@enduml
```


### MasterKey Generation
```plantuml
@startuml
Client->Client: Prompt MasterKey Password and passwordHint (MasterPass, Hint)
Client->Client: Generate Storage AES-256 Key (SKey)
Client->Client: Encrypt SKey on MasterPass (EncrSKey)
Client->GoKeeper: Add MasterKey Request (JWT, Hint, EncrSKey)
GoKeeper-->Client: 

@enduml
```


