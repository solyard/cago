# Cert Manager CRL Manager

**Cert Manager Issuer CRL Manager ** - it's a simple tool that can integrate with Kubernetes Cert Manager installation and implement functionality that need to use Cert Manager as Certificate Authority with CRL support

## How it works ?

1. This tool implement the CRL generator, so every certificate that issued by Cert Manager using CertificateRequests (CRD) can bee revoked by simply additing annotation like this one:
```yaml
kind: CertificateRequest
metadata:
  annotations:
    "certificate/revoked": "true"
```

> Warning: ⚠️ If CertificateRequest was deleted then CRL will be generated without deleted certificate. Fix that check previous CRL will come later.

2. Tool use Kubernetes Service Account with permissions to list, read for CRD's and list, read and create for Secrets in Kubernetes namespace where it's running.

## How to install?

Just apply files that introduced in `docs/kubernetes` folder.
