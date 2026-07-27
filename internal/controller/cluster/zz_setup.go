// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	accountauthenticationpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/accountauthenticationpolicyattachment"
	accountpasswordpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/accountpasswordpolicyattachment"
	alert "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/alert"
	apiintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegration"
	apiintegrationamazonapigateway "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationamazonapigateway"
	apiintegrationazureapimanagement "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationazureapimanagement"
	apiintegrationexternalmcpdynamicclient "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationexternalmcpdynamicclient"
	apiintegrationexternalmcpoauth2 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationexternalmcpoauth2"
	apiintegrationgitrepositorygithubapp "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationgitrepositorygithubapp"
	apiintegrationgitrepositoryoauth2 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationgitrepositoryoauth2"
	apiintegrationgitrepositoryprivatelink "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationgitrepositoryprivatelink"
	apiintegrationgitrepositorytoken "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationgitrepositorytoken"
	apiintegrationgooglecloudapigateway "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/apiintegrationgooglecloudapigateway"
	cortexagent "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/cortexagent"
	cortexsearchservice "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/cortexsearchservice"
	dynamictable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/dynamictable"
	emailnotificationintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/emailnotificationintegration"
	externalfunction "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/externalfunction"
	externaltable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/externaltable"
	failovergroup "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/failovergroup"
	fileformat "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/fileformat"
	functionjava "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/functionjava"
	functionjavascript "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/functionjavascript"
	functionpython "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/functionpython"
	functionscala "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/functionscala"
	functionsql "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/functionsql"
	icebergtablefromdeltafiles "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/icebergtablefromdeltafiles"
	icebergtablefromfiles "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/icebergtablefromfiles"
	jobservice "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/jobservice"
	managedaccount "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/managedaccount"
	materializedview "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/materializedview"
	networkpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/networkpolicyattachment"
	notebook "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/notebook"
	notificationintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/notificationintegration"
	objectparameter "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/objectparameter"
	pipe "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/pipe"
	postgresinstance "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/postgresinstance"
	procedurejava "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/procedurejava"
	procedurejavascript "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/procedurejavascript"
	procedurepython "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/procedurepython"
	procedurescala "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/procedurescala"
	proceduresql "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/proceduresql"
	semanticview "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/semanticview"
	sequence "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/sequence"
	share "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/share"
	stage "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/stage"
	storageintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/storageintegration"
	storagelifecyclepolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/storagelifecyclepolicy"
	table "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/table"
	tablecolumnmaskingpolicyapplication "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/tablecolumnmaskingpolicyapplication"
	tableconstraint "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/tableconstraint"
	tablestoragelifecyclepolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/tablestoragelifecyclepolicyattachment"
	userauthenticationpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/userauthenticationpolicyattachment"
	userpasswordpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/userpasswordpolicyattachment"
	userpublickeys "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/userpublickeys"
	warehouseadaptive "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/preview/warehouseadaptive"
	providerconfig "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/providerconfig"
	account "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/account"
	accountparameter "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/accountparameter"
	accountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/accountrole"
	accountsessionpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/accountsessionpolicyattachment"
	apiauthenticationintegrationwithauthorizationcodegrant "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/apiauthenticationintegrationwithauthorizationcodegrant"
	apiauthenticationintegrationwithclientcredentials "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/apiauthenticationintegrationwithclientcredentials"
	apiauthenticationintegrationwithjwtbearer "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/apiauthenticationintegrationwithjwtbearer"
	authenticationpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/authenticationpolicy"
	catalogintegrationawsglue "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/catalogintegrationawsglue"
	catalogintegrationicebergrest "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/catalogintegrationicebergrest"
	catalogintegrationobjectstorage "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/catalogintegrationobjectstorage"
	catalogintegrationopencatalog "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/catalogintegrationopencatalog"
	computepool "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/computepool"
	currentaccount "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/currentaccount"
	currentorganizationaccount "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/currentorganizationaccount"
	database "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/database"
	databaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/databaserole"
	execute "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/execute"
	externaloauthintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/externaloauthintegration"
	externalvolume "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/externalvolume"
	gitrepository "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/gitrepository"
	grantaccountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/grantaccountrole"
	grantapplicationrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/grantapplicationrole"
	grantdatabaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/grantdatabaserole"
	grantownership "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/grantownership"
	grantprivilegestoaccountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/grantprivilegestoaccountrole"
	grantprivilegestodatabaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/grantprivilegestodatabaserole"
	grantprivilegestoshare "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/grantprivilegestoshare"
	imagerepository "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/imagerepository"
	legacyserviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/legacyserviceuser"
	listing "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/listing"
	maskingpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/maskingpolicy"
	networkpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/networkpolicy"
	networkrule "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/networkrule"
	oauthintegrationforcustomclients "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/oauthintegrationforcustomclients"
	oauthintegrationforpartnerapplications "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/oauthintegrationforpartnerapplications"
	passwordpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/passwordpolicy"
	primaryconnection "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/primaryconnection"
	resourcemonitor "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/resourcemonitor"
	rowaccesspolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/rowaccesspolicy"
	saml2integration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/saml2integration"
	schema "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/schema"
	scimintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/scimintegration"
	secondaryconnection "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/secondaryconnection"
	secondarydatabase "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/secondarydatabase"
	secretwithauthorizationcodegrant "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/secretwithauthorizationcodegrant"
	secretwithbasicauthentication "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/secretwithbasicauthentication"
	secretwithclientcredentials "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/secretwithclientcredentials"
	secretwithgenericstring "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/secretwithgenericstring"
	service "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/service"
	serviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/serviceuser"
	sessionpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/sessionpolicy"
	shareddatabase "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/shareddatabase"
	stageexternalazure "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/stageexternalazure"
	stageexternalgcs "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/stageexternalgcs"
	stageexternals3 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/stageexternals3"
	stageexternals3compatible "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/stageexternals3compatible"
	stageinternal "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/stageinternal"
	storageintegrationaws "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/storageintegrationaws"
	storageintegrationazure "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/storageintegrationazure"
	storageintegrationgcs "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/storageintegrationgcs"
	streamlit "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/streamlit"
	streamondirectorytable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/streamondirectorytable"
	streamonexternaltable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/streamonexternaltable"
	streamontable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/streamontable"
	streamonview "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/streamonview"
	tag "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/tag"
	tagassociation "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/tagassociation"
	task "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/task"
	user "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/user"
	userprogrammaticaccesstoken "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/userprogrammaticaccesstoken"
	usersessionpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/usersessionpolicyattachment"
	view "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/view"
	warehouse "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/cluster/stable/warehouse"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accountauthenticationpolicyattachment.Setup,
		accountpasswordpolicyattachment.Setup,
		alert.Setup,
		apiintegration.Setup,
		apiintegrationamazonapigateway.Setup,
		apiintegrationazureapimanagement.Setup,
		apiintegrationexternalmcpdynamicclient.Setup,
		apiintegrationexternalmcpoauth2.Setup,
		apiintegrationgitrepositorygithubapp.Setup,
		apiintegrationgitrepositoryoauth2.Setup,
		apiintegrationgitrepositoryprivatelink.Setup,
		apiintegrationgitrepositorytoken.Setup,
		apiintegrationgooglecloudapigateway.Setup,
		cortexagent.Setup,
		cortexsearchservice.Setup,
		dynamictable.Setup,
		emailnotificationintegration.Setup,
		externalfunction.Setup,
		externaltable.Setup,
		failovergroup.Setup,
		fileformat.Setup,
		functionjava.Setup,
		functionjavascript.Setup,
		functionpython.Setup,
		functionscala.Setup,
		functionsql.Setup,
		icebergtablefromdeltafiles.Setup,
		icebergtablefromfiles.Setup,
		jobservice.Setup,
		managedaccount.Setup,
		materializedview.Setup,
		networkpolicyattachment.Setup,
		notebook.Setup,
		notificationintegration.Setup,
		objectparameter.Setup,
		pipe.Setup,
		postgresinstance.Setup,
		procedurejava.Setup,
		procedurejavascript.Setup,
		procedurepython.Setup,
		procedurescala.Setup,
		proceduresql.Setup,
		semanticview.Setup,
		sequence.Setup,
		share.Setup,
		stage.Setup,
		storageintegration.Setup,
		storagelifecyclepolicy.Setup,
		table.Setup,
		tablecolumnmaskingpolicyapplication.Setup,
		tableconstraint.Setup,
		tablestoragelifecyclepolicyattachment.Setup,
		userauthenticationpolicyattachment.Setup,
		userpasswordpolicyattachment.Setup,
		userpublickeys.Setup,
		warehouseadaptive.Setup,
		providerconfig.Setup,
		account.Setup,
		accountparameter.Setup,
		accountrole.Setup,
		accountsessionpolicyattachment.Setup,
		apiauthenticationintegrationwithauthorizationcodegrant.Setup,
		apiauthenticationintegrationwithclientcredentials.Setup,
		apiauthenticationintegrationwithjwtbearer.Setup,
		authenticationpolicy.Setup,
		catalogintegrationawsglue.Setup,
		catalogintegrationicebergrest.Setup,
		catalogintegrationobjectstorage.Setup,
		catalogintegrationopencatalog.Setup,
		computepool.Setup,
		currentaccount.Setup,
		currentorganizationaccount.Setup,
		database.Setup,
		databaserole.Setup,
		execute.Setup,
		externaloauthintegration.Setup,
		externalvolume.Setup,
		gitrepository.Setup,
		grantaccountrole.Setup,
		grantapplicationrole.Setup,
		grantdatabaserole.Setup,
		grantownership.Setup,
		grantprivilegestoaccountrole.Setup,
		grantprivilegestodatabaserole.Setup,
		grantprivilegestoshare.Setup,
		imagerepository.Setup,
		legacyserviceuser.Setup,
		listing.Setup,
		maskingpolicy.Setup,
		networkpolicy.Setup,
		networkrule.Setup,
		oauthintegrationforcustomclients.Setup,
		oauthintegrationforpartnerapplications.Setup,
		passwordpolicy.Setup,
		primaryconnection.Setup,
		resourcemonitor.Setup,
		rowaccesspolicy.Setup,
		saml2integration.Setup,
		schema.Setup,
		scimintegration.Setup,
		secondaryconnection.Setup,
		secondarydatabase.Setup,
		secretwithauthorizationcodegrant.Setup,
		secretwithbasicauthentication.Setup,
		secretwithclientcredentials.Setup,
		secretwithgenericstring.Setup,
		service.Setup,
		serviceuser.Setup,
		sessionpolicy.Setup,
		shareddatabase.Setup,
		stageexternalazure.Setup,
		stageexternalgcs.Setup,
		stageexternals3.Setup,
		stageexternals3compatible.Setup,
		stageinternal.Setup,
		storageintegrationaws.Setup,
		storageintegrationazure.Setup,
		storageintegrationgcs.Setup,
		streamlit.Setup,
		streamondirectorytable.Setup,
		streamonexternaltable.Setup,
		streamontable.Setup,
		streamonview.Setup,
		tag.Setup,
		tagassociation.Setup,
		task.Setup,
		user.Setup,
		userprogrammaticaccesstoken.Setup,
		usersessionpolicyattachment.Setup,
		view.Setup,
		warehouse.Setup,
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
		accountauthenticationpolicyattachment.SetupGated,
		accountpasswordpolicyattachment.SetupGated,
		alert.SetupGated,
		apiintegration.SetupGated,
		apiintegrationamazonapigateway.SetupGated,
		apiintegrationazureapimanagement.SetupGated,
		apiintegrationexternalmcpdynamicclient.SetupGated,
		apiintegrationexternalmcpoauth2.SetupGated,
		apiintegrationgitrepositorygithubapp.SetupGated,
		apiintegrationgitrepositoryoauth2.SetupGated,
		apiintegrationgitrepositoryprivatelink.SetupGated,
		apiintegrationgitrepositorytoken.SetupGated,
		apiintegrationgooglecloudapigateway.SetupGated,
		cortexagent.SetupGated,
		cortexsearchservice.SetupGated,
		dynamictable.SetupGated,
		emailnotificationintegration.SetupGated,
		externalfunction.SetupGated,
		externaltable.SetupGated,
		failovergroup.SetupGated,
		fileformat.SetupGated,
		functionjava.SetupGated,
		functionjavascript.SetupGated,
		functionpython.SetupGated,
		functionscala.SetupGated,
		functionsql.SetupGated,
		icebergtablefromdeltafiles.SetupGated,
		icebergtablefromfiles.SetupGated,
		jobservice.SetupGated,
		managedaccount.SetupGated,
		materializedview.SetupGated,
		networkpolicyattachment.SetupGated,
		notebook.SetupGated,
		notificationintegration.SetupGated,
		objectparameter.SetupGated,
		pipe.SetupGated,
		postgresinstance.SetupGated,
		procedurejava.SetupGated,
		procedurejavascript.SetupGated,
		procedurepython.SetupGated,
		procedurescala.SetupGated,
		proceduresql.SetupGated,
		semanticview.SetupGated,
		sequence.SetupGated,
		share.SetupGated,
		stage.SetupGated,
		storageintegration.SetupGated,
		storagelifecyclepolicy.SetupGated,
		table.SetupGated,
		tablecolumnmaskingpolicyapplication.SetupGated,
		tableconstraint.SetupGated,
		tablestoragelifecyclepolicyattachment.SetupGated,
		userauthenticationpolicyattachment.SetupGated,
		userpasswordpolicyattachment.SetupGated,
		userpublickeys.SetupGated,
		warehouseadaptive.SetupGated,
		providerconfig.SetupGated,
		account.SetupGated,
		accountparameter.SetupGated,
		accountrole.SetupGated,
		accountsessionpolicyattachment.SetupGated,
		apiauthenticationintegrationwithauthorizationcodegrant.SetupGated,
		apiauthenticationintegrationwithclientcredentials.SetupGated,
		apiauthenticationintegrationwithjwtbearer.SetupGated,
		authenticationpolicy.SetupGated,
		catalogintegrationawsglue.SetupGated,
		catalogintegrationicebergrest.SetupGated,
		catalogintegrationobjectstorage.SetupGated,
		catalogintegrationopencatalog.SetupGated,
		computepool.SetupGated,
		currentaccount.SetupGated,
		currentorganizationaccount.SetupGated,
		database.SetupGated,
		databaserole.SetupGated,
		execute.SetupGated,
		externaloauthintegration.SetupGated,
		externalvolume.SetupGated,
		gitrepository.SetupGated,
		grantaccountrole.SetupGated,
		grantapplicationrole.SetupGated,
		grantdatabaserole.SetupGated,
		grantownership.SetupGated,
		grantprivilegestoaccountrole.SetupGated,
		grantprivilegestodatabaserole.SetupGated,
		grantprivilegestoshare.SetupGated,
		imagerepository.SetupGated,
		legacyserviceuser.SetupGated,
		listing.SetupGated,
		maskingpolicy.SetupGated,
		networkpolicy.SetupGated,
		networkrule.SetupGated,
		oauthintegrationforcustomclients.SetupGated,
		oauthintegrationforpartnerapplications.SetupGated,
		passwordpolicy.SetupGated,
		primaryconnection.SetupGated,
		resourcemonitor.SetupGated,
		rowaccesspolicy.SetupGated,
		saml2integration.SetupGated,
		schema.SetupGated,
		scimintegration.SetupGated,
		secondaryconnection.SetupGated,
		secondarydatabase.SetupGated,
		secretwithauthorizationcodegrant.SetupGated,
		secretwithbasicauthentication.SetupGated,
		secretwithclientcredentials.SetupGated,
		secretwithgenericstring.SetupGated,
		service.SetupGated,
		serviceuser.SetupGated,
		sessionpolicy.SetupGated,
		shareddatabase.SetupGated,
		stageexternalazure.SetupGated,
		stageexternalgcs.SetupGated,
		stageexternals3.SetupGated,
		stageexternals3compatible.SetupGated,
		stageinternal.SetupGated,
		storageintegrationaws.SetupGated,
		storageintegrationazure.SetupGated,
		storageintegrationgcs.SetupGated,
		streamlit.SetupGated,
		streamondirectorytable.SetupGated,
		streamonexternaltable.SetupGated,
		streamontable.SetupGated,
		streamonview.SetupGated,
		tag.SetupGated,
		tagassociation.SetupGated,
		task.SetupGated,
		user.SetupGated,
		userprogrammaticaccesstoken.SetupGated,
		usersessionpolicyattachment.SetupGated,
		view.SetupGated,
		warehouse.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
