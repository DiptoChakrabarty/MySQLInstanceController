package controllers

import (
	"fmt"

	mysqlv1alpha1 "github.com/DiptoChakrabarty/MySQLInstanceController.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type BackupSchedule struct {
	BackupSchedule string
	MysqlName      string
	UserName       string
	Password       MysqlPasswords
	ServiceName    string
}

type MySQLInstanceConfig struct {
	Name      string
	Namespace string
	Instance  mysqlv1alpha1.MySQLInstance
}

type MysqlPasswords struct {
	RootPassword         string
	ClusterAdminPassword string
}

func NewMySQLStatefulSet(name string, namespace string, SecretName string) *appsv1.StatefulSet {
	// Define the StatefulSet Template
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql:8.0.26", // Replace with your MySQL image
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: SecretName,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return statefulSet
}

func NewMySQLSecret(name string, namespace string, password MysqlPasswords) *corev1.Secret {
	// Define the Secret Template
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"MYSQL_ROOT_PASSWORD": []byte(password.RootPassword),
			"MYSQL_USER":          []byte("clusteradmin"),
			"MYSQL_PASSWORD":      []byte(password.ClusterAdminPassword),
		},
	}

	return secret
}

func NewMySQLBackupCronJob(backupObject BackupSchedule, namespace string) *beta1.CronJob {
	// Defining the CronJOb Template
	mysqlDumpCommand := fmt.Sprintf(
		"mysqldump -h %s -u %s -p%s --all-databases > /backup/%s_backup.sql",
		backupObject.ServiceName,
		backupObject.UserName,
		backupObject.Password.ClusterAdminPassword,
		backupObject.MysqlName,
	)
	scheduleContainerName := fmt.Sprintf("%s-backup-container", backupObject.MysqlName)
	cronJob := &beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupObject.MysqlName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": backupObject.MysqlName,
			},
		},
		Spec: beta1.CronJobSpec{
			Schedule: backupObject.BackupSchedule,
			JobTemplate: beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            scheduleContainerName,
									Image:           "mysql:latest",
									ImagePullPolicy: corev1.PullAlways,
									Command: []string{
										"/bin/bash",
										"-c",
										mysqlDumpCommand,
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "backup-volume",
											MountPath: "/backup",
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "backup-volume",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return cronJob
}

func NewMySQLService(name string, namespace string) *corev1.Service {
	// Defining the Service template
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-service",
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       3306,
					TargetPort: intstr.IntOrString{IntVal: 3306},
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	return service
}
