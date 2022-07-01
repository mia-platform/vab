# Folder structure

Example:

    ./project-name
    ├── vendors
    |   ├── modules
    |   |   ├── ingress
    |   |   |   └── traefik
    |   |   |       ├── [...]
    |   |   |       └── kustomization.yaml
    |   |   └── cni
    |   |       ├── cilium
    |   |       |   ├── [...]
    |   |       |   └── kustomization.yaml
    |   |       └── calico
    |   |           ├── [...]
    |   |           └── kustomization.yaml
    |   └── add-ons
    |       └── ..
    ├── overrides
    |   └── group-name
    |       ├── all-clusters
    |       |   ├── ingress
    |       |   |   └── traefik
    |       |   ├── cni
    |       |   |   └── cilium
    |       |   └── kustomization.yaml
    |       ├── cluster1-name
    |       |   ├── cni
    |       |   |   └── calico
    |       |   └── kustomization.yaml    
    |       └── cluster2-name
    |           └── kustomization.yaml 
    └── config.yaml
