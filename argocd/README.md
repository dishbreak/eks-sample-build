# argocd install

## Preparing for the Install

This setup creates an IngressClass to enable the LoadBalancer controller.
Additionally, it creates an Ingress to expose the ArgoCD UI on an internal Load Balancer.
If you're going to use this setup for yourself, you'll need to make sure you generate an ACM certificate that matches a domain name you control and update the Ingress in `ui-service.yaml` with that same domain. 
Additionally, it's a good idea to have a Hosted Zone in Route53 that corresponds with the domain name, so that you're able to set A Alias records.

On a blank EKS cluster, use Kustomize via the `-k` flag on `kubectl` to apply the Kustomization that will install the latest version of ArgoCD on the cluster.
Note that the `--server-side` flag is needed to apply the Kustomization with [server-side apply](https://kubernetes.io/docs/reference/using-api/server-side-apply/). 
This is required because the resources will exceed the length limitations for client-side apply. 

```
k apply -n argocd -k . --server-side --force-conflicts
```

Once the apply is complete, manually create an A Alias record pointing at the Load Balancer created from the Ingress record.

## Connecting Without a Load Balancer

The setup creates a Load Balancer to allow you to interact with the HTTP/gRPC interfaces of ArgoCD. 
For some installations, it may not make sense to incur the cost of a Load Balancer.
As an alternative, use port forwarding to coneect to the `argocd-server` pod.

```
k port-forward -n argocd service/argocd-server 8080:80
```

You're now able to visit the ArgoCD UI on http://localhost:8080

## Logging in as the Admin User

You'll need to log in as the `admin` user. Argocd helpfully stores the initial admin password as a secret in the argocd namespace:

```
k get secret -o yaml argocd-initial-admin-secret
```

Note that the password itself is base64 encoded. On a macOS system with yq installed, this one-liner copies the decoded password straight to the clipboard:

```
k get secret -o yaml argocd-initial-admin-secret | yq '.data.password' | base64 -D | pbcopy
```
