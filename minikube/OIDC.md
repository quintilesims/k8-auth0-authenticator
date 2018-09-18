# OIDC

Instructions on setting up Auth0 as an OIDC provider for Minikube.

### Step 1: Setup Auth0
TODO: Instructions on setting up the Auth0 application.

(grab from `k8-auth0-authenticator` app: https://manage.auth0.com/#/applications/NBfx2YgzADXGcpw79k0zZbbCJ5iIoAGf/settings)
```
export AUTH0_DOMAIN=https://imshealth.auth0.com/
export AUTH0_CLIENT_ID=<client id here>
export AUTH0_CLIENT_SECRET=<client secret here>
export AUTH0_EMAIL=<your email here>
```

### Step 2: Start Minikube
```console
$ minikube start \
--vm-driver=hyperkit \
--extra-config=apiserver.Authorization.Mode=RBAC \
--extra-config=apiserver.Authentication.OIDC.IssuerURL=$AUTH0_DOMAIN \
--extra-config=apiserver.Authentication.OIDC.UsernameClaim=email \
--extra-config=apiserver.Authentication.OIDC.ClientID=$AUTH0_CLIENT_ID
``` 

### Step 3: Configure Kubeconfig 

```
```
Use:  https://github.com/int128/kubelogin



TODO: How to get ID token and refresh token from Auth0? 
(openid connect? cli helper?)

```
$ kubectl config set-context minikube --user $AUTH0_EMAIL

kubectl config set-credentials $AUTH0_EMAIL \
  --auth-provider oidc \
  --auth-provider-arg idp-issuer-url=$AUTH0_DOMAIN \
  --auth-provider-arg client-id=$AUTH0_CLIENT_ID \
  --auth-provider-arg client-secret=$AUTH0_CLIENT_SECRET \
  --auth-provider-arg extra-scopes=offline_access

$ kubelogin
```

### Step 4: Create RoleBinding for your User
```
kubectl --user=minikube apply -f cluster_admin.yaml 
```

TODO:
So right now, you need to delete `id-token` after `kubelogin`. 
The refresh token works, so it should be good. 
This also is working after `zack.patrick@appature.com` as used in the cluster_admin.yaml,
and the user ID   `google-oauth2|113669605019273744908` was used in `~/.kube/config`.

Trying with `zack.patrick@apature.com` instead of userid.
ok, that seemed to work.  
