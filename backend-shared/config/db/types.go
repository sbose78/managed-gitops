package db

import "time"

// GitopsEngineCluster is used to track clusters that host Argo CD instances
type GitopsEngineCluster struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"gitopsenginecluster"` //nolint

	Gitopsenginecluster_id string `pg:"gitopsenginecluster_id,pk"`

	SeqID int64 `pg:"seq_id"`

	// -- pointer to credentials for the cluster
	// -- Foreign key to: ClusterCredentials.clustercredentials_cred_id
	Clustercredentials_id string `pg:"clustercredentials_id"`
}

// GitopsEngineInstance is an Argo CD instance on a Argo CD cluster
type GitopsEngineInstance struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"gitopsengineinstance,alias:gei"` //nolint

	Gitopsengineinstance_id string `pg:"gitopsengineinstance_id,pk"`
	SeqID                   int64  `pg:"seq_id"`

	// -- An Argo CD cluster may host multiple Argo CD instances; these fields
	// -- indicate which namespace this specific instance lives in.
	Namespace_name string `pg:"namespace_name"`
	Namespace_uid  string `pg:"namespace_uid"`

	// -- Reference to the Argo CD cluster containing the instance
	// -- Foreign key to: GitopsEngineCluster.gitopsenginecluster_id
	EngineCluster_id string `pg:"enginecluster_id"`
}

// ManagedEnvironment is an environment (namespace(s) on a user's cluster) that they want to deploy applications to, using Argo CD
type ManagedEnvironment struct {
	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"managedenvironment,alias:me"` //nolint

	Managedenvironment_id string `pg:"managedenvironment_id,pk"`
	SeqID                 int64  `pg:"seq_id"`

	// -- human readable name
	Name string `pg:"name"`

	// -- pointer to credentials for the cluster
	// -- Foreign key to: ClusterCredentials.clustercredentials_cred_id
	Clustercredentials_id string `pg:"clustercredentials_id"`
}

// ClusterCredentials contains the credentials required to access a K8s cluster.
// The credentials may be in one of two forms:
// 1) Kubeconfig state: Kubeconfig file, plus a reference to a specific context within the
//     - This is the same content as can be found in your local '~/.kube/config' file
//     - This is what the user would initially provide via the Service/Web UI/CLI
//     - There may be (likely is) a better way of doing this, but this works for now.
// 2) ServiceAccount state: A bearer token for a service account on the target cluster
//     - Same mechanism Argo CD users for accessing remote clusters
//
// You can tell which state the credentials are in, based on whether 'serviceaccount_bearer_token' is null.
//
// It is the job of the cluster agent to convert state 1 (kubeconfig) into a service account
// bearer token on the target cluster (state 2).
//     - This is the same operation as the `argocd cluster add` command, and is the same
//       technique used by Argo CD to interface with remove clusters.
//     - See https://github.com/argoproj/argo-cd/blob/a894d4b128c724129752bac9971c903ab6c650ba/cmd/argocd/commands/cluster.go#L116
type ClusterCredentials struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"clustercredentials"` //nolint

	// -- Primary key for the credentials (UID)
	Clustercredentials_cred_id string `pg:"clustercredentials_cred_id,pk"`

	SeqID int64 `pg:"seq_id"`

	// -- API URL for the cluster
	// -- Example: https://api.ci-ln-dlfw0qk-f76d1.origin-ci-int-gce.dev.openshift.com:6443
	Host string `pg:"host"`

	// -- State 1) kube_config containing a token to a service account that has the permissions we need.
	Kube_config string `pg:"kube_config"`

	// -- State 1) The name of a context within the kube_config
	Kube_config_context string `pg:"kube_config_context"`

	// -- State 2) ServiceAccount bearer token from the target manager cluster
	Serviceaccount_bearer_token string `pg:"serviceaccount_bearer_token"`

	// -- State 2) The namespace of the ServiceAccount
	Serviceaccount_ns string `pg:"serviceaccount_ns"`
}

// ClusterUser is an individual user/customer
// Note: This is basically placeholder: a real implementation would need to be way more complex.
type ClusterUser struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"clusteruser"` //nolint

	Clusteruser_id string `pg:"clusteruser_id,pk"`
	User_name      string `pg:"user_name"`
	SeqID          int64  `pg:"seq_id"`
}

type ClusterAccess struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"clusteraccess"` //nolint

	// -- Describes whose cluster this is (UID)
	Clusteraccess_user_id string `pg:"clusteraccess_user_id,pk"`

	// -- Describes which managed environment the user has access to (UID)
	Clusteraccess_managed_environment_id string `pg:"clusteraccess_managed_environment_id,pk"`

	// -- Which Argo CD instance is managing the cluster?
	// clusteraccess_gitops_engine_instance_id VARCHAR (48) UNIQUE,
	Clusteraccess_gitops_engine_instance_id string `pg:"clusteraccess_gitops_engine_instance_id,pk"`

	// -- CONSTRAINT fk_cluster_access_target_inf_cluster   FOREIGN KEY(cluster_access_target_inf_cluster)  REFERENCES InfrastructureCluster(inf_cluster_id),
	// -- CONSTRAINT fk_cluster_access_user_id   FOREIGN KEY(cluster_access_user_id)  REFERENCES ClusterUser(user_id),

	// PRIMARY KEY(clusteraccess_user_id, clusteraccess_managed_environment_id, clusteraccess_gitops_engine_instance_id)

	SeqID int64 `pg:"seq_id"`
}

type OperationState string

const (
	OperationState_Waiting     = "Waiting"
	OperationState_In_Progress = "In_Progress"
	OperationState_Completed   = "Completed"
	OperationState_Failed      = "Failed"
)

type Operation struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"operation,alias:op"` //nolint

	// -- UID
	Operation_id string `pg:"operation_id,pk"`

	SeqID int64 `pg:"seq_id"`

	// -- Specifies which Argo CD instance is this operation against
	// -- Foreign key to: GitopsEngineInstance.gitopsengineinstance_id
	Instance_id string `pg:"instance_id"`

	// -- CONSTRAINT fk_target_inf_cluster   FOREIGN KEY(target_inf_cluster)  REFERENCES InfrastructureCluster(inf_cluster_id),

	// -- UID of the resource that was updated
	Resource_id string `pg:"resource_id"`

	// -- The user that initiated the operation.
	Operation_owner_user_id string `pg:"operation_owner_user_id"`

	// -- Resource type of the the resoutce that was updated.
	// -- This value lets the operation know which table contains the resource.
	// --
	// -- possible values:
	// -- * ManagedEnvironment (specified when we want Argo CD to C/R/U/D a user's cluster credentials)
	// -- * GitopsEngineInstance (specified to CRUD an Argo instance, for example to create a new namespace and put Argo CD in it, then signal when it's done)
	// -- * Application (user creates a new Application via service/web UI)
	Resource_type string `pg:"resource_type"`

	// -- When the operation was created. Used for garbage collection, as operations should be short lived.
	Created_on time.Time `pg:"created_on"`
	// created_on TIMESTAMP NOT NULL,

	// -- last_state_update is set whenever state changes
	// -- (initial value should be equal to created_on)
	Last_state_update time.Time `pg:"last_state_update"`
	// last_state_update TIMESTAMP NOT NULL,

	// -- possible values:
	// -- * Waiting
	// -- * In_Progress
	// -- * Completed
	// -- * Failed
	// state VARCHAR ( 30 ) NOT NULL,
	State string `pg:"state"`

	// -- If there is an error message from the operation, it is passed via this field.
	// human_readable_state VARCHAR ( 1024 )
	Human_readable_state string `pg:"human_readable_state"`
}

type Application struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"application"` //nolint

	Application_id string `pg:"application_id,pk"`

	SeqID int64 `pg:"seq_id"`

	// -- Name of the Application CR within the namespace
	Name string `pg:"name"`

	// -- resource_uid VARCHAR ( 48 ) NOT NULL UNIQUE,

	// -- '.spec' field of the Application CR
	// -- Note: Rather than converting individual JSON fields into SQL Table fields, we just pull the whole spec field.
	// -- In the future, it might be beneficial to pull out SOME of the fields, to reduce CPU time spent on json parsing
	Spec_field string `pg:"spec_field"`

	// -- Which Argo CD instance it's hosted on
	Engine_instance_inst_id string `pg:"engine_instance_inst_id"`

	// -- Which managed environment it is targetting
	// managed_environment_id VARCHAR(48) UNIQUE NOT NULL
	// -- CONSTRAINT fk_target_inf_cluster   FOREIGN KEY(target_inf_cluster)  REFERENCES InfrastructureCluster(inf_cluster_id),
	Managed_environment_id string `pg:"managed_environment_id"`
}

type ApplicationState struct {

	//lint:ignore U1000 used by go-pg
	tableName struct{} `pg:"applicationstate"` //nolint

	// -- Also a foreign key to Applicaiton.application_id
	Applicationstate_application_id string `pg:"applicationstate_application_id,pk"`
	// -- CONSTRAINT fk_app_id  PRIMARY KEY  FOREIGN KEY(app_id)  REFERENCES Application(appl_id),

	// -- Possible values:
	// -- * Healthy
	// -- * Progressing
	// -- * Degraded
	// -- * Suspended
	// -- * Missing
	// -- * Unknown
	Health string `pg:"health"`

	// -- Possible values:
	// -- * Synced
	// -- * OutOfSync
	Sync_Status string `pg:"sync_status"`

	// -- human_readable_health ( 512 ) NOT NULL,
	// -- human_readable_sync ( 512 ) NOT NULL,
	// -- human_readable_state ( 512 ) NOT NULL,

}
