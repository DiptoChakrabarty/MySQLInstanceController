## Creating your operator
```
operator-sdk init --domain=dipto.mysql.example --repo=github.com/DiptoChakrabarty/MySQLInstanceController.git --owner "DiptoChakrabarty"
```

## Creating Operator API and Controller Files
```
operator-sdk create api --group dipto.mysql.example --version v1alpha1 --kind MySQLInstance --resource --controller

```

## Create the manifests
```
make manifests
```

## Create docker image and push it
```
make docker-build docker-push
```

## Deploy to cluster
```
make deploy
```