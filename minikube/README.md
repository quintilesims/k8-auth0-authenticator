# Minikube Setup

https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication
https://stackoverflow.com/questions/46142990/configure-minikube-kubernetes-webhook-authentication-authorization

NOTE: K8s uses 'authn' and 'authz' to mean 'authentication' and 'authorization', respectively. 
* Authentication (authn): are the given credentials valid / which entity do the credentials belong to? 
* Authorization (authz): is the entity allowed to do the requested action? 

We are using the k8-auth0-authenticator for authentication and RBAC for authorization.


According to https://stackoverflow.com/questions/46142990/configure-minikube-kubernetes-webhook-authentication-authorization
"By default minikube mounts your Users directory therefore you can access config files over the /Users/username/path-to-file.yml",
so if we symlink our `authn.yaml` we shoudl be good:

```
# TODO: Update authn.yaml w/ certificate paths
  
# create a symlink for authn.yaml
ln -s $(pwd)/authn.yaml ~/authn.yaml

# start up minikube
minikube start --extra-config apiserver.Authentication.WebHook.ConfigFile=~/authn.yaml
```
