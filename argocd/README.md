# argocd install

On a blank EKS cluster, use Kustomize via the `-k` flag on `kubectl` to apply the Kustomization that will install the latest version of ArgoCD on the cluster.
Note that the `--server-side` flag is needed to apply the Kustomization with [server-side apply](https://kubernetes.io/docs/reference/using-api/server-side-apply/). 
This is required because the resources will exceed the length limitations for client-side apply. 

```
k apply -n argocd -k . --server-side --force-conflicts
```

In order to connect, use port forwarding to coneect to the `argocd-server` pod.

```
k port-forward -n argocd service/argocd-server 8080:80
```

You'll need to log in as the `admin` user. Argocd helpfully stores the initial admin password as a secret in the argocd namespace:

```
k get secret -o yaml argocd-initial-admin-secret
```

Note that the password itself is base64 encoded. On a macOS system with yq installed, this one-liner copies the decoded password straight to the clipboard:

```
k get secret -o yaml argocd-initial-admin-secret | yq '.data.password' | base64 -D | pbcopy
```
