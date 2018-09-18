# OIDC

Instructions on setting up Auth0 as an OIDC provider for Minikube.

### Step 1: Create Auth0 Application
TODO: Instructions on setting up the Auth0 application.

### Step 2: Run Minikube
```console
$ minikube start \
--extra-config=apiserver.Authorization.Mode=RBAC \
--extra-config=apiserver.Authentication.OIDC.IssuerURL=https://imshealth.auth0.com/ \
--extra-config=apiserver.Authentication.OIDC.UsernameClaim=email \
--extra-config=apiserver.Authentication.OIDC.ClientID="AUTH0_CLIENT_ID_HERE"
``` 

### Step 3: Get ID Token
TODO: How to get ID token and refresh token from Auth0? 
(openid connect? cli helper?)

```
kubectl config set-credentials EMAIL_HERE \
  --auth-provider oidc \
  --auth-provider-arg idp-issuer-url=https://imshealth.auth0.com/ \
  --auth-provider-arg client-id=CLIENT_ID_HERE \
  --auth-provider-arg client-secret=CLIENT_SECRET_HERE
  --auth-provider-arg id-token=ID_TOKEN_HERE \
  --auth-provider-arg refresh-token=REFRESH_TOKEN_HERE
```

### Step 4: Create RoleBinding for your User
```
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: oidc-admin-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: User
  name: EMAIL_HERE
```

#### Step 5: Set user for context
```console
$ kubectl config set-context minikube --user=EMAIL_HERE
```
