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

## Deploy Controller to cluster
```
make deploy
```

## Run a sample mysqlInstance example
```
kubectl apply -f mysql_instance_example.yaml
```

### Contributing
When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change.

If anyone wants to take up an issuse they are free to do so .


### Contribution Practices
Please be respectful of others , do not indulge in unacceptable behaviour.
If a person is working or has been assigned an issue and you want to work on it please ask him/her if he is working on it.
We are happy to allow you to work on your issues , but in case of long period of inactivity the issue will be approved to another volunteer.
If you report a bug please provide steps to reproduce the bug.
In case of changing the backend routes please submit an updated routes documentation for the same.
If there is an UI related change it would be great if you could attach a screenshot with the resultant changes so it is easier to review for the maintainers.