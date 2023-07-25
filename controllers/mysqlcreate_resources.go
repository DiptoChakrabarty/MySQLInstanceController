package controllers

import (
	"context"
	"math/rand"
	"time"
)

// Create StatefulSet method
func (rtx *MySQLInstanceReconciler) CreateMySQLStatefulset(mysqlInstanceConfig MySQLInstanceConfig) error {
	name := mysqlInstanceConfig.Name
	nameSpace := mysqlInstanceConfig.Namespace
	secretName := name + "-secret"

	// Create StatefulSet
	statefulset := NewMySQLStatefulSet(name, nameSpace, secretName)
	err := rtx.Create(context.TODO(), statefulset)
	if err != nil {
		return err
	}

	return nil
}

// Create the mysql secret
func (rtx *MySQLInstanceReconciler) CreateMySQLSecret(mysqlInstanceConfig MySQLInstanceConfig) error {
	name := mysqlInstanceConfig.Name
	nameSpace := mysqlInstanceConfig.Namespace
	// Generate a random password for the MySQL root user
	rand.Seed(time.Now().UnixNano())
	rootPwd := generateRandomPassword()
	clusteradminPwd := generateRandomPassword()

	// Create a new secret for the mysql statefulset
	secretName := name + "-secret"
	secret := NewMySQLSecret(secretName, nameSpace, rootPwd, clusteradminPwd)

	// Create the Secret
	err := rtx.Create(context.TODO(), secret)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}

// Create the mysql service
func (rtx *MySQLInstanceReconciler) CreateMySQLService(mysqlInstanceConfig MySQLInstanceConfig) error {
	name := mysqlInstanceConfig.Name
	nameSpace := mysqlInstanceConfig.Namespace

	// Create the mysql service
	service := NewMySQLService(name, nameSpace)

	// Create service
	err := rtx.Create(context.TODO(), service)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}

// Create the mysql CronJob
func (rtx *MySQLInstanceReconciler) CreateMySQLCronJOB(mysqlInstanceConfig MySQLInstanceConfig, backupObject BackupSchedule) error {
	nameSpace := mysqlInstanceConfig.Namespace

	// Create the mysql cronjob
	cronJob := NewMySQLBackupCronJob(backupObject, nameSpace)

	// Create cronJob
	err := rtx.Create(context.TODO(), cronJob)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}
