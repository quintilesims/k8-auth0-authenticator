# Kubernetes Auth0 Authenticator
This repo contains instructions for setting up [OIDC](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) with Auth0 using Minikube. 

## Step 1: Create Auth0 Application
* Type: `Regular Web Application`
* Allowed Callback URLs: `http://localhost:8000/`

Setup environment variables relating to your Auth0 account (they will be used later).

```console
$ export AUTH0_DOMAIN=https://imshealth.auth0.com/
$ export AUTH0_CLIENT_ID=<client id here>
$ export AUTH0_CLIENT_SECRET=<client secret here>
$ export AUTH0_EMAIL=<email used for Auth0 login here>
```

## Step 2: Start Minikube Cluster
You will need to configure your Kubernetes cluster to use OIDC authentication.
The flags are documented [here](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#configuring-the-api-server).

_Note: The `--vm-driver=hyperkit` flag is recommended, but not required._

```console
$ minikube start \
--vm-driver=hyperkit \
--extra-config=apiserver.Authorization.Mode=RBAC \
--extra-config=apiserver.Authentication.OIDC.IssuerURL=$AUTH0_DOMAIN \
--extra-config=apiserver.Authentication.OIDC.UsernameClaim=email \
--extra-config=apiserver.Authentication.OIDC.ClientID=$AUTH0_CLIENT_ID

Starting local Kubernetes v1.9.4 cluster...
Starting VM...
Getting VM IP address...
Moving files into cluster...
Setting up certs...
Connecting to cluster...
Setting up kubeconfig...
Starting cluster components...
Kubectl is now configured to use the cluster.
Loading cached images from config file.
```

### Step 3: Setup Kubeconfig
First, create a user with the same email you use when logging into Auth0:

```console
$ kubectl config set-context minikube --user $AUTH0_EMAIL
Context "minikube" modified.
```

Then, set part of the OIDC configuration for your user (the `id-token` and `refresh-token` values will be populated later):

```console
$ kubectl config set-credentials $AUTH0_EMAIL \
  --auth-provider oidc \
  --auth-provider-arg idp-issuer-url=$AUTH0_DOMAIN \
  --auth-provider-arg client-id=$AUTH0_CLIENT_ID \
  --auth-provider-arg client-secret=$AUTH0_CLIENT_SECRET \
  --auth-provider-arg extra-scopes=offline_access
  
User "<your Auth0 login email here>" set.
```

Finally, we need to populate the `id-token` and `refresh-token` values for your user. 
Download and execute [kubelogin](https://github.com/int128/kubelogin).
This tool will open up a web browser, prompt you to login to Auth0, and populate your `Kubeconfig` with the required `id-token` and `refresh-token` values.
Make sure to wait until the progam finishes executing (it may take about a minute).

```console
$ kubelogin
2018/09/18 14:47:44 Reading ~/.kube/config
2018/09/18 14:47:44 Using current-context: minikube
2018/09/18 14:47:44 Open http://localhost:8000 for authorization
2018/09/18 14:47:44 GET /
2018/09/18 14:47:49 GET /?code=XXXX&state=YYYY
2018/09/18 14:49:08 Got token for subject=ZZZZ
2018/09/18 14:49:08 Updated ~/.kube/config
```

This is what a valid `Kubeconfig` file should look like:
```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority: /Users/<user>/.minikube/ca.crt
    server: https://<ip>:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: <your auth0 login email>
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    as-user-extra: {}
    client-certificate: /Users/<user>/.minikube/client.crt
    client-key: /Users/<user>/.minikube/client.key
- name: <your auth0 login email>
  user:
    auth-provider:
      config:
        client-id: <auth0 client id>
        client-secret: <auth0 client secret>
        extra-scopes: offline_access
        id-token: <your auth0 id token>
        idp-issuer-url: https://imshealth.auth0.com/
        refresh-token: <your auth0 refresh token>
      name: oidc
```

## Step 4: Grant User Permissions
Now that you've created a user, you need to grant that user permissions.
Update the file `cluster_admin.yaml` in this directory to use your Auth0 login email.
```yaml
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
  name: <your auth0 login email>
```

Run the following to grant your user admin permissions:
```console
$ kubectl --user=minikube apply -f cluster_admin.yaml
TODO: show output
```

You should now be able to use `kubectl` with your user:
```console
$ kubectl get nodes
NAME       STATUS    ROLES     AGE       VERSION
minikube   Ready     <none>    1h        v1.9.4
```

TODO: If you get XXXXX, error message, try deleting `id-token`. 
