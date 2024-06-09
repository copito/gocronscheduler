MINIKUBE_STATUS := $(shell minikube status | grep -c Running)

install_packages:
	go get k8s.io/client-go@latest
	go get k8s.io/apimachinery@latest
	go get k8s.io/api@latest

minikube:
	if [ $(MINIKUBE_STATUS) -ge 3 ]; then 
		echo "Stopping Minikube" && minikube stop; 
	else 
		echo "Starting Minikube" && minikube start; 
	fi