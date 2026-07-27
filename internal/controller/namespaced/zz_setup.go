// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	account "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/account"
	accountauthenticationpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/accountauthenticationpolicyattachment"
	accountpasswordpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/accountpasswordpolicyattachment"
	adaptive "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/adaptive"
	agent "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/agent"
	alert "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/alert"
	authenticationpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/authenticationpolicyattachment"
	columnmaskingpolicyapplication "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/columnmaskingpolicyapplication"
	constraint "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/constraint"
	format "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/format"
	function "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/function"
	group "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/group"
	instance "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/instance"
	integration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integration"
	integrationamazonapigateway "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationamazonapigateway"
	integrationazureapimanagement "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationazureapimanagement"
	integrationexternalmcpdynamicclient "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationexternalmcpdynamicclient"
	integrationexternalmcpoauth2 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationexternalmcpoauth2"
	integrationgitrepositorygithubapp "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationgitrepositorygithubapp"
	integrationgitrepositoryoauth2 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationgitrepositoryoauth2"
	integrationgitrepositoryprivatelink "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationgitrepositoryprivatelink"
	integrationgitrepositorytoken "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationgitrepositorytoken"
	integrationgooglecloudapigateway "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/integrationgooglecloudapigateway"
	java "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/java"
	javascript "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/javascript"
	lifecyclepolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/lifecyclepolicy"
	notebook "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/notebook"
	notificationintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/notificationintegration"
	parameter "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/parameter"
	passwordpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/passwordpolicyattachment"
	pipe "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/pipe"
	policyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/policyattachment"
	publickeys "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/publickeys"
	python "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/python"
	scala "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/scala"
	searchservice "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/searchservice"
	sequence "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/sequence"
	service "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/service"
	share "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/share"
	sql "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/sql"
	stage "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/stage"
	storagelifecyclepolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/storagelifecyclepolicyattachment"
	table "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/table"
	tablefromdeltafiles "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/tablefromdeltafiles"
	tablefromfiles "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/tablefromfiles"
	view "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/view"
	providerconfig "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/providerconfig"
	accountstable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/account"
	accountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/accountrole"
	accountsessionpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/accountsessionpolicyattachment"
	applicationrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/applicationrole"
	association "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/association"
	authenticationintegrationwithauthorizationcodegrant "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/authenticationintegrationwithauthorizationcodegrant"
	authenticationintegrationwithclientcredentials "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/authenticationintegrationwithclientcredentials"
	authenticationintegrationwithjwtbearer "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/authenticationintegrationwithjwtbearer"
	authenticationpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/authenticationpolicy"
	connection "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/connection"
	database "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/database"
	databaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/databaserole"
	execute "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/execute"
	externalazure "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/externalazure"
	externalgcs "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/externalgcs"
	externals3 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/externals3"
	externals3compatible "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/externals3compatible"
	grantaccountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantaccountrole"
	integrationstable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integration"
	integrationaws "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationaws"
	integrationawsglue "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationawsglue"
	integrationazure "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationazure"
	integrationforcustomclients "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationforcustomclients"
	integrationforpartnerapplications "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationforpartnerapplications"
	integrationgcs "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationgcs"
	integrationicebergrest "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationicebergrest"
	integrationobjectstorage "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationobjectstorage"
	integrationopencatalog "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/integrationopencatalog"
	legacyserviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/legacyserviceuser"
	listing "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/listing"
	maskingpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/maskingpolicy"
	monitor "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/monitor"
	oauthintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/oauthintegration"
	ondirectorytable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/ondirectorytable"
	onexternaltable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/onexternaltable"
	ontable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/ontable"
	onview "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/onview"
	organizationaccount "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/organizationaccount"
	ownership "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/ownership"
	parameterstable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/parameter"
	passwordpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/passwordpolicy"
	policy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/policy"
	pool "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/pool"
	privilegestoaccountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/privilegestoaccountrole"
	privilegestodatabaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/privilegestodatabaserole"
	privilegestoshare "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/privilegestoshare"
	programmaticaccesstoken "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/programmaticaccesstoken"
	repository "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/repository"
	role "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/role"
	rowaccesspolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/rowaccesspolicy"
	rule "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/rule"
	schema "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/schema"
	secondarydatabase "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/secondarydatabase"
	servicestable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/service"
	serviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/serviceuser"
	sessionpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/sessionpolicy"
	sessionpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/sessionpolicyattachment"
	shareddatabase "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/shareddatabase"
	stageinternal "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/stageinternal"
	streamlit "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/streamlit"
	tag "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/tag"
	task "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/task"
	user "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/user"
	viewstable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/view"
	volume "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/volume"
	warehouse "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/warehouse"
	withauthorizationcodegrant "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/withauthorizationcodegrant"
	withbasicauthentication "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/withbasicauthentication"
	withclientcredentials "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/withclientcredentials"
	withgenericstring "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/withgenericstring"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		account.Setup,
		accountauthenticationpolicyattachment.Setup,
		accountpasswordpolicyattachment.Setup,
		adaptive.Setup,
		agent.Setup,
		alert.Setup,
		authenticationpolicyattachment.Setup,
		columnmaskingpolicyapplication.Setup,
		constraint.Setup,
		format.Setup,
		function.Setup,
		group.Setup,
		instance.Setup,
		integration.Setup,
		integration.Setup,
		integration.Setup,
		integrationamazonapigateway.Setup,
		integrationazureapimanagement.Setup,
		integrationexternalmcpdynamicclient.Setup,
		integrationexternalmcpoauth2.Setup,
		integrationgitrepositorygithubapp.Setup,
		integrationgitrepositoryoauth2.Setup,
		integrationgitrepositoryprivatelink.Setup,
		integrationgitrepositorytoken.Setup,
		integrationgooglecloudapigateway.Setup,
		java.Setup,
		java.Setup,
		javascript.Setup,
		javascript.Setup,
		lifecyclepolicy.Setup,
		notebook.Setup,
		notificationintegration.Setup,
		parameter.Setup,
		passwordpolicyattachment.Setup,
		pipe.Setup,
		policyattachment.Setup,
		publickeys.Setup,
		python.Setup,
		python.Setup,
		scala.Setup,
		scala.Setup,
		searchservice.Setup,
		sequence.Setup,
		service.Setup,
		share.Setup,
		sql.Setup,
		sql.Setup,
		stage.Setup,
		storagelifecyclepolicyattachment.Setup,
		table.Setup,
		table.Setup,
		table.Setup,
		tablefromdeltafiles.Setup,
		tablefromfiles.Setup,
		view.Setup,
		view.Setup,
		providerconfig.Setup,
		accountstable.Setup,
		accountstable.Setup,
		accountrole.Setup,
		accountsessionpolicyattachment.Setup,
		applicationrole.Setup,
		association.Setup,
		authenticationintegrationwithauthorizationcodegrant.Setup,
		authenticationintegrationwithclientcredentials.Setup,
		authenticationintegrationwithjwtbearer.Setup,
		authenticationpolicy.Setup,
		connection.Setup,
		connection.Setup,
		database.Setup,
		databaserole.Setup,
		execute.Setup,
		externalazure.Setup,
		externalgcs.Setup,
		externals3.Setup,
		externals3compatible.Setup,
		grantaccountrole.Setup,
		integrationstable.Setup,
		integrationstable.Setup,
		integrationaws.Setup,
		integrationawsglue.Setup,
		integrationazure.Setup,
		integrationforcustomclients.Setup,
		integrationforpartnerapplications.Setup,
		integrationgcs.Setup,
		integrationicebergrest.Setup,
		integrationobjectstorage.Setup,
		integrationopencatalog.Setup,
		legacyserviceuser.Setup,
		listing.Setup,
		maskingpolicy.Setup,
		monitor.Setup,
		oauthintegration.Setup,
		ondirectorytable.Setup,
		onexternaltable.Setup,
		ontable.Setup,
		onview.Setup,
		organizationaccount.Setup,
		ownership.Setup,
		parameterstable.Setup,
		passwordpolicy.Setup,
		policy.Setup,
		pool.Setup,
		privilegestoaccountrole.Setup,
		privilegestodatabaserole.Setup,
		privilegestoshare.Setup,
		programmaticaccesstoken.Setup,
		repository.Setup,
		repository.Setup,
		role.Setup,
		rowaccesspolicy.Setup,
		rule.Setup,
		schema.Setup,
		secondarydatabase.Setup,
		servicestable.Setup,
		serviceuser.Setup,
		sessionpolicy.Setup,
		sessionpolicyattachment.Setup,
		shareddatabase.Setup,
		stageinternal.Setup,
		streamlit.Setup,
		tag.Setup,
		task.Setup,
		user.Setup,
		viewstable.Setup,
		volume.Setup,
		warehouse.Setup,
		withauthorizationcodegrant.Setup,
		withbasicauthentication.Setup,
		withclientcredentials.Setup,
		withgenericstring.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		account.SetupGated,
		accountauthenticationpolicyattachment.SetupGated,
		accountpasswordpolicyattachment.SetupGated,
		adaptive.SetupGated,
		agent.SetupGated,
		alert.SetupGated,
		authenticationpolicyattachment.SetupGated,
		columnmaskingpolicyapplication.SetupGated,
		constraint.SetupGated,
		format.SetupGated,
		function.SetupGated,
		group.SetupGated,
		instance.SetupGated,
		integration.SetupGated,
		integration.SetupGated,
		integration.SetupGated,
		integrationamazonapigateway.SetupGated,
		integrationazureapimanagement.SetupGated,
		integrationexternalmcpdynamicclient.SetupGated,
		integrationexternalmcpoauth2.SetupGated,
		integrationgitrepositorygithubapp.SetupGated,
		integrationgitrepositoryoauth2.SetupGated,
		integrationgitrepositoryprivatelink.SetupGated,
		integrationgitrepositorytoken.SetupGated,
		integrationgooglecloudapigateway.SetupGated,
		java.SetupGated,
		java.SetupGated,
		javascript.SetupGated,
		javascript.SetupGated,
		lifecyclepolicy.SetupGated,
		notebook.SetupGated,
		notificationintegration.SetupGated,
		parameter.SetupGated,
		passwordpolicyattachment.SetupGated,
		pipe.SetupGated,
		policyattachment.SetupGated,
		publickeys.SetupGated,
		python.SetupGated,
		python.SetupGated,
		scala.SetupGated,
		scala.SetupGated,
		searchservice.SetupGated,
		sequence.SetupGated,
		service.SetupGated,
		share.SetupGated,
		sql.SetupGated,
		sql.SetupGated,
		stage.SetupGated,
		storagelifecyclepolicyattachment.SetupGated,
		table.SetupGated,
		table.SetupGated,
		table.SetupGated,
		tablefromdeltafiles.SetupGated,
		tablefromfiles.SetupGated,
		view.SetupGated,
		view.SetupGated,
		providerconfig.SetupGated,
		accountstable.SetupGated,
		accountstable.SetupGated,
		accountrole.SetupGated,
		accountsessionpolicyattachment.SetupGated,
		applicationrole.SetupGated,
		association.SetupGated,
		authenticationintegrationwithauthorizationcodegrant.SetupGated,
		authenticationintegrationwithclientcredentials.SetupGated,
		authenticationintegrationwithjwtbearer.SetupGated,
		authenticationpolicy.SetupGated,
		connection.SetupGated,
		connection.SetupGated,
		database.SetupGated,
		databaserole.SetupGated,
		execute.SetupGated,
		externalazure.SetupGated,
		externalgcs.SetupGated,
		externals3.SetupGated,
		externals3compatible.SetupGated,
		grantaccountrole.SetupGated,
		integrationstable.SetupGated,
		integrationstable.SetupGated,
		integrationaws.SetupGated,
		integrationawsglue.SetupGated,
		integrationazure.SetupGated,
		integrationforcustomclients.SetupGated,
		integrationforpartnerapplications.SetupGated,
		integrationgcs.SetupGated,
		integrationicebergrest.SetupGated,
		integrationobjectstorage.SetupGated,
		integrationopencatalog.SetupGated,
		legacyserviceuser.SetupGated,
		listing.SetupGated,
		maskingpolicy.SetupGated,
		monitor.SetupGated,
		oauthintegration.SetupGated,
		ondirectorytable.SetupGated,
		onexternaltable.SetupGated,
		ontable.SetupGated,
		onview.SetupGated,
		organizationaccount.SetupGated,
		ownership.SetupGated,
		parameterstable.SetupGated,
		passwordpolicy.SetupGated,
		policy.SetupGated,
		pool.SetupGated,
		privilegestoaccountrole.SetupGated,
		privilegestodatabaserole.SetupGated,
		privilegestoshare.SetupGated,
		programmaticaccesstoken.SetupGated,
		repository.SetupGated,
		repository.SetupGated,
		role.SetupGated,
		rowaccesspolicy.SetupGated,
		rule.SetupGated,
		schema.SetupGated,
		secondarydatabase.SetupGated,
		servicestable.SetupGated,
		serviceuser.SetupGated,
		sessionpolicy.SetupGated,
		sessionpolicyattachment.SetupGated,
		shareddatabase.SetupGated,
		stageinternal.SetupGated,
		streamlit.SetupGated,
		tag.SetupGated,
		task.SetupGated,
		user.SetupGated,
		viewstable.SetupGated,
		volume.SetupGated,
		warehouse.SetupGated,
		withauthorizationcodegrant.SetupGated,
		withbasicauthentication.SetupGated,
		withclientcredentials.SetupGated,
		withgenericstring.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
