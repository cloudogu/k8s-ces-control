# k8s-ces-control entwickeln

## Wie werden Service Accounts angelegt?
- Dogu-Operator execed in k8s-ces-control als CLI

## GRPC-Services generieren
Die GRPC-Services werden im Verzeichnis `grpc-protobuf` beschrieben. 
Mit dem Make-Target `make generate-grpc` können die Stubs für die GRPC-Services neu generiert werden.
Die generierten Sourcen sind im Verzeichnis `generated` zu finden. 