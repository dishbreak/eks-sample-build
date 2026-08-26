# EKS Auto Mode Cluster

This is an example EKS cluster, intended for studying and practice use only. I tried to make the best possible decisions I could, but this is definitely not a production grade setup.

> [!WARNING]
> Applying this code in your environment can cost you **real money**.
> Make sure to destroy it when it's not actively in use.

## What's the goal?

As a DevOps Engineer, I want to be able to specify a full Kubernetes environment end-to-end. This includes:

* Specifying the underlying cloud infrastructure to create the cluster
* Deploying a GitOps solution to manage applications on the cluster
* Creating a build/release pipeline that will automatically deploy successful builds of software to the cluster and open preview deploys for software. 

## What's in Here? 

* `eks_cluster/` - Infrastructure-as-code to handle creating an EKS cluster using Auto Mode.
* `argocd/` - Declarative YAML setup to manage an ArgoCD instance on said cluster

## What's NOT in Here?

There's a number of things that I'd want to see in a production-ready environment that simply aren't here.

* Integration with a Federated Identity System (Okta, Active Directory, etc.) 
* Role-based access control to Kubernetes 
* Use of managed database solutions like RDS
* Security defenses like supply-chain monitoring, intrusion protection/detection, and vulnerability scanning. 

These are all things I'd expect someone making money off of their workloads to have. As a solo practitioner, it's cost prohibitive and infeasible to implement much of this for infrastructure designed to be stood up and torn down in the span of an evening.
