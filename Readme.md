# Argocd Inventory
This app creates a simple inventory of the applications and their associated resources as defined in an ArgoCD instance.

# TODO
- [ ] Add a way to check for kubernetes API changes based on an initial version and a target version

# Usage
```bash
Usage of argo-inventory:
  -argo-password string
    	The ArgoCD password to use.
  -argo-url string
    	The ArgoCD URL to use.
  -argo-username string
    	The ArgoCD username to use.
  -log-level string
    	The log level to use. (default "info")
  -out string
    	The output file to write the report to. (default "report.csv")
```
