# Crane Migration Demo

## Getting Started

To start the demo, we will need a Kubernetes cluster, this can be a local minikube cluster, or it could be a managed cluster. I choose to test with an EKS cluster. 

We will also need an [OpenShift cluster](https://www.redhat.com/en/technologies/cloud-computing/openshift/try-it).

## Install the *KS Micro Service

Once you have the cluster, we will need to install the microservice application 'Sock Shop'. You can find out more about it [here](https://github.com/microservices-demo/microservices-demo) but we have moved the manifests you need into this repo. 

```bash
$ kubectl create -f manifests
```

This will create a `sock-shop` namespace and the application. Once we make sure that the application is functioning we will now try to migrate and move the application to an OpenShift Cluster. 

## Exporting the Application 

To export the application to disk, we will use the `Crane Export` command. 

For this command to work, make sure that `KUBECONFIG` is pointing to the *KS cluster. 

```bash
$ crane export -n sock-shop
```

This will look for all the resources that are namespaced in the cluster, and then get the instances of those resources in the namespace.

Now that we have all the resources, we can make the transformations and get these resources ready to be installed in OpenShift. 

## Transformations

Using the Crane plugin manager, we will see that there is an OpenShift plugin. 

```bash
$ crane plugin-manager list
Listing from the repo default
+-------------------+-----------------+
| Name              | OpenShiftPlugin |
| ShortDescription  | OpenShiftPlugin |
| AvailableVersions | v0.0.2, v0.0.3  |
+-------------------+-----------------+
```

We will install this plugin, as someone has done the hard work of making decisions on what should be changed in a Kubernetes resource to move to OpenShift.

To Install the plugin:

```bash
$ crane plugin-manager add OpenShiftPlugin --version v0.0.3 --plugin-dir=~/.crane/plugin
```

Now that the Plug-In is installed let's go ahead and run our transform.

```bash
$ crane transform --plugin-dir=~/.crane/plugin
```

we will see a transform directory that will contain the changes that will be applied to export, such that you can get the output that can be created in OpenShift.

But while it ran, we should have seen some errors such as: 

```
security context for container: {container-name} has security context. Updates may be needed
```

### Dealing With Security Constraints

There are many ways that you could fix this warning, but I choose to remove the security context. I could have updated each file in the export, but what if we need to do this again! I am going to use the default security constraints in the [OpenShift](https://docs.openshift.com/container-platform/4.10/authentication/managing-security-context-constraints.html) platform.

To achieve this, we will use the ```test-python-plugin.py``` plugin. to use this, we will download it, make sure it is executable, and install it into the plugins directory. 

```bash
$ curl {path_to_python_raw} > ~/.crane/plugin/test-python-plugin.py && chmod +x ~/.crane/plguin/test-pythong-plugin.py
```

Once this is done, let's go ahead and re-run crane transform, the warnings may still be present, but we know that we are taking care of this!

### Filesystem Access Transformation
Now, this is where knowing your application will be very helpful, this may require trial and error, but I know that I need some filesystem access for the root filesystem, for the databases to work. To fix this particular problem I wrote a plugin in go. You will need to download, build the binary and install this into the plugins directory.

You will then need to re-run crane transform.

Once this happens you will see, for the `carts-db`, `orders-db`, `user-db` transformations file will also be adding an emptyDir volume for the databases to use!


## Apply

Now that we are ready, let's go ahead and apply the transformations.

```bash
$ crane apply
```

Now we will switch the kubeconfig, to point to the OpenShift cluster, and apply the resource, we are almost done with the migration!

```bash
$ kubectl create -f ./output
```

Well then, it seems everything worked! and everything should be up and running, but we did see those deprecation warnings. 
```
beta.k8s.io/os is deprecated ... use k8s.io/os
```
We probably want to future-proof the automation. For this, you will need to download the `test-ruby-plugin.rb` and install it into the plugin directory.

```bash
$ curl {path_to_ruby_raw} > ~/.crane/plugin/test-ruby-plugin.rb && chmod +x ~/.crane/plugin/test-ruby-plugin.rb
```

I would personally re-run ```crane transform -p ~/.crane/plugins && crane apply && kubectl replace -f ./output```


