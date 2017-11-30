#!/bin/bash
set -x

kubectl delete deployment -l app=cloud-controller-manager -n kube-system
kubectl delete service -l app=cloud-controller-manager -n kube-system

# Delete RBAC objects, if --rbac flag was used.
kubectl delete serviceaccount -l app=cloud-controller-manager -n kube-system
kubectl delete clusterrolebindings -l app=cloud-controller-manager -n kube-system
kubectl delete clusterrole -l app=cloud-controller-manager -n kube-system
