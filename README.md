# MySQLInstanceController

# Project Description

The MySQLInstance Controller is a Kubernetes Operator designed to automate the deployment and management of MySQL instances in a Kubernetes cluster. It simplifies the process of setting up and maintaining MySQL databases with backups , providing a seamless experience for developers and administrators. 

## Flow of how controller works

<details><summary>show</summary>
<p>

```bash
- Custom Resource Definition (CRD):
The project defines a new custom resource named MySQLInstance. This custom resource allows users to specify the configuration details of their MySQL instances, such as version, storage requirements, and backup preferences.

- Reconciliation Logic:
The core of the Operator is the reconciliation loop, which continuously ensures that the actual state of the system matches the desired state defined in the MySQLInstance custom resource's spec. In each iteration, the Operator compares the current state of the resources with the desired state and takes corrective actions as needed.

- StatefulSet Creation:
When a new MySQLInstance custom resource is created, the Operator generates a StatefulSet manifest based on the specifications provided in the resource's spec. The StatefulSet defines the MySQL pods and their persistent storage.

- Service Creation:
The Operator creates a Service to expose the MySQL pods within the Kubernetes cluster. The Service allows other applications to interact with the MySQL database using the appropriate endpoint and port.

- Secrets Management:
For secure communication with MySQL pods, the Operator creates and manages Kubernetes Secrets containing credentials, such as the MySQL root password. These Secrets are used by the MySQL pods for authentication.

- Backup Scheduler:
The MySQLInstance custom resource allows users to specify whether backups are required for their MySQL instances. If backups are enabled, the Operator sets up a backup schedule using Kubernetes CronJobs.

- Backup Logic:
When the backup schedule is triggered, the Operator orchestrates the backup process for the MySQL instances. It interacts with the MySQL pods, using the credentials from the Secrets, and performs a backup using either mysqldump or a custom backup tool.
```

</p>
</details>
