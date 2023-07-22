package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BackupSchedule struct {
	BackupSchedule string
	MysqlName      string
	UserName       string
	Password       string
}

func NewMySQLStatefulSet(name string, namespace string, SecretName string) *appsv1.StatefulSet {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32(3),
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

func NewMySQLSecret(name string, namespace string, rootPwd string, clusteradminPwd string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"MYSQL_ROOT_PASSWORD": []byte(rootPwd),
			"MYSQL_USER":          []byte("clusteradmin"),
			"MYSQL_PASSWORD":      []byte(clusteradminPwd),
		},
	}

	return secret
}

func NewMySQLBackupCronJob(backupObject BackupSchedule, namespace string) *beta1.CronJob {

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
				Spec: beta1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            "backup-container",
									Image:           "mysql:latest",
									ImagePullPolicy: corev1.PullAlways,
									Command: []string{
										"/bin/bash",
										"-c",
										"mysqldump -h mysql-prd-0.mysql-prd -u <mysql-username> -p<mysql-password> <database-name> > /backup/db_backup.sql",
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
